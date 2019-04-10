package main

import (
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	DB      *redis.Client
	BOT     *tgbotapi.BotAPI
	CONFIGS Configs
	VERSION string
	SETTING *botSetting
)

func init() {
	VERSION = "0.1"

	version := flag.Bool("version", false, "Print version")
	command := flag.Bool("manual", false, "Print bot manual")
	flag.Parse()

	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	} else if *command {
		fmt.Print(StartMsg, HelpMsg)
		os.Exit(0)
	}
}

func main() {
	log.Printf("*** Garbage Collector Bot. Version: %s ***", VERSION)
	// load SETTING
	SETTING = parseSetting(loadSettingFromEnv())
	log.Printf("Bot setting: %s", SETTING)

	DB = redis.NewClient(&redis.Options{
		DB:       SETTING.dbRedisDB,
		Addr:     SETTING.dbRedisAddress,
		Password: SETTING.dbRedisPassword,
	})

	// try ping redis
	_, err := DB.Ping().Result()
	if err != nil {
		log.Fatal("Error occurred with connect to Redis:", err)
	}

	CONFIGS = GetChatConfigs()
	log.Println("Loading configurations:", len(CONFIGS))

	// chan for BOT command handler
	cmdChan := make(chan *tgbotapi.Message, 50)
	// chan for signal handler
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	if SETTING.useSocksProxy {
		socks := SETTING.socksParams
		socksClient := socksProxyClient(socks.socksAddress, socks.socksUser, socks.socksPassword)
		BOT, err = tgbotapi.NewBotAPIWithClient(SETTING.botToken, socksClient)
	} else {
		BOT, err = tgbotapi.NewBotAPI(SETTING.botToken)
	}

	if err != nil {
		log.Fatal("Connection error to bot API telegram:", err)
	}

	log.Printf("Authorized on account %s", BOT.Self.UserName)

	BOT.Debug = SETTING.botDebug

	go botUpdateMsgHandler(cmdChan)
	go botCommandHandler(cmdChan)
	go garbageCollectorHandler(SETTING.gcTimeout)

	for true {
		sig := <-signals
		log.Println("Catch signal", sig)
		BOT.StopReceivingUpdates()
		close(cmdChan)
		log.Println("Bot exit")
		os.Exit(0)
	}
}
