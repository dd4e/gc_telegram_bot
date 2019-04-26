package test

import (
	"gc_telegram_bot"
	"os"
	"testing"
)

func TestParseSettingFromEnv(t *testing.T) {

	testData := map[string]string{
		"gc_token":         "TOKEN",
		"gc_bot_debug":     "false",
		"gc_check_timeout": "60",
		"gc_redis_addr":    "127.0.0.1:6379",
		"gc_redis_db":      "0",
		"gc_redis_pwd":     "qwerty",
		"gc_use_socks5":    "true",
		"gc_socks5_user":   "user",
		"gc_socks5_pwd":    "password",
		"gc_socks5_addr":   "172.0.0.1:8080",
	}

	for key, value := range testData {
		_ = os.Setenv(key, value)
	}

	testSetting := botSetting{
		botToken:        "TOKEN",
		botDebug:        false,
		gcTimeout:       60,
		dbRedisAddress:  "127.0.0.1:6379",
		dbRedisDB:       0,
		dbRedisPassword: "qwerty",
		useSocksProxy:   true,
		socksParams: struct {
			socksAddress  string
			socksUser     string
			socksPassword string
		}{socksAddress: "172.0.0.1:8080",
			socksUser:     "user",
			socksPassword: "password"},
		timeoutLimit: 604800,
	}

	setting := parseSetting(loadSettingFromEnv())

	if *setting != testSetting {
		t.Error("Error with compare setting")
	}
}

func TestParseFromFile(t *testing.T) {
	testSetting := botSetting{
		botToken:        "BOT_TOKEN",
		botDebug:        true,
		gcTimeout:       60,
		dbRedisAddress:  "127.0.0.1:6379",
		dbRedisDB:       0,
		dbRedisPassword: "",
		useSocksProxy:   false,
		timeoutLimit:    604800,
	}

	setting := parseSetting(loadSettingFromFile("config.json"))

	if *setting != testSetting {
		t.Error("Error with compare setting from file")
	}
}
