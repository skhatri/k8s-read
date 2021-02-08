package controller

import (
	"github.com/skhatri/api-router-go/router"
	"github.com/skhatri/api-router-go/router/model"
	"github.com/skhatri/k8s-read/k8s/middleware"
)

func fetchPods(web *router.WebRequest) *model.Container {
	var namespace = web.GetQueryParam("namespace")
	var nodeName = web.GetQueryParam("node")
	return getPodWorkload(namespace, nodeName)
}

func getPodWorkload(namespace string, nodeName string) *model.Container {
	workload, err := middleware.GetPods(namespace, nodeName)
	if err != nil {
		return model.ErrorResponse(model.MessageItem{
			Code:    "list-error",
			Message: err.Error(),
		}, 500)
	}
	return model.WithDataOnly(workload)
}
