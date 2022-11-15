package botSystem

import (
	"errors"
	"github.com/ssharifzoda/bot/internal/service"
	"github.com/ssharifzoda/bot/pkg/logging"
)

func RegisterUser(userId int64, s *service.Service) (string, error) {
	log := logging.GetLogger()
	text, err := s.Bot.RegisterUser(userId)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return text, nil
}
func RegisterUsernames(userId int64, username string, s *service.Service) string {
	log := logging.GetLogger()
	if err := Validate(username); err != nil {
		return err.Error()
	}
	text, err := s.Bot.RegisterUsernames(userId, username)
	if err != nil {
		log.Println(err)
		return err.Error()
	}
	return text
}
func RegisterPassword(userId int64, password string, s *service.Service) string {
	log := logging.GetLogger()
	text, err := s.Bot.RegisterPassword(userId, password)
	if err != nil {
		log.Println(err)
		return err.Error()
	}
	return text
}

func Validate(msg string) error {
	count := 0
	textRune := []rune(msg)
	for _, i := range textRune {
		switch i {
		case '@':
			count++
		case '.':
			count++
		}
	}
	if count < 2 {
		return errors.New("incorrect username")
	}
	return nil

}
