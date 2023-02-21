package database

import (
	"github.com/ssharifzoda/bot/internal/types"
	"github.com/ssharifzoda/bot/pkg/logging"
	"gorm.io/gorm"
)

type MailPostgres struct {
	conn *gorm.DB
}

func NewMailPostgres(conn *gorm.DB) *MailPostgres {
	return &MailPostgres{conn: conn}
}

func (m *MailPostgres) GetAllUsers() ([]*types.Users, error) {
	log := logging.GetLogger()
	var users []*types.Users
	row := m.conn.Table("test_bot_users").Find(&users)
	if row.Error != nil {
		log.Println(row.Error)
	}
	return users, row.Error
}
func (m *MailPostgres) UpdateCounts(userId, unseenMsg, totalMsg int) error {
	log := logging.GetLogger()
	tx := m.conn.Table("test_bot_users").Where("user_id", userId).Update("unseen_msg_count", unseenMsg)
	tx = m.conn.Table("test_bot_users").Where("user_id", userId).Update("total_msg_count", totalMsg)
	if tx.Error != nil {
		log.Println(tx.Error)
		return tx.Error
	}
	return nil
}
