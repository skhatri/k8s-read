package controller

import (
	"github.com/skhatri/api-router-go/router"
	"github.com/skhatri/api-router-go/router/model"
	"github.com/skhatri/k8s-read/k8s/middleware"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func getCrdInstance(web *router.WebRequest) *model.Container {
	gvr := schema.GroupVersionResource{
		Resource: web.GetQueryParam("resource-type"),
		Group:    web.GetQueryParam("resource-group"),
		Version:  web.GetQueryParam("resource-version"),
	}
	cres, err := middleware.GetCrdByName(
		web.GetQueryParam("namespace"), gvr, web.GetQueryParam("resource-name"))
	if err != nil {
		return model.ErrorResponse(model.MessageItem{
			Code:    "crd get error",
			Message: err.Error(),
		}, 500)
	}
	return model.Response(cres)
}

func getCrdInstanceList(web *router.WebRequest) *model.Container {
	gvr := schema.GroupVersionResource{
		Resource: web.GetQueryParam("resource-type"),
		Group:    web.GetQueryParam("resource-group"),
		Version:  web.GetQueryParam("resource-version"),
	}
	cresList, err := middleware.GetCrdInstanceList(
		web.GetQueryParam("namespace"), gvr)
	if err != nil {
		return model.ErrorResponse(model.MessageItem{
			Code:    "crd instance list error",
			Message: err.Error(),
		}, 500)
	}
	return model.Response(cresList)
}

func getCrds(web *router.WebRequest) *model.Container {
	crdList, err := middleware.GetCrds()
	if err != nil {
		return model.ErrorResponse(model.MessageItem{
			Code:    "crd list error",
			Message: err.Error(),
		}, 500)
	}
	return model.Response(crdList)
}
