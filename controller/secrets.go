package controller

import (
	"github.com/skhatri/api-router-go/router"
	"github.com/skhatri/api-router-go/router/model"
	"github.com/skhatri/k8s-read/k8s/middleware"
)

func getSecrets(web *router.WebRequest) *model.Container {
	var namespace = web.GetQueryParam("namespace")
	var publicKey = web.GetHeader("X-Request-Public-Key")
	var algoKey = web.GetHeader("X-Request-Encrypt-Algorithm")
	var secretType = web.GetQueryParam("type")
	if secretType == "" {
		secretType = "tls"
	}
	return getSecretItems(namespace, algoKey, publicKey, &secretType)
}

func getSecretItems(namespace string, algo string, publicKey string, secretType *string) *model.Container {
	if algo != "age" {
		return model.ErrorResponse(model.MessageItem{
			Code:    "request-error",
			Message: "Algorithm not supported",
		}, 400)
	}
	workload, err := middleware.GetSecretItems(namespace, middleware.SecretRequest{
		Algorithm:  algo,
		PublicKey:  publicKey,
		SecretType: secretType,
	})
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
