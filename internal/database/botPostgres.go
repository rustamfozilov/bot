package database

import (
	"fmt"
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
	log := logging.GetLogger()
	query := fmt.Sprintf("insert into my_users (user_id) values(?);")
	if err := b.conn.Exec(query, userId); err.Error != nil {
		log.Println(err)
		return "", err.Error
	}
	msg := "Вы успешно прошли регистрацию"
	return msg, nil
}
func (b *BotPostgres) RegisterUsernames(userid int64, userName string) (string, error) {
	log := logging.GetLogger()
	s := strings.Split(userName, "@")
	if err := b.conn.Table("my_users").Where("user_id = ?", userid).Update("mail_login", s[0]); err.Error != nil {
		log.Println(err)
		return "", err.Error
	}
	if err := b.conn.Table("my_users").Where("user_id = ?", userid).Update("mail_service", s[1]); err.Error != nil {
		log.Println(err)
		return "", err.Error
	}
	msg := "Гуд. Я сохранил ваш логин"
	return msg, nil
}
func (b *BotPostgres) RegisterPassword(userid int64, password string) (string, error) {
	log := logging.GetLogger()
	query := fmt.Sprintf("update my_users set mail_password = ? where user_id = ?;")
	var username string
	if err := b.conn.Raw(query, password, userid).Scan(&username); err.Error != nil {
		log.Println(err)
		return "", err.Error
	}
	msg := "Всё. Теперь я буду вас оповещать о новых письмах"
	return msg, nil
}

func (b *BotPostgres) UserValidate(userParams []string) (string, error) {
	log := logging.GetLogger()
	switch userParams[0] {
	case "Салом":
		userId := 0
		if err := b.conn.Raw("select id from my_users where user_id = ?", userParams[1]); err.Error != nil {
			log.Println(err.Error)
			err.Scan(&userId)
			if userId > 0 {
				return "You already registered", nil
			}
		}
	case "login":
		userId := 0
		if err := b.conn.Raw("select id from my_users where mail_login = ?", userParams[1]); err.Error != nil {
			log.Println(err.Error)
			err.Scan(&userId)
			if userId > 0 {
				return "You already registered", nil
			}
		}
	}
	return "", nil
}
