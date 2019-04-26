package test

import (
	"fmt"
	"testing"
	"time"
)

func TestDBConnection(t *testing.T) {
	DB = redisDB{
		Addr:  "127.0.0.1:6379",
		numDB: 1,
	}

	if err := DB.Connect(); err != nil {
		t.Error("Failed to connect to redis:", err)
	}
}

func TestNewChatConfig(t *testing.T) {
	testChatID := int64(12345)

	testConfigs := GetChatConfigs()
	cfg := NewChatConfig(testChatID, 60, "test chat config")

	testConfigs[cfg.ChatID] = cfg

	if !testConfigs.Exist(testChatID) {
		t.Error("Error with creating chat configuration")
	}
}

func TestNewMessage(t *testing.T) {
	testChatID := int64(12345)
	testMsgID := 1234567
	NewMessage(testChatID, testMsgID, int(time.Now().Unix()))
	_, err := DB.LoadFromDB(fmt.Sprintf("msg_%d_%d", testChatID, testMsgID))
	if err != nil {
		t.Error("Error with load message")
	}
}

func TestDeleteMessage(t *testing.T) {
	//testChatID := int64(12345)
	chatConfigs := GetChatConfigs()

	for _, value := range chatConfigs {
		value.DeleteAllChatMessages()
	}
}
