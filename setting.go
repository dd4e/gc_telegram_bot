package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/proxy"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const HelpMsg = `
Supported command:
/help 		-- print this message
/on   		-- the bot will delete outdated messages
/off		-- the bot will be disabled
/timeout	-- new timeout after which the messages will be deleted
/delete		-- delete all messages
/setting	-- print current settings
/stop		-- !!! Delete all messages, delete settings and stop the bot !!!

Timeout format:
Timeout is set in the format: <decimal><unit suffix>
unit suffix one of "s", "m", "h"
Example: 1h15m, 24h, 30m, 60s, 10h30m15s
`

const StartMsg = `
Instructions to get started:
1. Add bot to group
2. Give him admin rights to delete messages
3. In the group send a command to the bot /on to run
4. You can change the timeout setting
5. To get the list command send /help
`

// bot setting type
type botSetting struct {
	botToken        string
	botDebug        bool
	gcTimeout       time.Duration
	dbRedisAddress  string
	dbRedisDB       int
	dbRedisPassword string
	useSocksProxy   bool
	socksParams     struct {
		socksAddress  string
		socksUser     string
		socksPassword string
	}
	timeoutLimit int
	// todo:
	//useHTTPSProxy bool
	//httpsParams struct{
	//	httpsAddress string
	//	httpsUser string
	//	httpsPassword string
	//}
}

func (s botSetting) String() string {
	return fmt.Sprint("botDebug:", s.botDebug,
		", gcTimeout:", int(s.gcTimeout),
		", dbRedisAddress:", s.dbRedisAddress,
		", dbRedisDB:", s.dbRedisDB,
		", useSocksProxy:", s.useSocksProxy,
		", socksAddress:", s.socksParams.socksAddress,
		", socksUser:", s.socksParams.socksUser,
		", timeoutLimit:", s.timeoutLimit)
}

// parsing and create setting
func parseSetting(rawData map[string]string) *botSetting {
	var setting botSetting

	// if token not set
	if _, ok := rawData["gc_token"]; !ok {
		log.Fatal("Bot token not set!")
	} else {
		setting.botToken = rawData["gc_token"]
	}

	for key, value := range rawData {
		switch key {
		case "gc_bot_debug":
			debug, err := strconv.ParseBool(value)
			if err != nil {
				log.Fatal("Bot debug must be boolean")
			}
			setting.botDebug = debug
		case "gc_check_timeout":
			timeoutInt, err := strconv.Atoi(value)
			if err != nil {
				log.Fatal("Invalid garbage collector timeout")
			}
			setting.gcTimeout = time.Duration(timeoutInt)
		case "gc_redis_addr":
			setting.dbRedisAddress = value
		case "gc_redis_db":
			db, err := strconv.Atoi(value)
			if err != nil {
				log.Fatal("Invalid redis DB number")
			}
			setting.dbRedisDB = db
		case "gc_redis_pwd":
			setting.dbRedisPassword = value
		case "gc_use_socks5":
			useSOCKS5Bool, err := strconv.ParseBool(value)
			if err != nil {
				log.Fatal("Use socks5 must be boolean")
			}
			setting.useSocksProxy = useSOCKS5Bool
		case "gc_socks5_user":
			setting.socksParams.socksUser = value
		case "gc_socks5_pwd":
			setting.socksParams.socksPassword = value
		case "gc_socks5_addr":
			setting.socksParams.socksAddress = value
		case "gc_timeout_limit":
			timeout, err := strconv.Atoi(value)
			if err != nil {
				log.Fatal("Invalid timeout limit")
			}
			setting.timeoutLimit = timeout
		}
	}

	// check socks5 settings
	if setting.useSocksProxy {
		for _, socksParam := range []string{setting.socksParams.socksAddress,
			setting.socksParams.socksPassword,
			setting.socksParams.socksUser} {
			if len(socksParam) == 0 {
				log.Fatal("Not enough parameters to configure SOCKS5 proxy")
			}
		}
	}

	// setup default gc timeout
	if setting.gcTimeout == 0 {
		setting.gcTimeout = 60
	}

	// set default timeout limit
	if setting.timeoutLimit == 0 {
		setting.timeoutLimit = 604800
	}

	// setup default redis
	if len(setting.dbRedisAddress) == 0 {
		setting.dbRedisAddress = "127.0.0.1:6379"
	}

	return &setting
}

// loading setting from system env
func loadSettingFromEnv() map[string]string {
	settingMap := make(map[string]string)
	for _, item := range os.Environ() {
		setting := strings.Split(item, "=")
		settingMap[strings.ToLower(setting[0])] = setting[1]
	}
	return settingMap
}

// load setting from config file
func loadSettingFromFile(file string) map[string]string {
	var rawJSON map[string]string

	jsonConfig, err := os.Open(file)
	if err != nil {
		log.Fatal("Error occurred with loading config file", err)
	}

	defer func() {
		if err := jsonConfig.Close(); err != nil {
			log.Print("Error occurred with closing config file", err)
		}
	}()

	decode := json.NewDecoder(jsonConfig)
	if err := decode.Decode(&rawJSON); err != nil {
		log.Fatal("Error occurred with decode config file", err)
	}

	settingMap := make(map[string]string, len(rawJSON))

	for key, value := range rawJSON {
		newKey := strings.ToLower(key)
		settingMap[newKey] = value
	}

	return settingMap
}

// todo:
//func httpsProxyClient() {}

// creates http.Client connection through SOCKS5 proxy
func socksProxyClient(address, user, password string) *http.Client {
	socksAuth := proxy.Auth{User: user, Password: password}
	dialSocksProxy, err := proxy.SOCKS5(
		"tcp",
		address,
		&socksAuth,
		proxy.Direct,
	)

	if err != nil {
		fmt.Println("Error connecting to proxy:", err)
	}

	// create client
	socksClient := &http.Client{Transport: &http.Transport{Dial: dialSocksProxy.Dial}}
	return socksClient
}
