package http

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mdh67899/mail-provider/model"
	"github.com/mdh67899/mail-provider/sender"
)

type httpServer struct {
	srv       *http.Server
	exit      chan bool
	exit_done chan bool
}

func NewhttpServer() *httpServer {
	return &httpServer{srv: nil, exit: make(chan bool, 1), exit_done: make(chan bool, 1)}
}

var HttpServer = NewhttpServer()

func mail(c *gin.Context) {
	mail := &model.Mail{}
	err := c.ShouldBind(mail)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if mail.Tos == "" || mail.Subject == "" || mail.Content == "" {
		c.JSON(http.StatusOK, gin.H{"status": "tos,subject,content is empty,please check", "message": mail.String()})
	}

	c.JSON(http.StatusOK, gin.H{"message": mail.String()})

	sender.Queue.SafePush(mail)
}

func (this *httpServer) Init(addr string) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		time.Sleep(1 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	router.POST("/mail", mail)

	if addr == "" {
		log.Fatalln("[error] http listen address is: ", addr)
	}

	this.srv = &http.Server{
		Addr:    addr,
		Handler: router,
	}
}

func (this *httpServer) Start() {

	go func() {
		// service connections
		if err := this.srv.ListenAndServe(); err != nil {
			log.Println("listen: ", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 3 seconds.
	<-this.exit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := this.srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")

	this.exit_done <- true
}

func (this *httpServer) Exit() {
	if this.srv != nil {
		this.exit <- true
		<-this.exit_done
		this.srv = nil
	}
}
