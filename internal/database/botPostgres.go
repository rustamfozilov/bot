package database

import (
	"fmt"
	"github.com/ssharifzoda/bot/internal/types"
	"github.com/ssharifzoda/bot/pkg/logging"
	"gorm.io/gorm"
	"strings"
)

type BotPostgres struct {
	conn *gorm.DB
}

func NewBotPostgres(conn *gorm.DB) *BotPostgres {
	return &BotPostgres{conn: conn}
}

func (b *BotPostgres) RegisterUser(userId int64) (string, error) {
	query := fmt.Sprintf("insert into test_bot_users (user_id) values(?);")
	if err := b.conn.Exec(query, userId); err.Error != nil {
		return "", err.Error
	}
	msg := "Салом. Первый шаг успешно пройден"
	return msg, nil
}
func (b *BotPostgres) RegisterUsernames(userid int64, userName string) (string, error) {
	log := logging.GetLogger()
	s := strings.Split(userName, "@")
	if err := b.conn.Table("test_bot_users").Where("user_id = ?", userid).Update("mail_login", s[0]); err.Error != nil {
		log.Println(err)
		return "", err.Error
	}
	if err := b.conn.Table("test_bot_users").Where("user_id = ?", userid).Update("mail_service", s[1]); err.Error != nil {
		log.Println(err)
		return "", err.Error
	}
	msg := "Логин, успешно сохранён"
	return msg, nil
}
func (b *BotPostgres) RegisterPassword(userid int64, password string) (string, error) {
	log := logging.GetLogger()
	query := fmt.Sprintf("update test_bot_users set mail_password = ? where user_id = ?;")
	var username string
	if err := b.conn.Raw(query, password, userid).Scan(&username); err.Error != nil {
		log.Println(err)
		return "", err.Error
	}
	msg := "Поздравляю. Вы успешно прошли регистрацию. Теперь все ваши письма я вам сюда отправлю"
	return msg, nil
}

func (b *BotPostgres) CheckUser(userId int64) (bool, error) {
	var user *types.Users
	if db := b.conn.Where("user_id = ?", userId).First(&user); db.Error != nil {
		return true, db.Error
	}
	if user.UserId != 0 {
		return true, nil
	}
	return false, nil
}
