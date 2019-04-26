package main

import (
	"flag"
	"fmt"
	"github.com/dadmoscow/gc_telegram_bot/app"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	VERSION string
	SOURCE  *string
)

func init() {
	VERSION = "0.2"

	version := flag.Bool("version", false, "Print version")
	command := flag.Bool("manual", false, "Print help")
	SOURCE = flag.String("setting", "", "Use setting from config file")
	flag.Parse()

	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	} else if *command {
		fmt.Print(app.StartMsg, app.HelpMsg)
		os.Exit(0)
	}
}

func main() {
	log.Printf("*** Garbage Collector BotAPI. Version: %s ***", VERSION)

	setting := app.BotSetting{}
	GCBot := app.BotApp{}

	// load SETTING
	if len(*SOURCE) > 0 {
		log.Print("Loading setting from file ", *SOURCE)
		setting.Load(app.LoadSettingFromFile(*SOURCE))
	} else {
		log.Print("Loading setting from system environment")
		setting.Load(app.LoadSettingFromEnv())
	}

	log.Printf("BotAPI setting: %s", setting)

	GCBot.Init(setting)

	// chan for BOT command handler
	cmdChan := make(chan *tgbotapi.Message, 50)
	// chan for signal handler
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go GCBot.BotUpdateMsgHandler(cmdChan)
	go GCBot.BotCommandHandler(cmdChan)
	go GCBot.GarbageCollectorHandler()

	for true {
		sig := <-signals
		log.Println("Catch signal", sig)
		GCBot.Bot.BotAPI.StopReceivingUpdates()
		close(cmdChan)
		log.Println("BotAPI exit")
		os.Exit(0)
	}
}
