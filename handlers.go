package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
	"time"
)

// new message handler
func botUpdateMsgHandler(cmdChan chan *tgbotapi.Message) {
	log.Println("Start new message handler")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, _ := BOT.GetUpdatesChan(u)

	for update := range updates {
		msg := update.Message
		// skip non message updates
		if msg == nil {
			continue
		}

		// process message only for group or super group
		if msg.Chat.IsGroup() || msg.Chat.IsSuperGroup() {

			// if chat exist in config and enabled save new message
			if CONFIGS.ExistAndEnable(msg.Chat.ID) {
				NewMessage(msg.Chat.ID, msg.MessageID, msg.Date)
				log.Printf("New message %d handled for chat %d", msg.MessageID, msg.Chat.ID)
			} else {
				log.Printf(
					"Message %d not saved. Saving is disabled or there is no configuration for chat %d",
					msg.MessageID, msg.Chat.ID)
			}

			// process command message
			if msg.IsCommand() {
				cmdChan <- msg
			}
			continue
		}

		// process help and start command to bot
		if msg.IsCommand() {
			cmd := strings.ToLower(msg.Command())

			if strings.Contains("start help ping", cmd) {
				log.Printf("Bot %s command hanling", cmd)
				cmdChan <- msg
			}
		}
	}
}

// send and save reply message
func replyTo(chatID int64, msgID int, msgText string) *tgbotapi.Message {
	newMsg := tgbotapi.NewMessage(chatID, msgText)
	newMsg.ReplyToMessageID = msgID
	replyMsg, err := BOT.Send(newMsg)
	if err != nil {
		log.Println("Error occurred with sending reply:", err)
		return nil
	}

	log.Printf("Reply to %d sent", msgID)

	// save reply message
	if CONFIGS.ExistAndEnable(chatID) {
		NewMessage(chatID, replyMsg.MessageID, replyMsg.Date)
		log.Println("Reply message saved to Redis")
	}

	return &replyMsg
}

// bot command handler
func botCommandHandler(cmdChan chan *tgbotapi.Message) {
	log.Println("Start command handler")

	for msg := range cmdChan {
		command := strings.ToLower(msg.Command())
		log.Printf("Receive <%s> command from chat %d", command, msg.Chat.ID)

		switch command {
		case "help":
			replyTo(msg.Chat.ID, msg.MessageID, HelpMsg)
		case "start":
			replyTo(msg.Chat.ID, msg.MessageID, StartMsg)
		case "on":
			if CONFIGS.Exist(msg.Chat.ID) {
				// if saving is disabled
				if !CONFIGS[msg.Chat.ID].Enabled && CONFIGS[msg.Chat.ID].ChangeStatus(true) {

					// save /on command message
					NewMessage(msg.Chat.ID, msg.MessageID, msg.Date)

					replyTo(msg.Chat.ID, msg.MessageID, "Enabled saving messages")
					log.Printf("Enable saved message for chat %d", msg.Chat.ID)
				} else {
					replyTo(msg.Chat.ID, msg.MessageID, "Saving message already enabled")
					log.Printf("Saved message already enabled for chat %d", msg.Chat.ID)
				}
				// create new configuration
			} else {
				CONFIGS[msg.Chat.ID] = NewChatConfig(msg.Chat.ID, 3600, msg.Chat.Title)
				log.Printf("Create new configuration for chat %s", CONFIGS[msg.Chat.ID])

				// save /on command message
				NewMessage(msg.Chat.ID, msg.MessageID, msg.Date)

				replyTo(msg.Chat.ID, msg.MessageID,
					"Create new configuration, default message timeout 1 hour")
			}
		case "off":
			if CONFIGS.ExistAndEnable(msg.Chat.ID) {
				var replyMsg *tgbotapi.Message

				if CONFIGS[msg.Chat.ID].ChangeStatus(false) {
					replyMsg = replyTo(msg.Chat.ID, msg.MessageID, "Disabled saving messages")
					log.Printf("Disable saved message for chat %s", CONFIGS[msg.Chat.ID])
				} else {
					replyMsg = replyTo(msg.Chat.ID, msg.MessageID, "Saving message already disabled")
					log.Printf("Saved message already disabled for chat %s", CONFIGS[msg.Chat.ID])
				}

				// save reply message
				NewMessage(replyMsg.Chat.ID, replyMsg.MessageID, replyMsg.Date)
			}
		case "timeout":
			if CONFIGS.Exist(msg.Chat.ID) {
				newTime, err := time.ParseDuration(msg.CommandArguments())
				if err != nil {
					log.Printf("WARNING: Invalid new timeout value: %s", err)
					replyTo(msg.Chat.ID, msg.MessageID,
						"Error! Invalid new timeout value. Send a /help command to get help")
					break
				}
				if err := CONFIGS[msg.Chat.ID].ChangeTimeout(int(newTime.Seconds())); err != nil {
					replyMsg := fmt.Sprintf("Unable to set timeout! %s", err)
					log.Printf("WARNING: %s", replyMsg)
					replyTo(msg.Chat.ID, msg.MessageID, replyMsg)
					break
				}
				log.Printf("New timeout %s for chat %s", newTime, CONFIGS[msg.Chat.ID])
				replyTo(msg.Chat.ID, msg.MessageID, "Timeout changed")
			}
		case "delete":
			if CONFIGS.Exist(msg.Chat.ID) {
				CONFIGS[msg.Chat.ID].DeleteAllChatMessages()
			}
		case "setting":
			if CONFIGS.Exist(msg.Chat.ID) {
				timeSec := fmt.Sprintf("%ds", CONFIGS[msg.Chat.ID].Timeout)
				timeHuman, _ := time.ParseDuration(timeSec)

				status := "enable"
				if !CONFIGS[msg.Chat.ID].Enabled {
					status = "disable"
				}

				setting := fmt.Sprintf("Status: %s, Timeout: %s", status, timeHuman)
				replyTo(msg.Chat.ID, msg.MessageID, setting)
			}
		case "stop":
			if CONFIGS.Exist(msg.Chat.ID) {
				chatConfig := CONFIGS[msg.Chat.ID]

				// delete all saved message
				chatConfig.DeleteAllChatMessages()
				log.Printf("All chat %d messages have been deleted.", chatConfig.ChatID)

				chatConfig.DeleteConfig()
				delete(CONFIGS, msg.Chat.ID)
				log.Println("Chat configuration have been deleted.")

				replyTo(msg.Chat.ID, msg.MessageID, "Good by!")
			}
		case "ping":
			replyTo(msg.Chat.ID, msg.MessageID, "pong")
		default:
			log.Println("Receive unknown command:", command)
			replyTo(msg.Chat.ID, msg.MessageID,
				"Unknown command. Please send 'help' for all possible commands.")
		}
	}
}

// garbage collector for deleting older messages
func garbageCollectorHandler(timeout time.Duration) {
	log.Println("Start garbage collector handler")
	for true {
		log.Println("Garbage collector awake")
		if len(CONFIGS) > 0 {
			for _, message := range GetAllMessages(CONFIGS) {
				if message.IsOutdated() {
					message.Delete()
				}
			}
		} else {
			log.Println("No chat configurations")
		}
		log.Printf("Garbage collector falls asleep for %d seconds...", timeout)
		time.Sleep(timeout * time.Second)
	}
}
