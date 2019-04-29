package data

import (
	"fmt"
	"testing"
)

var testChat = Chat{
	ChatID:       12345,
	TimeoutLimit: 604800,
	ChatTitle:    "TestChat",
	Enabled:      true,
	Timeout:      3600,
}

func TestChangeStatus(t *testing.T) {
	testChat.ChangeStatus(false)
	if testChat.IsEnabled() {
		t.Error("change status error")
	}
}

func TestChangeTimeout(t *testing.T) {
	_ = testChat.ChangeTimeout(180)
	if testChat.Timeout != 180 {
		t.Error("error change timeout")
	}
}

func TestExportToDB(t *testing.T) {
	key, value := testChat.ExportToDB()
	if len(key) == 0 && len(value) == 0 {
		t.Error("error with export to db")
	}
}

func TestDBKey(t *testing.T) {
	if testChat.DBKey() != fmt.Sprintf("chat_%d", testChat.ChatID) {
		t.Error("db key error")
	}
}
