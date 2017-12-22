package sender

import (
	"log"
	"time"

	"github.com/mdh67899/go-utils/cron"
	"github.com/mdh67899/mail-provider/config"
	"github.com/mdh67899/mail-provider/mail"
	"gopkg.in/gomail.v2"
)

var QueueCron = cron.NewCronScheduler(time.Second * 3)

func Job(Queue *safeLinklist, cron *cron.CronScheduler, mailer *mail.MailSender) {
	for {
		select {
		case <-cron.Ticker.C:
			size := Queue.Len()
			if size == 0 {
				continue
			}

			for i := 0; i < size; i++ {
				item := Queue.PopBack()

				if item == nil {
					continue
				}

				m := gomail.NewMessage()
				m.SetHeader("From", config.Cfg.Auth.UserName)
				m.SetHeader("To", item.Tos)
				m.SetHeader("Subject", item.Subject)
				m.SetBody("text/plain", item.Content)

				mailer.Msg <- m

				log.Println("Process message successfull: ", item.String())
			}

		case <-cron.Quit:
			cron.Quit_Done <- struct{}{}
			return
		}
	}
}

func StartCron() {
	go Job(Queue, QueueCron, mail.Mailer)
	log.Println("Job start successfull...")
}

func StopCron() {
	QueueCron.Quit <- struct{}{}
	<-QueueCron.Quit_Done

	QueueCron.Destory()
	log.Println("Job stop successfull...")
}
