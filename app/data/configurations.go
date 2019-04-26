package data

import (
	"errors"
)

// all chat configuration type
type Configs map[int64]*Chat

func (c Configs) Get(chatID int64) (*Chat, error) {
	if c.Exist(chatID) {
		return c[chatID], nil
	}

	return nil, errors.New("unknown chat configuration")
}

// chat check method
func (c Configs) Exist(chatID int64) bool {
	if _, ok := c[chatID]; ok {
		return true
	}
	return false
}

// chat config exist and enable
func (c Configs) ExistAndEnable(chatID int64) bool {
	if config, ok := c[chatID]; ok {
		return config.Enabled
	}
	return false
}
