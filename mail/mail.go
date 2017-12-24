package mail

import (
	"errors"
	"log"
	"time"

	"github.com/mdh67899/go-utils/cron"
	"github.com/mdh67899/mail-provider/model"
	"gopkg.in/gomail.v2"
	"sync"
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
	sync.RWMutex

	Cfg       *model.AuthUser
	Dialer    *gomail.Dialer
	S         gomail.SendCloser
	Msg       chan *gomail.Message
	Scheduler *cron.JobScheduler
	Ready     bool
}

func InitMailSender(cfg *model.Config) *MailSender {
	if !cfg.Auth.Islegal() {
		log.Fatalln("[ERROR] authenticate user is not legal:", cfg.Auth.String())
	}

	Mailer := &MailSender{
		Cfg:       &cfg.Auth,
		Dialer:    gomail.NewDialer(cfg.Auth.Host, cfg.Auth.Port, cfg.Auth.UserName, cfg.Auth.PassWord),
		S:         nil,
		Msg:       make(chan *gomail.Message),
		Scheduler: cron.NewJobScheduler(),
		Ready:     false,
	}

	return Mailer
}

func (this *MailSender) open() {
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

func (this *MailSender) close() {
	if !this.Ready {
		return
	}

	if err := this.S.Close(); err != nil {
		log.Println("[ERROR] Close dialer smtp server connnection failed:", err)
		return
	}

	this.Ready = false
}

func (this *MailSender) Send(m *gomail.Message) error {
	this.Lock()
	defer this.Unlock()

	for i := 0; i < 3; i++ {
		if this.Ready {
			break
		}

		this.open()

		time.Sleep(time.Millisecond * 100)
	}

	if !this.Ready {
		return errors.New("dialer to smtp server failed")
	}

	if err := gomail.Send(this.S, m); err != nil {
		this.close()
		return err
	}

	return nil
}

func (this *MailSender) StartCron() {
	go this.Process()
	log.Println("Worker start successfull...")
}

func (this *MailSender) StopCron() {
	this.Scheduler.Quit <- struct{}{}
	<-this.Scheduler.Quit_Done

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
			this.close()

		case <-this.Scheduler.Quit:
			this.close()
			this.Scheduler.Quit_Done <- struct{}{}
			return
		}
	}
}
