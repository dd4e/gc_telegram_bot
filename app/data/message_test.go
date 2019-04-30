package data

import (
	"bytes"
	"encoding/json"
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
	if testMessage.DBKey() != "msg_12345_123" {
		t.Error("error db key")
	}
}

func TestMessage_ExportToDB(t *testing.T) {
	key, value := testMessage.ExportToDB()

	if key != "msg_12345_123" {
		t.Error("export db error: key")
	}

	testJSON, _ := json.Marshal(testMessage)
	if !bytes.Equal(value, testJSON) {
		t.Error("export db error: value")
	}
}
