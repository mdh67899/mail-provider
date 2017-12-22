package mail

import (
	"errors"
	"log"
	"time"

	"github.com/mdh67899/go-utils/cron"
	"github.com/mdh67899/mail-provider/model"
	"gopkg.in/gomail.v2"
)

func NewAuthUser(host string, port int, username string, password string) *model.AuthUser {
	return &model.AuthUser{
		Host:     host,
		Port:     port,
		UserName: username,
		PassWord: password,
	}
}

type MailSender struct {
	Cfg       *model.AuthUser
	Dialer    *gomail.Dialer
	S         gomail.SendCloser
	Ready     bool
	Msg       chan *gomail.Message
	Scheduler *cron.JobScheduler
}

func InitMailSender(cfg model.AuthUser) {
	if !cfg.Islegal() {
		log.Fatalln("[ERROR] authenticate user is not legal:", cfg.String())
	}

	Mailer = &MailSender{
		Cfg:       &cfg,
		Dialer:    gomail.NewDialer(cfg.Host, cfg.Port, cfg.UserName, cfg.PassWord),
		Msg:       make(chan *gomail.Message),
		Scheduler: cron.NewJobScheduler(),
	}
}

func (this *MailSender) Open() {
	var err error
	if this.Ready {
		return
	}

	if this.S, err = this.Dialer.Dial(); err != nil {
		log.Println("[ERROR] Dialer to smtp server failed:", err)
		return
	}

	this.Ready = true
}

func (this *MailSender) Close() {
	if !this.Ready {
		return
	}

	if err := this.S.Close(); err != nil {
		this.Ready = false
		log.Println("[ERROR] Close dialer smtp server connnection failed:", err)
		return
	}

	this.Ready = false
}

func (this *MailSender) Send(m *gomail.Message) error {
	for i := 0; i < 3; i++ {
		if this.Ready {
			break
		}

		this.Open()

		time.Sleep(time.Millisecond * 100)
	}

	if !this.Ready {
		return errors.New("dialer to smtp server failed")
	}

	if err := gomail.Send(this.S, m); err != nil {
		this.Close()
		return err
	}

	return nil
}

var (
	Mailer *MailSender
)

func StartCron() {
	go Mailer.Process()
	log.Println("Worker start successfull...")
}

func StopCron() {
	Mailer.Scheduler.Quit <- struct{}{}
	<-Mailer.Scheduler.Quit_Done
	Mailer.Scheduler.Destory()

	log.Println("Worker stop successfull...")
}

func (this *MailSender) Process() {
	for {
		select {
		case m, ok := <-this.Msg:
			if !ok {
				return
			}

			err := this.Send(m)

			if err != nil {
				log.Println("[ERROR] send email failed, error is: ", err)
			}

			log.Println("send message to smtp server successfull")

		// Close the connection to the SMTP server if no email was sent in
		// the last 30 seconds.
		case <-time.After(30 * time.Second):
			this.Close()

		case <-this.Scheduler.Quit:
			this.Close()
			this.Scheduler.Quit_Done <- struct{}{}
			return
		}
	}
}
