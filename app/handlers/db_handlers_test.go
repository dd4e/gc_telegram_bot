package handlers

import (
	"github.com/dadmoscow/gc_telegram_bot/app/data"
	"math/rand"
	"testing"
	"time"
)

var testDBInstance = DBHandlers{
	DBNum: 1,
	Addr:  "127.0.0.1:6379",
}

var testChat = data.Chat{
	Timeout:   60,
	Enabled:   true,
	ChatTitle: "testChat",
	ChatID:    12345,
}

func TestDBHandlers_Connect(t *testing.T) {
	if err := testDBInstance.Connect(); err != nil {
		t.Fatal(err)
	}
}

func TestDBHandlers_SaveData(t *testing.T) {
	if err := testDBInstance.SaveData(testChat); err != nil {
		t.Error(err)
	}

	for _, value := range rand.Perm(5) {
		testMsg := data.Message{
			ChatID:     12345,
			ChatConfig: &testChat,
			MsgID:      value,
			TimeStamp:  int(time.Now().Unix()),
		}

		if err := testDBInstance.SaveData(testMsg); err != nil {
			t.Error(err)
		}
	}
}

func TestDBHandlers_GetAllChatMessage(t *testing.T) {
	message := testDBInstance.GetAllChatMessage(testChat)
	if len(message) != 5 {
		t.Errorf("error with get all chat message: %d", len(message))
	}
}

func TestDBHandlers_GetAllConfigs(t *testing.T) {
	configs := testDBInstance.GetAllConfigs()
	config, err := configs.Get(12345)
	if err != nil {
		t.Error(err)
	}

	if *config != testChat {
		t.Error("error loading chat config")
	}
}

func TestDBHandlers_GetAllMessages(t *testing.T) {

	configs := testDBInstance.GetAllConfigs()
	allMessage := testDBInstance.GetAllMessages(configs)

	if len(allMessage) != 5 {
		t.Error("error loading all message")
	}
}

func TestDBHandlers_LoadData(t *testing.T) {
	_, err := testDBInstance.LoadData(testChat)
	if err != nil {
		t.Error(err)
	}
}

func TestDBHandlers_DeleteData(t *testing.T) {
	chatMessages := testDBInstance.GetAllChatMessage(testChat)
	for _, message := range chatMessages {
		err := testDBInstance.DeleteData(message)
		if err != nil {
			t.Error(err)
		}
	}

	err := testDBInstance.DeleteData(testChat)
	if err != nil {
		t.Error(err)
	}
}
