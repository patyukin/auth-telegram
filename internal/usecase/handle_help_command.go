package usecase

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func (uc *UseCase) handleHelpCommand(message *tgbotapi.Message) {
	helpText := `
Available commands:
/start <uuid> - Start the bot with an invite link
/generate_code - Generate a temporary code with 60 seconds validity
/help - Show this help message
`
	msg := tgbotapi.NewMessage(message.Chat.ID, helpText)
	if _, err := uc.bot.API.Send(msg); err != nil {
		log.Error().Msgf("Error sending message: %v", err)
		return
	}

	log.Info().Msgf("Sent message to chat %d", message.Chat.ID)
}
