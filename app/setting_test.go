package app

import (
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

	testSetting := BotSetting{
		BotToken:      "TOKEN",
		BotDebug:      false,
		SleepTimeout:  60,
		RedisAddress:  "127.0.0.1:6379",
		RedisDB:       0,
		RedisPassword: "qwerty",
		UseSocksProxy: true,
		SocksParams: struct {
			SocksAddress  string
			SocksUser     string
			SocksPassword string
		}{
			SocksAddress:  "172.0.0.1:8080",
			SocksUser:     "user",
			SocksPassword: "password",
		},
		TimeoutLimit: 604800,
	}

	for key, value := range testData {
		_ = os.Setenv(key, value)
	}

	setting := BotSetting{}
	setting.Load(LoadSettingFromEnv())

	if setting != testSetting {
		t.Error("Error with compare setting")
	}
}

func TestParseFromFile(t *testing.T) {
	testSetting := BotSetting{
		BotToken:      "BOT_TOKEN",
		BotDebug:      true,
		SleepTimeout:  60,
		RedisAddress:  "127.0.0.1:6379",
		RedisDB:       0,
		RedisPassword: "",
		UseSocksProxy: false,
		TimeoutLimit:  604800,
	}

	setting := BotSetting{}
	setting.Load(LoadSettingFromFile("../config/config.json"))

	if setting != testSetting {
		t.Error("Error with compare setting from file")
	}
}
