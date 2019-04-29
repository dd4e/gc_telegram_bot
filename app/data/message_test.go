package data

import (
	"fmt"
	"testing"
	"time"
)

var testMessage = Message{
	ChatID:    12345,
	TimeStamp: int(time.Now().Unix()),
	MsgID:     123,
	ChatConfig: &Chat{
		ChatID:       12345,
		TimeoutLimit: 100500,
		ChatTitle:    "testChat",
		Enabled:      true,
		Timeout:      0,
	},
}

func TestMessage_IsOutdated(t *testing.T) {
	if !testMessage.IsOutdated() {
		t.Error("error with outdated func")
	}
}

func TestMessage_DBKey(t *testing.T) {
	key := fmt.Sprintf("msg_%d_%d", testMessage.ChatID, testMessage.MsgID)
	if testMessage.DBKey() != key {
		t.Error("error db key")
	}
}

func TestMessage_ExportToDB(t *testing.T) {
	key, value := testMessage.ExportToDB()
	if len(key) == 0 && len(value) == 0 {
		t.Error("error with export to db")
	}
}
