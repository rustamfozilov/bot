package mail

import (
	"github.com/emersion/go-imap"
	_ "github.com/emersion/go-imap"
	"github.com/spf13/viper"
	"github.com/ssharifzoda/bot/internal/service"
	"github.com/ssharifzoda/bot/internal/types"
	"time"
)

func GetNewMails(s *service.Service, m chan types.Response) {
	tick := time.NewTicker(viper.GetDuration(refreshTimeMail))
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			EMailSystem(m, s)
		}
	}
}
func EMailSystem(ch chan types.Response, s *service.Service) {
	users, err := s.Mail.GetAllUsers()
	if err != nil {
		log.Error(err)
	}
	for _, user := range users {
		resp := Connector(user, s)
		if resp == nil {
			continue
		}
		for _, i := range resp {
			i.UserId = user.UserId
			ch <- i
			continue
		}
	}
}

func Connector(user *types.Users, s *service.Service) []types.Response {
	if user.MailLogin == "" || user.MailPassword == "" || user.UserId == 0 || user.MailService == "" {
		return nil
	}
	c, err := Conductor(user)
	if err != nil {
		log.Error(err)
		return nil
	}
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Error(err)
	}
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{"\\Seen"}
	unseenMsg, err := c.Search(criteria)
	if err != nil {
		log.Error(err)
	}
	count := CountMsgAnalyze(user, len(unseenMsg), int(mbox.Messages), s)
	if count == 0 {
		return nil
	}
	responses := GetBodyMassage(mbox, c, count)
	return responses
}
