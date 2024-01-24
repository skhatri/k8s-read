package controller

import (
	"github.com/skhatri/api-router-go/router"
	"github.com/skhatri/api-router-go/router/model"
)

type ApiHandler func(web *router.WebRequest) *model.Container

func disabled(web *router.WebRequest) *model.Container {
	return model.ErrorResponse(model.MessageItem{
		Code:    "NotEnabled",
		Message: "Forbidden Path",
	}, 403)
}
