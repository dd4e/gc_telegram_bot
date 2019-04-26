package handlers

import (
	"github.com/dadmoscow/gc_telegram_bot/app/data"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type TGBotHandlers struct {
	BotAPI *tgbotapi.BotAPI
}

func (bh TGBotHandlers) DeleteMessage(msg data.Message) error {
	delMsg := tgbotapi.DeleteMessageConfig{
		MessageID: msg.MsgID,
		ChatID:    msg.ChatID,
	}
	resp, err := bh.BotAPI.DeleteMessage(delMsg)

	if err != nil {
		switch resp.ErrorCode {
		case 400:
			log.Printf("Warning: %s from chat %s", resp.Description, msg.ChatConfig)
		default:
			log.Printf("Error: %s. The message %s will be deleted later from chat %s.",
				err, msg, msg.ChatConfig)
			return err
		}
	}

	log.Printf("The message %s has been deleted from chat %s",
		msg, msg.ChatConfig)
	return nil
}
