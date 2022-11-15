package service

import (
	"github.com/ssharifzoda/bot/internal/database"
	"github.com/ssharifzoda/bot/internal/types"
)

type Bot interface {
	RegisterUser(userId int64) (string, error)
	RegisterUsernames(userid int64, userName string) (string, error)
	RegisterPassword(userid int64, password string) (string, error)
}

type Mail interface {
	UpdateCounts(userId, unseenMsg, totalMsg int) error
	GetAllUsers() ([]*types.Users, error)
}

type Service struct {
	Bot
	Mail
}

func NewService(db *database.Database) *Service {
	return &Service{
		Bot:  NewBotService(db),
		Mail: NewMailService(db),
	}
}
