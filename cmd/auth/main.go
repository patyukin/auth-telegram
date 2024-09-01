package main

import (
	"auth-telegram/internal/cacher"
	"auth-telegram/internal/config"
	"auth-telegram/internal/cronjob"
	"auth-telegram/internal/db"
	"auth-telegram/internal/dbconn"
	"auth-telegram/internal/handler"
	"auth-telegram/internal/migrator"
	"auth-telegram/internal/server"
	"auth-telegram/internal/server/router"
	"auth-telegram/internal/telegram"
	"auth-telegram/internal/usecase"
	"auth-telegram/pkg/meter"
	"auth-telegram/pkg/tracer"
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Msgf("unable to load config: %v", err)
	}

	dbConn, err := dbconn.New(ctx, cfg)
	if err != nil {
		log.Fatal().Msgf("failed connecting to db: %v", err)
	}

	err = migrator.UpMigrations(ctx, dbConn)
	if err != nil {
		log.Fatal().Msgf("failed migrating db: %v", err)
	}

	chr, err := cacher.New(ctx, cfg)
	if err != nil {
		log.Fatal().Msgf("failed creating redis client: %v", err)
	}

	bot, err := telegram.New(cfg)
	if err != nil {
		log.Fatal().Msgf("failed creating telegram bot: %v", err)
	}

	traceProv, err := tracer.InitTracer("http://localhost:14268/api/traces", "Auth Service")
	if err != nil {
		log.Fatal().Msgf("init tracer, err: %v", err)
	}

	meterProv, err := meter.InitMeter(ctx, "Auth Service")
	if err != nil {
		log.Fatal().Msgf("init meter, err: %v", err)
	}

	srvAddress := fmt.Sprintf(":%d", cfg.HttpPort)
	dbClient := db.New(dbConn)
	uc := usecase.New(dbClient, chr, bot, cfg.JwtSecret)
	h := handler.New(uc)
	rtr := router.Init(h, srvAddress)

	errCh := make(chan error)

	srv := server.New(rtr)
	cj := cronjob.NewCronJob(uc)

	go func() {
		if err = cj.Run(ctx); err != nil {
			log.Error().Msgf("failed adding cron job, err: %v", err)
			errCh <- err
		}
	}()

	go func() {
		if err = srv.Run(srvAddress, cfg); err != nil {
			log.Error().Msgf("failed running http server: %v", err)
			errCh <- err
		}
	}()

	go uc.StartTelegramBot(ctx)

	log.Info().Msg("Auth App Started")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err = <-errCh:
		log.Error().Msgf("Failed to run, err: %v", err)
	case res := <-sigChan:
		if res == syscall.SIGINT || res == syscall.SIGTERM {
			log.Info().Msgf("Signal received, exiting...")
		} else if res == syscall.SIGHUP {
			log.Info().Msgf("Signal received, sighup")
		}
	}

	log.Info().Msgf("Shutting Down")

	if err = srv.Shutdown(ctx); err != nil {
		log.Error().Msgf("failed server shutting down: %s", err.Error())
	}

	if err = dbClient.Close(); err != nil {
		log.Error().Msgf("failed db connection close: %s", err.Error())
	}

	if err = traceProv.Shutdown(ctx); err != nil {
		log.Error().Msgf("Error shutting down tracer provider: %v", err)
	}

	if err = meterProv.Shutdown(ctx); err != nil {
		log.Error().Msgf("Error shutting down meter provider: %v", err)
	}

	cj.Stop()
}
