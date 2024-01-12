package controller

import (
	"github.com/skhatri/api-router-go/router"
	"github.com/skhatri/api-router-go/router/model"
	"github.com/skhatri/k8s-read/k8s/middleware"
)

func getIngresses(web *router.WebRequest) *model.Container {
	var namespace = web.GetQueryParam("namespace")
	return getIngressItem(namespace)
}

func getIngressItem(namespace string) *model.Container {
	workload, err := middleware.GetIngress(namespace)
	if err != nil {
		return model.ErrorResponse(model.MessageItem{
			Code:    "ingress-list-error",
			Message: err.Error(),
		}, 500)
	}
	var items = make([]interface{}, 0, len(workload))
	for _, w := range workload {
		items = append(items, w)
	}
	return model.ListResponse(items)
}
