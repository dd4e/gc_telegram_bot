package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/dadmoscow/gc_telegram_bot/app/data"
	"github.com/go-redis/redis"
	"log"
)

type DBHandlers struct {
	client   *redis.Client
	Addr     string
	Password string
	DBNum    int
}

type BotDB interface {
	Set(string, []byte) error
	Get(string) ([]string, error)
	Delete(string) error
	SaveData(DBMethods) error
	DeleteData(DBMethods) error
	LoadData(DBMethods) ([]string, error)
}

type DBMethods interface {
	ExportToDB() (string, []byte)
	DBKey() string
}

// initial connection to Redis
func (db *DBHandlers) Connect() error {
	db.client = redis.NewClient(
		&redis.Options{
			DB:       db.DBNum,
			Addr:     db.Addr,
			Password: db.Password,
		},
	)

	// try ping redis
	if err := db.client.Ping().Err(); err != nil {
		return err
	}
	return nil
}

func (db DBHandlers) SaveData(d DBMethods) error {
	return db.set(d.ExportToDB())
}

func (db DBHandlers) DeleteData(d DBMethods) error {
	return db.delete(d.DBKey())
}

func (db DBHandlers) LoadData(d DBMethods) ([]string, error) {
	return db.get(d.DBKey())
}

// save key value to Redis
func (db DBHandlers) set(key string, value []byte) error {
	err := db.client.Set(key, value, 0).Err()
	if err != nil {
		log.Println("Error occurred with save to Redis:", err)
		return err
	}
	return nil
}

// load data from Redis by filtered key
func (db DBHandlers) get(filter string) ([]string, error) {
	keys, err := db.client.Keys(filter).Result()
	if err != nil {
		log.Println("Error occurred with loading data from Redis:", err)
		return nil, err
	}

	values := make([]string, len(keys))

	for i, key := range keys {
		value, err := db.client.Get(key).Result()
		if err != nil {
			log.Printf(
				"Error occurred with getting value by key %s: %s. Skip...", key, err)
			continue
		}
		values[i] = value
	}
	return values, nil
}

// delete value by key from Redis
func (db DBHandlers) delete(key string) error {
	err := db.client.Del(key).Err()
	if err != nil {
		log.Printf("Error occurred with deleting value by key %s: %s", key, err)
		return err
	}
	return nil
}

// get all chat configuration
func (db DBHandlers) GetAllConfigs() data.Configs {
	chatConfigs := make(data.Configs)

	values, err := db.get("chat_*")
	if err != nil {
		log.Println("Error occurred with loading chat configurations", err)
		return chatConfigs
	}

	for _, item := range values {
		var config data.Chat

		err := json.Unmarshal([]byte(item), &config)
		if err != nil {
			log.Println("Error occurred with unmarshal configuration:", err)
			continue
		}
		chatConfigs[config.ChatID] = &config
	}
	return chatConfigs
}

// method getting all message for chat
func (db DBHandlers) GetAllChatMessage(chat data.Chat) []data.Message {
	key := fmt.Sprintf("msg_%d_*", chat.ChatID)
	jsonMessages, err := db.get(key)
	if err != nil {
		log.Printf("Error occurred with loading all chat %s messages: %s", chat, err)
		return make([]data.Message, 0)
	}

	chatMessages := make([]data.Message, len(jsonMessages))
	for n, item := range jsonMessages {
		var message data.Message
		if err := json.Unmarshal([]byte(item), &message); err != nil {
			log.Println("Error occurred with unmarshal message:", err)
			continue
		}
		if message.ChatID != chat.ChatID {
			log.Printf("Warning! The message %s does not belong to the chat %s. Skip!",
				message, chat)
			continue
		}
		// add chat configuration to message
		message.ChatConfig = &chat
		chatMessages[n] = message
	}
	return chatMessages
}

// get all messages for all chats
func (db DBHandlers) GetAllMessages(configs data.Configs) []data.Message {
	jsonMessages, err := db.get("msg_*")
	if err != nil {
		log.Println("Error occurred with loading all messages:", err)
		return make([]data.Message, 0)
	}

	allMessages := make([]data.Message, len(jsonMessages))
	for n, item := range jsonMessages {
		var message data.Message
		if err := json.Unmarshal([]byte(item), &message); err != nil {
			log.Println("Error occurred with unmarshal message:", err)
			continue
		}

		// set chat configuration in to message object
		if message.ChatConfig, err = configs.Get(message.ChatID); err != nil {
			log.Printf("Chat %d not found for message %s. Skip", message.ChatID, message)
			continue
		}

		allMessages[n] = message
	}
	return allMessages
}
