package service

import (
	"github.com/ssharifzoda/bot/internal/database"
	"github.com/ssharifzoda/bot/internal/types"
)

type MailService struct {
	db database.Mail
}

func NewMailService(db database.Mail) *MailService {
	return &MailService{db: db}
}

func (m *MailService) GetAllUsers() ([]*types.Users, error) {

	return m.db.GetAllUsers()
}

func (m *MailService) UpdateCounts(userId, unseenMsg, totalMsg int) error {
	return m.db.UpdateCounts(userId, unseenMsg, totalMsg)
}
