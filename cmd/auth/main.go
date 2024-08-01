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
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
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

	dbClient := db.New(dbConn)
	uc := usecase.New(dbClient, chr, bot, cfg.JwtSecret)
	h := handler.New(uc)
	rtr := router.Init(h)

	errCh := make(chan error)

	srv := server.New(rtr)
	cj := cronjob.NewCronJob(uc)
	err = cj.Run(ctx)
	if err != nil {
		log.Fatal().Msgf("failed adding cron job, err: %v", err)
	}

	go func() {
		if err = srv.Run("0.0.0.0:1234"); err != nil {
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
		logrus.Errorf("Failed to run, err: %v", err)
	case res := <-sigChan:
		if res == syscall.SIGINT || res == syscall.SIGTERM {
			logrus.Info("Signal received")
		} else if res == syscall.SIGHUP {
			logrus.Info("Signal received")
		}
	}

	logrus.Print("Shutting Down")

	if err = srv.Shutdown(ctx); err != nil {
		logrus.Errorf("failed server shutting down: %s", err.Error())
	}

	if err = dbClient.Close(); err != nil {
		logrus.Errorf("failed db connection close: %s", err.Error())
	}

	cj.Stop()
}
