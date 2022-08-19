package main

import (
	"context"
	"notifier/cmd"
	"notifier/config"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

func init() {

}
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gracefullyShutdown(ctx, cancel)

	var conf config.Notifier
	err := envconfig.Process("notifier", &conf)
	if err != nil {
		log.Panic(err.Error())
	}

	l := parseLogLevel(conf.LogLevel)
	log.SetLevel(l)

	command := cmd.BuildCommand(ctx, &conf)
	err = cmd.ConfigFlags(command, &conf)
	if err != nil {
		log.Error(err.Error())
		return
	}

	if err := command.Execute(); err != nil {
		log.Error("error executing the command")
		return
	}
}

func parseLogLevel(level string) log.Level {
	l, err := log.ParseLevel(level)
	if err != nil {
		l = log.InfoLevel
		log.Errorf("setting log info %v", err)
	}
	return l
}

func gracefullyShutdown(ctx context.Context, cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)

	go func() {
		select {
		case <-c:
			log.Debugf("Shutting down gracefully.")
			cancel()
		case <-ctx.Done():
		}
	}()
}
