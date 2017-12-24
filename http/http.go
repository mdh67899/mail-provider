package http

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mdh67899/go-utils/cron"
	"github.com/mdh67899/mail-provider/model"
	"github.com/mdh67899/mail-provider/sender"
	"sync"
)

type HttpServer struct {
	sync.RWMutex
	srv   *http.Server
	queue *sender.Store
	job   *cron.JobScheduler
}

func NewhttpServer() *HttpServer {
	return &HttpServer{srv: nil, queue: nil, job: cron.NewJobScheduler()}
}

func (this *HttpServer) mail(c *gin.Context) {
	msg := &model.Mail{}
	err := c.ShouldBind(msg)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if msg.Tos == "" || msg.Subject == "" || msg.Content == "" {
		c.JSON(http.StatusOK, gin.H{"status": "tos,subject,content is empty,please check", "message": msg.String()})
	}

	this.queue.SafePush(msg)

	c.JSON(http.StatusOK, gin.H{"message": msg.String()})
}

func (this *HttpServer) Init(cfg *model.Config, queue *sender.Store) {
	this.queue = queue

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	router.POST("/mail", this.mail)

	if cfg.Addr == "" {
		log.Fatalln("[error] http listen address is: ", cfg.Addr)
	}

	this.srv = &http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
}

func (this *HttpServer) start() {
	this.RLock()
	defer this.RUnlock()

	go func() {
		// service connections
		if err := this.srv.ListenAndServe(); err != nil {
			log.Println("listen: ", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 3 seconds.
	<-this.job.Quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := this.srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")

	this.job.Quit_Done <- struct{}{}
}

func (this *HttpServer) StartServer() {
	go this.start()
	log.Println("Server Start succeessfull.")
}

func (this *HttpServer) StopServer() {
	this.stop()
}


func (this *HttpServer) stop() {
	this.RLock()
	defer this.RUnlock()

	if this.srv != nil {
		this.job.Quit <- struct{}{}
		<-this.job.Quit_Done
		this.srv = nil
	}
}
