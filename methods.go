package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
	"time"
)

// bot message type
type tMessage struct {
	chatConfig *tChatConfig
	ChatID     int64
	MsgID      int
	TimeStamp  int
}

// method delete message from Redis and telegram
func (msg tMessage) Delete() {
	// delete from telegram
	delMsg := tgbotapi.DeleteMessageConfig{
		MessageID: msg.MsgID,
		ChatID:    msg.ChatID,
	}
	resp, err := BOT.DeleteMessage(delMsg)

	if err != nil {
		switch resp.ErrorCode {
		case 400:
			log.Printf("Warning: %s from chat %s", resp.Description, msg.chatConfig)
		default:
			log.Printf("Error: %s. The message %s will be deleted later from chat %s.",
				err, msg, msg.chatConfig)
			return
		}
	}

	log.Printf("The message %s has been deleted from chat %s",
		msg, msg.chatConfig)

	// delete from Redis
	key := fmt.Sprintf("msg_%d_%d", msg.ChatID, msg.MsgID)
	if err := DB.DeleteFromDB(key); err != nil {
		log.Printf("Error: message %s has not been deleted from redis: %s", msg, err)
	}
}

// aging test message method
func (msg tMessage) IsOutdated() bool {
	delta := int(time.Now().Unix()) - msg.TimeStamp
	if delta >= msg.chatConfig.Timeout {
		return true
	}
	return false
}

// method saving message to Redis
func (msg tMessage) Save() bool {
	jsonMessage, _ := json.Marshal(msg)
	key := fmt.Sprintf("msg_%d_%d", msg.ChatID, msg.MsgID)
	if err := DB.SaveToDB(key, jsonMessage); err != nil {
		log.Println("Failed to save message:", err)
		return false
	}
	return true
}

func (msg tMessage) String() string {
	return strconv.Itoa(msg.MsgID)
}

// chat configuration type
type tChatConfig struct {
	ChatID    int64
	Timeout   int
	ChatTitle string
	Enabled   bool
}

func (cnf tChatConfig) String() string {
	return cnf.ChatTitle
}

// method saving configuration to Redis
func (cnf tChatConfig) Save() bool {
	jsonConfig, _ := json.Marshal(cnf)
	key := fmt.Sprintf("chat_%d", cnf.ChatID)
	if err := DB.SaveToDB(key, jsonConfig); err != nil {
		log.Println("Failed to save message:", err)
		return false
	}
	return true
}

// change method garbage collector timeout in configuration
func (cnf *tChatConfig) ChangeTimeout(timeout int) error {
	if timeout > 0 && timeout <= SETTING.timeoutLimit {
		cnf.Timeout = timeout
		cnf.Save()
		return nil
	}

	if timeout <= 0 {
		return errors.New("timeout must be greater than 0")
	} else if timeout >= SETTING.timeoutLimit {
		maxTimeHuman, _ := time.ParseDuration(fmt.Sprintf("%ds", SETTING.timeoutLimit))
		return errors.New(fmt.Sprintf("maximum timeout value: %s", maxTimeHuman))
	}
	log.Printf("ERROR: unknown timeout error, raw timeout: %d", timeout)
	return errors.New("unknown timeout error")
}

// enable\disable saving message method
func (cnf *tChatConfig) ChangeStatus(enabled bool) bool {
	cnf.Enabled = enabled
	return cnf.Save()
}

// method checking enable status
func (cnf tChatConfig) IsEnabled() bool {
	return cnf.Enabled
}

// delete configuration method
func (cnf tChatConfig) DeleteConfig() bool {
	key := fmt.Sprintf("chat_%d", cnf.ChatID)
	err := DB.DeleteFromDB(key)
	if err != nil {
		log.Printf("Failed to delete chat %s configuration: %s", cnf, err)
		return false
	}
	return true
}

// method getting all message for chat
func (cnf tChatConfig) GetAllChatMessage() []tMessage {
	key := fmt.Sprintf("msg_%d_*", cnf.ChatID)
	jsonMessages, err := DB.LoadFromDB(key)
	if err != nil {
		log.Printf("Error occurred with loading all chat %s messages: %s", cnf, err)
		return make([]tMessage, 0)
	}

	chatMessages := make([]tMessage, len(jsonMessages))
	for n, item := range jsonMessages {
		var message tMessage
		if err := json.Unmarshal([]byte(item), &message); err != nil {
			log.Println("Error occurred with unmarshal message:", err)
			continue
		}
		if message.ChatID != cnf.ChatID {
			log.Printf("Warning! The message %s does not belong to the chat %s. Skip!",
				message, cnf)
			continue
		}
		// add chat configuration to message
		message.chatConfig = &cnf
		chatMessages[n] = message
	}
	return chatMessages
}

// method deleting all chat messages
func (cnf tChatConfig) DeleteAllChatMessages() {
	for _, message := range cnf.GetAllChatMessage() {
		message.Delete()
	}
}

// all chat configuration type
type Configs map[int64]*tChatConfig

// chat check method
func (c Configs) Exist(chatID int64) bool {
	if _, ok := c[chatID]; ok {
		return true
	}
	return false
}

// chat config exist and enable
func (c Configs) ExistAndEnable(chatID int64) bool {
	if config, ok := c[chatID]; ok {
		return config.Enabled
	}
	return false
}

// get all messages for all chats
func GetAllMessages(configs Configs) []tMessage {
	jsonMessages, err := DB.LoadFromDB("msg_*")
	if err != nil {
		log.Println("Error occurred with loading all messages:", err)
		return make([]tMessage, 0)
	}
	allMessages := make([]tMessage, len(jsonMessages))
	for n, item := range jsonMessages {
		var message tMessage
		if err := json.Unmarshal([]byte(item), &message); err != nil {
			log.Println("Error occurred with unmarshal message:", err)
			continue
		}
		// if message chat exist
		if !configs.Exist(message.ChatID) {
			log.Printf("Chat %d not found for message %s. Skip", message.ChatID, message)
			continue
		}
		// set chat configuration in to message object
		message.chatConfig = configs[message.ChatID]
		allMessages[n] = message
	}
	return allMessages
}

// create and save new message
func NewMessage(chatID int64, msgID, timestamp int) {
	newMsg := tMessage{
		ChatID:    chatID,
		MsgID:     msgID,
		TimeStamp: timestamp,
	}
	if !newMsg.Save() {
		log.Printf("Message %d from chat %d don't save", msgID, chatID)
	}
}

// create and save new configuration
func NewChatConfig(chatID int64, timeout int, title string) *tChatConfig {
	newConfig := tChatConfig{
		ChatTitle: title,
		Timeout:   timeout,
		Enabled:   true,
		ChatID:    chatID,
	}

	if !newConfig.Save() {
		log.Printf("New chat configuration fro chat %d (%s) don't save", chatID, title)
	}
	return &newConfig
}

// get all chat configuration
func GetChatConfigs() Configs {
	chatConfigs := make(Configs)
	values, err := DB.LoadFromDB("chat_*")
	if err != nil {
		log.Println("Error occurred with loading chat configurations", err)
		return chatConfigs
	}

	for _, item := range values {
		var config tChatConfig

		err := json.Unmarshal([]byte(item), &config)
		if err != nil {
			log.Println("Error occurred with unmarshal configuration:", err)
			continue
		}
		chatConfigs[config.ChatID] = &config
	}
	return chatConfigs
}
