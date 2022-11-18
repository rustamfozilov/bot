package mail

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message"
	_ "github.com/emersion/go-message/charset"
	"github.com/ssharifzoda/bot/internal/service"
	"github.com/ssharifzoda/bot/internal/types"
	"github.com/ssharifzoda/bot/pkg/logging"
	_ "github.com/ssharifzoda/bot/pkg/logging"
	"io"
	"io/ioutil"
)

const refreshTimeMail = "refreshmail"

func UpdateMsgCounts(userID, unseenMsg, totalMsg int, s *service.Service) error {
	log := logging.GetLogger()
	err := s.Mail.UpdateCounts(userID, unseenMsg, totalMsg)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func Conductor(user *types.Users) (*client.Client, error) {
	log := logging.GetLogger()
	switch user.MailService {
	case "gmail.com":
		c, err := Gmail(user)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return c, nil
	case "mail.ru":
		c, err := MailRu(user)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return c, nil
	case "inbox.ru":
		c, err := MailRu(user)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return c, nil
	case "icloud.com":
		c, err := Icloud(user)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return c, nil
	}
	return nil, nil
}

func GetBodyMassage(mbox *imap.MailboxStatus, c *client.Client, countMsg int) []types.Response {
	log := logging.GetLogger()
	var responses []types.Response
	var r types.Response
	var from uint32
	to := mbox.Messages
	from = mbox.Messages - uint32(countMsg-1)
	seqSet := new(imap.SeqSet)
	seqSet.AddRange(from, to)
	section := &imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem()}
	messages := make(chan *imap.Message, countMsg)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqSet, items, messages)
	}()
	for msg := range messages {
		bodyParams := msg.GetBody(section)
		m, err := message.Read(bodyParams)
		if err != nil {
			log.Println(err)
		}
		r.From = m.Header.Get("From")
		mr := m.MultipartReader()
		if mr == nil {
			return nil
		}
		defer func(mr message.MultipartReader) {
			err := mr.Close()
			if err != nil {
				log.Println(err)
				return
			}
		}(mr)
		for {
			p, err := mr.NextPart()
			defer func() {
				if err := recover(); err != nil {
					log.Println("panic occurred: ", err)
				}
			}() //Вероятное место возникновение паники, поэтому сделан перехват паники
			if err == io.EOF {
				break
			}
			contentType, _, _ := p.Header.ContentType()
			if contentType == "text/plain" {
				body, err := ioutil.ReadAll(p.Body)
				if err != nil {
					log.Println(err)
				}
				r.Body = string(body)
				responses = append(responses, r)
			}
		}
	}
	return responses
}

func CountMsgAnalyze(user *types.Users, unseenMsg, totalMsg int, s *service.Service) int {
	log := logging.GetLogger()
	//Первое обновление почты пользователя
	if user.UnseenMsgCount == 0 && user.TotalMsgCount == 0 {
		if err := UpdateMsgCounts(user.UserId, unseenMsg, totalMsg, s); err != nil {
			log.Println(err)
		}
		return 0
	}
	// У пользователя новое/новые письма
	if unseenMsg > user.UnseenMsgCount && totalMsg > user.TotalMsgCount {
		if err := UpdateMsgCounts(user.UserId, user.UnseenMsgCount, totalMsg, s); err != nil {
			log.Println(err)
		}
		return unseenMsg - user.UnseenMsgCount
	}
	if err := UpdateMsgCounts(user.UserId, unseenMsg, totalMsg, s); err != nil {
		log.Println(err)
	}
	return 0
}
