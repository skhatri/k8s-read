package controller

import (
	"github.com/skhatri/api-router-go/router"
	"github.com/skhatri/api-router-go/router/model"
	"github.com/skhatri/k8s-read/k8s/middleware"
)

func nodeHandler(_ *router.WebRequest) *model.Container {
	nodeList, err := middleware.GetNodes()
	if err != nil {
		return model.ErrorResponse(model.MessageItem{
			Code:    "node-error",
			Message: "Could not get nodes list",
		}, 500)
	}
	return model.Response(nodeList)
}
