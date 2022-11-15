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
		log.Error(err)
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
	countMsg--
	to := mbox.Messages
	// We're using unsigned integers here, only subtract if the result is > 0
	if countMsg == 0 {
		return nil
	}
	from = mbox.Messages - uint32(countMsg)
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
		defer func(mr message.MultipartReader) {
			err := mr.Close()
			if err != nil {
				log.Println(err)
			}
		}(mr)
		for {
			p, err := mr.NextPart()
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
			}
		}
		responses = append(responses, r)
	}
	return responses
}

func CountMsgAnalyze(user *types.Users, unseenMsg, totalMsg int, s *service.Service) int {
	log := logging.GetLogger()
	//Первый раз когда регистрируются данные в базу данных
	if user.UnseenMsgCount == 0 && user.TotalMsgCount == 0 {
		if err := UpdateMsgCounts(user.UserId, unseenMsg, totalMsg, s); err != nil {
			log.Print(err)
		}
		return 0
	}
	if totalMsg < user.TotalMsgCount {
		if err := UpdateMsgCounts(user.UserId, unseenMsg, totalMsg, s); err != nil {
			log.Print(err)
		}
		return 0
	}
	//Пользователь прочитал какое - то количество сообщений у себя
	if unseenMsg < user.UnseenMsgCount && totalMsg == user.TotalMsgCount {
		if err := UpdateMsgCounts(user.UserId, unseenMsg, totalMsg, s); err != nil {
			log.Print(err)
		}
		return 0
	}
	// Новое сообщение
	if unseenMsg > user.UnseenMsgCount {
		if err := UpdateMsgCounts(user.UserId, user.UnseenMsgCount, totalMsg, s); err != nil {
			log.Print(err)
		}
		return unseenMsg - user.UnseenMsgCount
	}
	//Пользователь решил прочитать непрочитанные сообщения и вдруг ему приходит новое сообщение и новое тоже прочтет
	if unseenMsg < user.UnseenMsgCount && totalMsg > user.TotalMsgCount {
		if err := UpdateMsgCounts(user.UserId, unseenMsg, totalMsg, s); err != nil {
			log.Print(err)
		}
		return 0
	}
	//Момент, когда пользователю пришли новые сообщения, он зашёл, прочитал некоторые из них,
	// но, некоторое количество так и оставил непрочитанным и сразу же приходит ещё одно новое сообщение
	if unseenMsg > user.UnseenMsgCount && user.TotalMsgCount > totalMsg {
		totalDiff := totalMsg - user.TotalMsgCount
		unseenDiff := unseenMsg - user.UnseenMsgCount
		if totalDiff > unseenDiff {
			if err := UpdateMsgCounts(user.UserId, unseenMsg-totalDiff, totalMsg, s); err != nil {
				log.Print(err)
			}
			return totalDiff
		}
	}
	return 0
}
