package data

import (
	"bytes"
	"encoding/json"
	"testing"
)

var testChat = Chat{
	ChatID:       12345,
	TimeoutLimit: 604800,
	ChatTitle:    "TestChat",
	Enabled:      true,
	Timeout:      3600,
}

func TestChat_ChangeStatus(t *testing.T) {
	testChat.ChangeStatus(false)
	if testChat.IsEnabled() {
		t.Error("change status error")
	}
}

func TestChat_ChangeTimeout(t *testing.T) {
	_ = testChat.ChangeTimeout(180)
	if testChat.Timeout != 180 {
		t.Error("error change timeout")
	}
}

func TestChat_ExportToDB(t *testing.T) {
	key, value := testChat.ExportToDB()
	if key != "chat_12345" {
		t.Error("export db error: key")
	}

	testJSON, _ := json.Marshal(testChat)
	if !bytes.Equal(value, testJSON) {
		t.Error("export db error: value")
	}
}

func TestChat_DBKey(t *testing.T) {
	if testChat.DBKey() != "chat_12345" {
		t.Error("export db key error")
	}
}
