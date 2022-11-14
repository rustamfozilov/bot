package botSystem

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func BotCommandHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	switch update.Message.Command() {
	case "command1":
		msg.Text = "Напишите мне приветсвие для регистрации: Салом"

	case "command2":
		msg.Text = "Отправь свой логин от почти в следующем виде\n" +
			"Пример: login - test@gmail.com"
	case "command3":
		msg.Text = "Отправь свой пароль от почти в следующем виде\n" +
			"Пример: password - 6552856sc"
	}
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}
}
