package mail

import (
	"github.com/emersion/go-imap/client"
	"github.com/ssharifzoda/bot/internal/service"
	"github.com/ssharifzoda/bot/internal/types"
	"github.com/ssharifzoda/bot/pkg/logging"
)

func Gmail(user *types.Users) (*client.Client, error) {
	log := logging.GetLogger()
	c, err := client.DialTLS(gmailAddress, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	pass, _ := service.DeHash(user.MailPassword)
	err = c.Login(user.MailLogin, pass)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return c, nil
}
func MailRu(user *types.Users) (*client.Client, error) {
	log := logging.GetLogger()
	c, err := client.DialTLS(mailAddress, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	login := user.MailLogin + "@" + user.MailService
	pass, _ := service.DeHash(user.MailPassword)
	err = c.Login(login, pass)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return c, nil
}

func Icloud(user *types.Users) (*client.Client, error) {
	log := logging.GetLogger()
	c, err := client.DialTLS(idcloudAddress, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	pass, _ := service.DeHash(user.MailPassword)
	err = c.Login(user.MailLogin, pass)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return c, nil
}

func HumoTj(user *types.Users) (*client.Client, error) {
	c, err := client.Dial(humoAddress)
	if err != nil {
		return nil, err
	}
	pass, err := service.DeHash(user.MailPassword)
	if err != nil {
		return nil, err
	}
	err = c.Login(user.MailLogin, pass)
	if err != nil {
		return nil, err
	}
	return c, nil
}
