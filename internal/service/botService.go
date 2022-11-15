package service

import (
	"github.com/ssharifzoda/bot/internal/database"
)

type BotService struct {
	db database.Bot
}

func NewBotService(db database.Bot) *BotService {
	return &BotService{db: db}
}

func (b *BotService) RegisterUser(userId int64) (string, error) {
	return b.db.RegisterUser(userId)
}
func (b *BotService) RegisterUsernames(userid int64, userName string) (string, error) {

	return b.db.RegisterUsernames(userid, userName)
}
func (b *BotService) RegisterPassword(userid int64, password string) (string, error) {
	pass := Hash(password)
	return b.db.RegisterPassword(userid, pass)
}
func (b *BotService) UserValidate(userParams []string) (string, error) {
	return b.db.UserValidate(userParams)
}
