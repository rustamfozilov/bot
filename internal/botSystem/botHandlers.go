package botSystem

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ssharifzoda/bot/pkg/logging"
)

func BotCommandHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	log := logging.GetLogger()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	switch update.Message.Command() {
	case "command1":
		msg.Text = "Напишите мне приветствие для регистрации: Салом"
	case "command2":
		msg.Text = "Отправьте свой логин от почты в следующем виде\n" +
			"Пример: login - test@gmail.com"
	case "command3":
		msg.Text = "Отправьте свой пароль от почты в следующем виде\n" +
			"Пример: password - 6552856sc"
	}
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}
