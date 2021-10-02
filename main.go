package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/arriqaaq/kubediff/config"
	"github.com/arriqaaq/kubediff/pkg/watcher"
)

var (
	configPath = flag.String("config", "", "config file path")
)

func main() {
	flag.Parse()
	conf, err := config.New(*configPath)
	if err != nil {
		log.Fatalf("Error in loading configuration. Error:%s", err.Error())
	}
	watcher, err := watcher.NewWatcher(conf)
	if err != nil {
		log.Fatalf("Error in loading configuration. Error:%s", err.Error())
	}

	stopCh := make(chan struct{})
	defer close(stopCh)

	watcher.Run(stopCh)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGSTOP)

	<-sigterm
}
