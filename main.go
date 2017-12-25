package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/mdh67899/mail-provider/config"
)

func process_signal(pid int, fn func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	log.Println(pid, "register signal notify")
	for {
		s := <-sigs
		log.Println("recv", s)

		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			log.Println("gracefull shut down")
			fn()
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

	service := config.NewProgram()
	service.Init(*cfg)
	service.Start()

	go process_signal(os.Getpid(), service.Stop)

	select {}
}
