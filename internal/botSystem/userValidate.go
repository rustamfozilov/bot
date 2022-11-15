package botSystem

import (
	"github.com/ssharifzoda/bot/internal/service"
	"github.com/ssharifzoda/bot/pkg/logging"
)

func UserValidate(userParams []string, s *service.Service) (string, error) {
	log := logging.GetLogger()
	text, err := s.Bot.UserValidate(userParams)
	if err != nil {
		log.Println(err)
	}
	return text, err
}
