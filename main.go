package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/arriqaaq/kubediff/config"
	"github.com/arriqaaq/kubediff/pkg/watcher"
)

var (
	configFilePath = "config.yaml"
	configPath     = flag.String("config", "", "config folder path")
)

func main() {
	flag.Parse()
	filepath := filepath.Join(*configPath, configFilePath)
	conf, err := config.New(filepath)
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
