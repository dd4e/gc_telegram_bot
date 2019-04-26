package app

import (
	"fmt"
	"github.com/dadmoscow/gc_telegram_bot/app/data"
	"github.com/dadmoscow/gc_telegram_bot/app/handlers"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
	"time"
)

type BotApp struct {
	Bot     handlers.TGBotHandlers
	DB      handlers.DBHandlers
	Configs data.Configs
	setting *BotSetting
}

func (b *BotApp) Init(setting *BotSetting) {
	b.setting = setting

	// init db connection
	b.DB.Addr = setting.RedisAddress
	b.DB.Password = setting.RedisPassword
	b.DB.DBNum = setting.RedisDB

	if err := b.DB.Connect(); err != nil {
		log.Fatal("Error occurred with connect to Redis:", err)
	}

	// init telegram bot api connection
	var bot *tgbotapi.BotAPI
	var err error

	if setting.UseSocksProxy {
		socks := setting.SocksParams
		socksClient := SOCKS5ProxyClient(socks.SocksAddress, socks.SocksUser, socks.SocksPassword)
		bot, err = tgbotapi.NewBotAPIWithClient(setting.BotToken, socksClient)
	} else {
		bot, err = tgbotapi.NewBotAPI(setting.BotToken)
	}

	if err != nil {
		log.Fatal("Connection error to bot API telegram:", err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)
	bot.Debug = setting.BotDebug
	b.Bot.BotAPI = bot

	// init configs
	b.Configs = b.DB.GetAllConfigs()
	log.Println("Loading configurations:", len(b.Configs))
}

func (b BotApp) DeleteMessage(message data.Message) {
	if err := b.Bot.DeleteMessage(message); err == nil {
		_ = b.DB.DeleteData(message)
	}
}

func (b BotApp) NewMessage(chatID int64, messageID int, date int) {
	newMsg := data.Message{
		ChatID:    chatID,
		MsgID:     messageID,
		TimeStamp: date,
	}

	if err := b.DB.SaveData(newMsg); err != nil {
		log.Printf("%s", err)
	}
}

// send and save reply message
func (b BotApp) replyAndSave(msg *tgbotapi.Message, msgText string) *tgbotapi.Message {
	newMsg := tgbotapi.NewMessage(msg.Chat.ID, msgText)
	newMsg.ReplyToMessageID = msg.MessageID
	replyMsg, err := b.Bot.BotAPI.Send(newMsg)
	if err != nil {
		log.Println("Error occurred with sending reply:", err)
		return nil
	}

	log.Printf("Reply to %d sent", msg.MessageID)

	// save reply message
	if b.Configs.ExistAndEnable(msg.Chat.ID) {
		b.NewMessage(msg.Chat.ID, replyMsg.MessageID, replyMsg.Date)
		log.Println("Reply message saved to Redis")
	}
	return &replyMsg
}

func (b BotApp) deleteAllChatMessage(chat data.Chat) {
	for _, message := range b.DB.GetAllChatMessage(chat) {
		b.DeleteMessage(message)
	}
}

// new message handler
func (b BotApp) BotUpdateMsgHandler(cmdChan chan *tgbotapi.Message) {
	log.Println("Start new message handler")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, _ := b.Bot.BotAPI.GetUpdatesChan(u)

	for update := range updates {
		msg := update.Message

		// skip non message updates
		if msg == nil {
			continue
		}

		// process message only for group or super group
		if msg.Chat.IsGroup() || msg.Chat.IsSuperGroup() {

			// if chat exist in config and enabled save new message
			if b.Configs.ExistAndEnable(msg.Chat.ID) {
				b.NewMessage(msg.Chat.ID, msg.MessageID, msg.Date)
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
				log.Printf("BotAPI %s command hanling", cmd)
				cmdChan <- msg
			}
		}
	}
}

// bot command handler
func (b BotApp) BotCommandHandler(cmdChan chan *tgbotapi.Message) {
	log.Println("Start command handler")

	for msg := range cmdChan {
		command := strings.ToLower(msg.Command())
		log.Printf("Receive <%s> command from chat %d", command, msg.Chat.ID)

		switch command {
		case "help":
			b.replyAndSave(msg, HelpMsg)
		case "start":
			b.replyAndSave(msg, StartMsg)
		case "on":
			if b.Configs.Exist(msg.Chat.ID) {
				config := b.Configs[msg.Chat.ID]
				// if saving is disabled
				if !config.Enabled {
					// enable and save
					config.ChangeStatus(true)
					if err := b.DB.SaveData(config); err != nil {
						log.Printf("%s", err)
					}

					// save /on command message
					b.NewMessage(msg.Chat.ID, msg.MessageID, msg.Date)

					b.replyAndSave(msg, "Enabled saving messages")
					log.Printf("Enable saved message for chat %d", msg.Chat.ID)
				} else {
					b.replyAndSave(msg, "Saving message already enabled")
					log.Printf("Saved message already enabled for chat %d", msg.Chat.ID)
				}
				// create new configuration
			} else {
				newConfig := data.Chat{
					ChatID:       msg.Chat.ID,
					Timeout:      3600,
					Enabled:      true,
					ChatTitle:    msg.Chat.Title,
					TimeoutLimit: b.setting.TimeoutLimit,
				}

				if err := b.DB.SaveData(newConfig); err != nil {
				}

				b.Configs[msg.Chat.ID] = &newConfig
				log.Printf("Create new configuration for chat %s", b.Configs[msg.Chat.ID])

				// save /on command message
				b.NewMessage(msg.Chat.ID, msg.MessageID, msg.Date)

				b.replyAndSave(msg, "Create new configuration, default message timeout 1 hour")
			}
		case "off":
			if b.Configs.ExistAndEnable(msg.Chat.ID) {
				var replyMsg *tgbotapi.Message

				b.Configs[msg.Chat.ID].ChangeStatus(false)

				if err := b.DB.SaveData(b.Configs[msg.Chat.ID]); err == nil {
					replyMsg = b.replyAndSave(msg, "Disabled saving messages")
					log.Printf("Disable saved message for chat %s", b.Configs[msg.Chat.ID])
				} else {
					replyMsg = b.replyAndSave(msg, "Saving message already disabled")
					log.Printf("Saved message already disabled for chat %s", b.Configs[msg.Chat.ID])
				}

				// save reply message
				b.NewMessage(replyMsg.Chat.ID, replyMsg.MessageID, replyMsg.Date)
			}
		case "timeout":
			if b.Configs.Exist(msg.Chat.ID) {
				newTime, err := time.ParseDuration(msg.CommandArguments())
				if err != nil {
					log.Printf("WARNING: Invalid new timeout value: %s", err)
					b.replyAndSave(msg, "Error! Invalid new timeout value. Send a /help command to get help")
					break
				}
				if err := b.Configs[msg.Chat.ID].ChangeTimeout(int(newTime.Seconds())); err != nil {
					replyMsg := fmt.Sprintf("Unable to set timeout! %s", err)
					log.Printf("WARNING: %s", replyMsg)
					b.replyAndSave(msg, replyMsg)
					break
				}

				_ = b.DB.SaveData(b.Configs[msg.Chat.ID])

				log.Printf("New timeout %s for chat %s", newTime, b.Configs[msg.Chat.ID])
				b.replyAndSave(msg, "Timeout changed")
			}
		case "delete":
			if b.Configs.Exist(msg.Chat.ID) {
				b.deleteAllChatMessage(*b.Configs[msg.Chat.ID])
			}
		case "setting":
			if b.Configs.Exist(msg.Chat.ID) {
				timeSec := fmt.Sprintf("%ds", b.Configs[msg.Chat.ID].Timeout)
				timeHuman, _ := time.ParseDuration(timeSec)

				status := "enable"
				if !b.Configs[msg.Chat.ID].Enabled {
					status = "disable"
				}

				setting := fmt.Sprintf("Status: %s, Timeout: %s", status, timeHuman)
				b.replyAndSave(msg, setting)
			}
		case "stop":
			if b.Configs.Exist(msg.Chat.ID) {
				chatConfig := b.Configs[msg.Chat.ID]

				// delete all saved message
				b.deleteAllChatMessage(*chatConfig)
				log.Printf("All chat %d messages have been deleted.", chatConfig.ChatID)

				_ = b.DB.DeleteData(chatConfig)
				delete(b.Configs, msg.Chat.ID)
				log.Println("Chat configuration have been deleted.")

				b.replyAndSave(msg, "Good by!")
			}
		case "ping":
			b.replyAndSave(msg, "pong")
		default:
			log.Println("Receive unknown command:", command)
			b.replyAndSave(msg, "Unknown command. Please send 'help' for all possible commands.")
		}
	}
}

// garbage collector for deleting older messages
func (b BotApp) GarbageCollectorHandler() {
	log.Println("Start garbage collector handler")
	for true {
		log.Println("Garbage collector awake")
		if len(b.Configs) > 0 {
			for _, message := range b.DB.GetAllMessages(b.Configs) {
				if message.IsOutdated() {
					b.DeleteMessage(message)
				}
			}
		} else {
			log.Println("No chat configurations")
		}
		log.Printf("Garbage collector falls asleep for %d seconds...", b.setting.SleepTimeout)
		time.Sleep(b.setting.SleepTimeout * time.Second)
	}
}
