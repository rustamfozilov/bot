package mail

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message"
	_ "github.com/emersion/go-message/charset"
	"github.com/ssharifzoda/bot/internal/service"
	"github.com/ssharifzoda/bot/internal/types"
	"io"
	"io/ioutil"
	"log"
)

const refreshTimeMail = "refreshmail"

func UpdateMsgCounts(userID, unseenMsg, totalMsg int, s *service.Service) error {
	err := s.Mail.UpdateCounts(userID, unseenMsg, totalMsg)
	if err != nil {
		return err
	}
	return nil
}

func Conductor(user *types.Users) (*client.Client, error) {
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

func GetBodyMassage(mbox *imap.MailboxStatus, c *client.Client) types.Response {
	var r types.Response
	seqSet := new(imap.SeqSet)
	seqSet.AddRange(mbox.Messages, mbox.Messages)
	section := &imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem()}
	messages := make(chan *imap.Message, 1)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqSet, items, messages)
	}()
	msg := <-messages
	bodyParams := msg.GetBody(section)
	m, err := message.Read(bodyParams)
	if err != nil {
		log.Println(err)
	}
	r.From = m.Header.Get("From")
	mr := m.MultipartReader()
	defer mr.Close()
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		contentType, _, _ := p.Header.ContentType()
		if contentType == "text/plain" {
			body, err := ioutil.ReadAll(p.Body)
			r.Body = string(body)
			if err != nil {
				log.Println(err)
			}
		}
	}
	return r
}
