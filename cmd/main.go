package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/ssharifzoda/bot/internal/botSystem"
	"github.com/ssharifzoda/bot/internal/database"
	"github.com/ssharifzoda/bot/internal/database/postgres"
	"github.com/ssharifzoda/bot/internal/mail"
	"github.com/ssharifzoda/bot/internal/service"
	"github.com/ssharifzoda/bot/internal/types"
	"log"
	"os"
)

// Канал, куда будем помещать ответы к запросу почти
var mailResponseChan chan types.Response

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error initializing env value: %s", err.Error())
	}
	conn, _ := postgres.NewPostgresGorm()
	db := database.NewDatabase(conn)
	src := service.NewService(db)
	mailResponseChan = make(chan types.Response)
	go mail.GetNewMails(src, mailResponseChan)
	RunBot(src, mailResponseChan)
}

func initConfig() error {
	viper.AddConfigPath("internal/config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func RunBot(s *service.Service, ch chan types.Response) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	var ucfg = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	updates := bot.GetUpdatesChan(ucfg)
	for {
		select {
		case update := <-updates:
			botSystem.BotCommandHandler(update, bot)
			botSystem.BotLogicCommands(bot, update, s)
		case r := <-ch:
			response := fmt.Sprintf("У вас новых сообщений: %d. \n Последнее сообщение от: %s. Текст: %s",
				r.UnseenMsgCount, r.From, r.Body)
			msg := tgbotapi.NewMessage(int64(r.UserId), response)
			bot.Send(msg)
		}
	}
}
