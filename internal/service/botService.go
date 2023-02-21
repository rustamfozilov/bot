package service

import (
	"errors"
	"github.com/ssharifzoda/bot/internal/database"
)

type BotService struct {
	db database.Bot
}

func NewBotService(db database.Bot) *BotService {
	return &BotService{db: db}
}

func (b *BotService) RegisterUser(userId int64) (string, error) {
	check, err := b.db.CheckUser(userId)
	if err != nil {
		return "", err
	}
	if check == true {
		return "", errors.New("эти данные уже есть в системе")
	}
	return b.db.RegisterUser(userId)
}
func (b *BotService) RegisterUsernames(userid int64, userName string) (string, error) {
	return b.db.RegisterUsernames(userid, userName)
}
func (b *BotService) RegisterPassword(userid int64, password string) (string, error) {
	pass, err := Hash(password)
	if err != nil {
		return "", err
	}
	return b.db.RegisterPassword(userid, pass)
}
