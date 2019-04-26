package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
)

// chat configuration type
type Chat struct {
	ChatID       int64
	Timeout      int
	ChatTitle    string
	Enabled      bool
	TimeoutLimit int `json:"-"`
}

func (cnf Chat) String() string {
	return cnf.ChatTitle
}

func (cnf Chat) ExportToDB() (string, []byte) {
	jsonConfig, _ := json.Marshal(cnf)
	key := fmt.Sprintf("chat_%d", cnf.ChatID)
	return key, jsonConfig
}

func (cnf Chat) DBKey() string {
	return fmt.Sprintf("chat_%d", cnf.ChatID)
}

// change method garbage collector timeout in configuration
func (cnf *Chat) ChangeTimeout(timeout int) error {
	if timeout > 0 && timeout <= cnf.TimeoutLimit {
		cnf.Timeout = timeout
		return nil
	}

	if timeout <= 0 {
		return errors.New("timeout must be greater than 0")
	} else if timeout >= cnf.TimeoutLimit {
		maxTimeHuman, _ := time.ParseDuration(fmt.Sprintf("%ds", cnf.TimeoutLimit))
		return errors.New(fmt.Sprintf("maximum timeout value: %s", maxTimeHuman))
	}
	log.Printf("ERROR: unknown timeout error, raw timeout: %d", timeout)
	return errors.New("unknown timeout error")
}

// enable\disable saving message method
func (cnf *Chat) ChangeStatus(enabled bool) {
	cnf.Enabled = enabled
}

// method checking enable status
func (cnf Chat) IsEnabled() bool {
	return cnf.Enabled
}
