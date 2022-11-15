package mail

import (
	"github.com/emersion/go-imap/client"
	"github.com/ssharifzoda/bot/internal/service"
	"github.com/ssharifzoda/bot/internal/types"
	"github.com/ssharifzoda/bot/pkg/logging"
)

var log *logging.Logger

func Gmail(user *types.Users) (*client.Client, error) {
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	pass := service.DeHash(user.MailPassword)
	err = c.Login(user.MailLogin, pass)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return c, nil
}
func MailRu(user *types.Users) (*client.Client, error) {
	c, err := client.DialTLS("imap.mail.ru:993", nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	login := user.MailLogin + "@" + user.MailService
	pass := service.DeHash(user.MailPassword)
	err = c.Login(login, pass)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return c, nil
}

func Icloud(user *types.Users) (*client.Client, error) {
	c, err := client.DialTLS("imap.mail.me.com:993", nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	pass := service.DeHash(user.MailPassword)
	err = c.Login(user.MailLogin, pass)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return c, nil
}
