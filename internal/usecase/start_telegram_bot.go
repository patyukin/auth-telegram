package usecase

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func (uc *UseCase) StartTelegramBot(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := uc.bot.API.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Info().Msgf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			switch update.Message.Command() {
			case "start":
				uc.handleStartCommand(ctx, update.Message)
			case "help":
				uc.handleHelpCommand(update.Message)
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command. Use /help to see the list of available commands.")
				if _, err := uc.bot.API.Send(msg); err != nil {
					log.Error().Msgf("Error sending message: %v", err)
					return
				}

				log.Info().Msgf("Sent message to chat %d", update.Message.Chat.ID)
			}
		}
	}
}
