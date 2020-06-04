package main

import (
	"github.com/sirupsen/logrus"
	"github.com/skhatri/api-router-go/router"
	"github.com/skhatri/api-router-go/starter"
	"github.com/skhatri/k8s-read/controller"
	"os"
)

func main() {

	starter.StartAppWithOptions(os.Args, 6100, func(configurer router.ApiConfigurer) {
		controller.Configure(configurer)
	}, newLogger())
}

func newLogger() func(...interface{}) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.999-0700",
	})

	var loggingFunc = func(args ...interface{}) {
		logger.Info(args)
	}
	return loggingFunc
}
