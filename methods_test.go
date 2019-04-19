package main

import (
	"fmt"
	"testing"
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

func TestCreateNewChatConfig(t *testing.T) {
	var testChatID int64
	testChatID = 12345

	testConfigs := GetChatConfigs()
	cfg := NewChatConfig(testChatID, 60, "test chat config")

	testConfigs[cfg.ChatID] = cfg

	if !testConfigs.Exist(testChatID) {
		t.Error("Error with creating chat configuration")
	}
}

func TestAddNewMessage(t *testing.T) {
	NewMessage(12345, 1234567, 987654325454)
	_, _ = DB.LoadFromDB(fmt.Sprintf("msg_%d_%d", 12345, 1234567))
}
