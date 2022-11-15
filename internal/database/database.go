package database

import (
	"github.com/ssharifzoda/bot/internal/types"
	"gorm.io/gorm"
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

type Database struct {
	Bot
	Mail
}

func NewDatabase(conn *gorm.DB) *Database {
	return &Database{
		Bot:  NewBotPostgres(conn),
		Mail: NewMailPostgres(conn),
	}
}
