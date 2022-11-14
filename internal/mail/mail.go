package mail

import (
	"github.com/emersion/go-imap"
	"github.com/spf13/viper"
	"github.com/ssharifzoda/bot/internal/service"
	"github.com/ssharifzoda/bot/internal/types"
	"log"
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
		log.Println(err)
	}
	for _, user := range users {
		resp := Connector(user)
		//Первый раз когда регистрируются данные в базу данных
		if user.UnseenMsgCount == 0 && user.TotalMsgCount == 0 {
			if err := UpdateMsgCounts(resp.UserId, resp.UnseenMsgCount, resp.TotalMsgCount, s); err != nil {
				log.Print(err)
			}
			continue
		}
		// Нет новых сообщений
		if user.UnseenMsgCount == resp.UnseenMsgCount && user.TotalMsgCount == resp.TotalMsgCount {
			continue
		}
		//Пользователь прочитал какое - то количество сообщений у себя
		if resp.UnseenMsgCount < user.UnseenMsgCount && resp.TotalMsgCount == user.TotalMsgCount {
			if err := UpdateMsgCounts(resp.UserId, resp.UnseenMsgCount, resp.TotalMsgCount, s); err != nil {
				log.Print(err)
			}
			continue
		}
		// Новое сообщение
		if resp.UnseenMsgCount > user.UnseenMsgCount {
			if err := UpdateMsgCounts(resp.UserId, resp.UnseenMsgCount, resp.TotalMsgCount, s); err != nil {
				log.Print(err)
			}
			//Отправляем количество непрочитанных - новых
			resp.UnseenMsgCount = resp.UnseenMsgCount - user.UnseenMsgCount
			ch <- resp
			continue
		}
		//Момент, когда пользователю пришли новые сообщения, он зашёл, прочитал некоторые из них,
		// но, некоторое количество так и оставил непрочитанным и сразу же приходит ещё одно новое сообщение
		if resp.UnseenMsgCount > user.UnseenMsgCount && user.TotalMsgCount > resp.TotalMsgCount {
			totalDiff := resp.TotalMsgCount - user.TotalMsgCount
			unseenDiff := resp.UnseenMsgCount - user.TotalMsgCount
			if totalDiff > unseenDiff {
				if err := UpdateMsgCounts(resp.UserId, resp.UnseenMsgCount, resp.TotalMsgCount, s); err != nil {
					log.Print(err)
				}
			}
			resp.UnseenMsgCount = totalDiff
			ch <- resp
			continue
		}
		//Пользователь решил прочитать непрочитанные сообщения и вдруг ему приходит новое сообщение и новое тоже прочтет
		if resp.UnseenMsgCount < user.UnseenMsgCount && resp.TotalMsgCount > user.TotalMsgCount {
			if err := UpdateMsgCounts(resp.UserId, resp.UnseenMsgCount, resp.

func Connector(user *types.Users) types.Response {
	if user.MailLogin == "" || user.MailPassword == "" || user.UserId == 0 || user.MailService == "" {
		return types.Response{}
	}
	c, err := Conductor(user)
	if err != nil {
		log.Println(err)
		return types.Response{}
	}
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Println(err)
	}
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{"\\Seen"}
	unseenMsg, err := c.Search(criteria)
	if err != nil {
		log.Println(err)
	}
	response := GetBodyMassage(mbox, c)
	response.UnseenMsgCount = len(unseenMsg)
	response.TotalMsgCount = int(mbox.Messages)
	response.UserId = user.UserId
	return response
}
