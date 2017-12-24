package config

import (
	"github.com/mdh67899/mail-provider/http"
	"github.com/mdh67899/mail-provider/mail"
	"github.com/mdh67899/mail-provider/sender"
	"github.com/mdh67899/mail-provider/model"
	"sync"
)

type Program struct {
	sync.RWMutex

	queue  *sender.Store
	cfg    *model.Config
	server *http.HttpServer
	mailer *mail.MailSender
}

func NewProgram() *Program {
	return &Program{queue: nil, cfg: nil, server: nil, mailer: nil}
}

func (this *Program) Init(cfg string) {
	this.cfg = ParseConfig(cfg)
	this.queue = sender.NewStore()
	this.mailer = mail.InitMailSender(this.cfg)
	this.server = http.NewhttpServer()
	this.server.Init(this.cfg, this.queue)
}

func (this *Program) Start() {
	this.mailer.StartCron()
	this.queue.StartCron(this.cfg, this.mailer)
	this.server.StartServer()
}

func (this *Program) Stop() {
	this.server.StopServer()
	this.queue.StopCron()
	this.mailer.StopCron()
}
