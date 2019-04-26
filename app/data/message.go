package data

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Message struct {
	ChatConfig *Chat `json:"-"`
	ChatID     int64
	MsgID      int
	TimeStamp  int
}

func (msg Message) ExportToDB() (string, []byte) {
	jsonMessage, _ := json.Marshal(msg)
	key := fmt.Sprintf("msg_%d_%d", msg.ChatID, msg.MsgID)
	return key, jsonMessage
}

func (msg Message) DBKey() string {
	return fmt.Sprintf("msg_%d_%d", msg.ChatID, msg.MsgID)
}

// aging test message method
func (msg Message) IsOutdated() bool {
	delta := int(time.Now().Unix()) - msg.TimeStamp
	if delta >= msg.ChatConfig.Timeout {
		return true
	}
	return false
}

func (msg Message) String() string {
	return strconv.Itoa(msg.MsgID)
}
