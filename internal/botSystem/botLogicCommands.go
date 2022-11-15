package botSystem

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ssharifzoda/bot/internal/service"
	"github.com/ssharifzoda/bot/pkg/logging"
	"strings"
)

func BotLogicCommands(bot *tgbotapi.BotAPI, update tgbotapi.Update, s *service.Service) {
	log := logging.GetLogger()
	userId := update.Message.From.ID
	text := update.Message.Text
	repl := strings.ReplaceAll(text, " ", "")
	sliceStr := strings.Split(repl, "-")
	for _, val := range sliceStr {
		switch val {
		case "Салом":
			msg, err := RegisterUser(userId, s)
			if err != nil {
				log.Println(err)
			}
			t := tgbotapi.NewMessage(userId, msg)
			if _, err := bot.Send(t); err != nil {
				log.Println(err)
			}
		case "password":
			res := RegisterPassword(userId, sliceStr[1], s)
			t := tgbotapi.NewMessage(userId, res)
			if _, err := bot.Send(t); err != nil {
				log.Println(err)
			}
		case "login":
			userId = update.Message.From.ID
			res := RegisterUsernames(userId, sliceStr[1], s)
			t := tgbotapi.NewMessage(userId, res)
			if _, err := bot.Send(t); err != nil {
				log.Println(err)
			}
		}
	}
}
