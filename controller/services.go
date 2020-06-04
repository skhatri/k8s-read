package controller

import (
	"github.com/skhatri/api-router-go/router"
	"github.com/skhatri/api-router-go/router/model"
	"github.com/skhatri/k8s-read/k8s/middleware"
)

func getEndpoints(web *router.WebRequest) *model.Container {
	var namespace = web.GetQueryParam("namespace")
	return getServiceItems(namespace, "endpoint")
}
func getServices(web *router.WebRequest) *model.Container {
	var namespace = web.GetQueryParam("namespace")
	return getServiceItems(namespace, "service")
}

func getServiceItems(namespace string, kind string) *model.Container {
	workload, err := middleware.GetServiceItems(namespace, kind)
	if err != nil {
		return model.ErrorResponse(model.MessageItem{
			Code:    "list-error",
			Message: err.Error(),
		}, 500)
	}
	var items = make([]interface{}, 0, len(workload))
	for _, w := range workload {
		items = append(items, w)
	}
	return model.ListResponse(items)
}
