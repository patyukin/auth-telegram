package usecase

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"strings"
)

func (uc *UseCase) handleStartCommand(ctx context.Context, message *tgbotapi.Message) {
	args := message.CommandArguments()

	if args == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Welcome! Please provide a valid invite link to start.")
		if _, err := uc.bot.API.Send(msg); err != nil {
			log.Error().Msgf("Error sending message: %v", err)
			return
		}

		log.Info().Msgf("Sent message to chat %d", message.Chat.ID)
		return
	}

	code, err := uuid.Parse(args)
	if err != nil {
		log.Error().Msgf("Error parsing invite code: %v", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid invite link.")
		if _, err = uc.bot.API.Send(msg); err != nil {
			log.Error().Msgf("Error sending message: %v", err)
			return
		}

		log.Info().Msgf("Sent message to chat %d", message.Chat.ID)
		return
	}

	allSignUpCode, err := uc.redis.GetSignUpCode(ctx, message.From.UserName)
	if err != nil {
		log.Error().Msgf("Error getting sign up code: %v", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid invite link.")
		if _, err = uc.bot.API.Send(msg); err != nil {
			log.Error().Msgf("Error sending message: %v", err)
			return
		}

		log.Info().Msgf("Sent message to chat %d", message.Chat.ID)
		return
	}

	result := strings.SplitN(allSignUpCode, ":", 2)

	signUpCode, err := uuid.Parse(result[0])
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Error was encountered. Please try again.")
		if _, err = uc.bot.API.Send(msg); err != nil {
			log.Error().Msgf("Error sending message: %v", err)
			return
		}

		log.Info().Msgf("Sent message to chat %d", message.Chat.ID)
		return
	}

	if signUpCode.String() != code.String() {
		log.Error().Msgf("Error parsing invite code: %v != %v", signUpCode.String(), code.String())
		msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid invite link.")
		if _, err = uc.bot.API.Send(msg); err != nil {
			log.Error().Msgf("Error sending message: %v", err)
			return
		}

		log.Info().Msgf("Sent message to chat %d", message.Chat.ID)
		return
	}

	err = uc.redis.DeleteSignUpCode(ctx, message.From.UserName)
	if err != nil {
		log.Error().Msgf("Error deleting sign up code: %v", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Error was encountered. Please try again.")
		if _, err = uc.bot.API.Send(msg); err != nil {
			log.Error().Msgf("Error sending message: %v", err)
			return
		}

		log.Info().Msgf("Sent message to chat %d", message.Chat.ID)
		return
	}

	userUUID, err := uuid.Parse(result[1])
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Error was encountered. Please try again.")
		if _, err = uc.bot.API.Send(msg); err != nil {
			log.Error().Msgf("Error sending message: %v", err)
			return
		}

		log.Info().Msgf("Sent message to chat %d", message.Chat.ID)
		return
	}

	err = uc.registry.GetRepo().UpdateTelegramUserAfterSignUp(ctx, userUUID, message.Chat.ID, message.From.ID)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Error was encountered. Please try again.")
		if _, err = uc.bot.API.Send(msg); err != nil {
			log.Error().Msgf("Error sending message: %v", err)
			return
		}

		log.Info().Msgf("Sent message to chat %d", message.Chat.ID)
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "You have successfully signed up!")
	if _, err = uc.bot.API.Send(msg); err != nil {
		log.Error().Msgf("Error sending message: %v", err)
		return
	}

	log.Info().Msgf("Sent message to chat %d", message.Chat.ID)
}
