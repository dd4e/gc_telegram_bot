package app

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

// bot setting type
type BotSetting struct {
	BotToken      string
	BotDebug      bool
	SleepTimeout  time.Duration
	RedisAddress  string
	RedisDB       int
	RedisPassword string
	UseSocksProxy bool
	SocksParams   struct {
		SocksAddress  string
		SocksUser     string
		SocksPassword string
	}
	TimeoutLimit int
	// todo:
	//useHTTPSProxy bool
	//httpsParams struct{
	//	httpsAddress string
	//	httpsUser string
	//	httpsPassword string
	//}
}

func (s BotSetting) String() string {
	return fmt.Sprint("BotDebug:", s.BotDebug,
		", SleepTimeout:", int(s.SleepTimeout),
		", RedisAddress:", s.RedisAddress,
		", DBHandlers:", s.RedisDB,
		", UseSocksProxy:", s.UseSocksProxy,
		", SocksAddress:", s.SocksParams.SocksAddress,
		", SocksUser:", s.SocksParams.SocksUser,
		", TimeoutLimit:", s.TimeoutLimit)
}

// parsing and create setting
func (s *BotSetting) Load(rawData map[string]string) {
	// if token not set
	if _, ok := rawData["gc_token"]; !ok {
		log.Fatal("BotAPI token not set!")
	} else {
		s.BotToken = rawData["gc_token"]
	}

	for key, value := range rawData {
		switch key {
		case "gc_bot_debug":
			debug, err := strconv.ParseBool(value)
			if err != nil {
				log.Fatal("BotAPI debug must be boolean")
			}
			s.BotDebug = debug
		case "gc_sleep_timeout":
			timeoutInt, err := strconv.Atoi(value)
			if err != nil {
				log.Fatal("Invalid garbage collector timeout")
			}
			s.SleepTimeout = time.Duration(timeoutInt)
		case "gc_redis_addr":
			s.RedisAddress = value
		case "gc_redis_db":
			db, err := strconv.Atoi(value)
			if err != nil {
				log.Fatal("Invalid redis DB number")
			}
			s.RedisDB = db
		case "gc_redis_pwd":
			s.RedisPassword = value
		case "gc_use_socks5":
			useSOCKS5Bool, err := strconv.ParseBool(value)
			if err != nil {
				log.Fatal("Use socks5 must be boolean")
			}
			s.UseSocksProxy = useSOCKS5Bool
		case "gc_socks5_user":
			s.SocksParams.SocksUser = value
		case "gc_socks5_pwd":
			s.SocksParams.SocksPassword = value
		case "gc_socks5_addr":
			s.SocksParams.SocksAddress = value
		case "gc_timeout_limit":
			timeout, err := strconv.Atoi(value)
			if err != nil {
				log.Fatal("Invalid timeout limit")
			}
			s.TimeoutLimit = timeout
		}
	}

	// check socks5 settings
	if s.UseSocksProxy {
		for _, socksParam := range []string{s.SocksParams.SocksAddress,
			s.SocksParams.SocksPassword,
			s.SocksParams.SocksUser} {
			if len(socksParam) == 0 {
				log.Fatal("Not enough parameters to configure SOCKS5 proxy")
			}
		}
	}

	// setup default gc timeout
	if s.SleepTimeout == 0 {
		s.SleepTimeout = 60
	}

	// set default timeout limit
	if s.TimeoutLimit == 0 {
		s.TimeoutLimit = 604800
	}

	// setup default redis
	if len(s.RedisAddress) == 0 {
		s.RedisAddress = "127.0.0.1:6379"
	}
}

// loading setting from system env
func LoadSettingFromEnv() map[string]string {
	settingMap := make(map[string]string)
	for _, item := range os.Environ() {
		setting := strings.Split(item, "=")
		settingMap[strings.ToLower(setting[0])] = setting[1]
	}
	return settingMap
}

// load setting from config file
func LoadSettingFromFile(file string) map[string]string {
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
//func HTTPProxyClient() {}

// creates http.Client connection through SOCKS5 proxy
func SOCKS5ProxyClient(address, user, password string) *http.Client {
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
