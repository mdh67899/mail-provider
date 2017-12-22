package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/mdh67899/mail-provider/config"
	"github.com/mdh67899/mail-provider/http"
	"github.com/mdh67899/mail-provider/mail"
	"github.com/mdh67899/mail-provider/sender"
)

func process_signal(pid int) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	log.Println(pid, "register signal notify")
	for {
		s := <-sigs
		log.Println("recv", s)

		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			log.Println("gracefull shut down")
			http.HttpServer.Exit()
			sender.StopCron()
			mail.StopCron()

			log.Println(pid, "exit")
			os.Exit(0)
		}
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	cfg := flag.String("c", "cfg.yaml", "configuration file")
	flag.Parse()

	config.ParseConfig(*cfg)

	mail.InitMailSender(config.Cfg.Auth)
	http.HttpServer.Init(config.Cfg.Addr)

	mail.StartCron()
	sender.StartCron()

	go http.HttpServer.Start()
	process_signal(os.Getpid())
}
