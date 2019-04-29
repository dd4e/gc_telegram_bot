package data

import "testing"

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
	testChat.ExportToDB()
}

func TestDBKey(t *testing.T) {
	testChat.DBKey()
}
