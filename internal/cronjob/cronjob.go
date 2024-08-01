package cronjob

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

type CronJob struct {
	c  *cron.Cron
	uc UseCase
}

type UseCase interface {
	CleanTokens(ctx context.Context) error
}

func NewCronJob(uc UseCase) *CronJob {
	return &CronJob{
		c:  cron.New(),
		uc: uc,
	}
}

func (cj *CronJob) Stop() {
	cj.c.Stop()
}

func (cj *CronJob) Run(ctx context.Context) error {
	var err error
	_, err = cj.c.AddFunc("*/10 * * * *", func() {
		log.Info().Msg("run cleaning tokens")

		err = cj.uc.CleanTokens(ctx)
		if err != nil {
			log.Error().Msgf("failed cleaning tokens, err: %v", err)
		}
	})
	if err != nil {
		return fmt.Errorf("failed adding cron job: %w", err)
	}

	cj.c.Start()
	return nil
}

// https://crontab.guru/
