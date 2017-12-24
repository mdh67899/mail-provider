package sender

import (
	"log"

	"github.com/mdh67899/mail-provider/mail"
	"github.com/mdh67899/mail-provider/model"
	"gopkg.in/gomail.v2"
)

func (this *Store) flush(cfg *model.Config, mailer *mail.MailSender) {
	size := this.Len()
	if size == 0 {
		return
	}

	for i := 0; i < size; i++ {
		item := this.PopBack()

		if item == nil {
			continue
		}

		m := gomail.NewMessage()
		m.SetHeader("From", cfg.Auth.UserName)
		m.SetHeader("To", item.Tos)
		m.SetHeader("Subject", item.Subject)
		m.SetBody("text/plain", item.Content)

		mailer.Msg <- m

		log.Println("Process message successfull: ", item.String())
	}
}

func (this *Store) Job(cfg *model.Config, mailer *mail.MailSender) {
	for {
		select {
		case <-this.Cron.Ticker.C:
			this.flush(cfg, mailer)

		case <-this.Cron.Quit:
			this.flush(cfg, mailer)
			this.Cron.Quit_Done <- struct{}{}
			return
		}
	}
}

func (this *Store) StartCron(cfg *model.Config, mailer *mail.MailSender) {
	go this.Job(cfg, mailer)
	log.Println("Job start successfull...")
}

func (this *Store) StopCron() {
	this.Cron.Quit <- struct{}{}
	<-this.Cron.Quit_Done

	log.Println("Job stop successfull...")
}
