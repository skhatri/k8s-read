package main

import (
	"github.com/sirupsen/logrus"
	"github.com/skhatri/api-router-go/router"
	"github.com/skhatri/api-router-go/starter"
	"github.com/skhatri/k8s-read/controller"
)

func main() {
  starter.RunAppWithOptions(func(configurer router.ApiConfigurer){
    controller.Configure(configurer)
  }, newLogger())
}

func newLogger() func(router.RequestSummary) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.999-0700",
	})

	var loggingFunc = func(requestSummary router.RequestSummary) {
                if requestSummary.Status >= 400 && requestSummary.Status != 404 {
		  logger.WithField("uri", requestSummary.Uri).
                    WithField("status_code", requestSummary.Status).
                    WithField("time_taken", requestSummary.TimeTaken).WithField("unit", requestSummary.Unit).Error()
                }
	}
	return loggingFunc
}

