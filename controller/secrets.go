package controller

import (
	"github.com/skhatri/api-router-go/router"
	"github.com/skhatri/api-router-go/router/model"
	"github.com/skhatri/k8s-read/k8s/middleware"
	"strings"
)

func getSecretsFunc(keys *string) ApiHandler {
	publicKeys := strings.Split(*keys, ",")
	return func(web *router.WebRequest) *model.Container {
		var namespace = web.GetQueryParam("namespace")
		filterNames := getfilterList(web.GetQueryParam("names"))
		annotations := web.GetQueryParam("annotations")
		labels := web.GetQueryParam("labels")
		var publicKey = web.GetHeader("X-Request-Public-Key")
		var algoKey = web.GetHeader("X-Request-Encrypt-Algorithm")
		var secretType = web.GetQueryParam("type")
		if secretType == "" {
			secretType = "tls"
		}
		displayOptions := middleware.DisplayOptions{
			Names:       filterNames,
			Annotations: annotations == "true",
			Labels:      labels == "true",
		}
		return getSecretItems(publicKeys, namespace, algoKey, publicKey, &secretType, displayOptions)
	}
}
func contains(items []string, item string) bool {
	for _, value := range items {
		if value == item {
			return true
		}
	}
	return false
}

func getSecretItems(publicKeys []string, namespace string, algo string, publicKey string, secretType *string, options middleware.DisplayOptions) *model.Container {
	if algo != "age" {
		return model.ErrorResponse(model.MessageItem{
			Code:    "request-error",
			Message: "Algorithm not supported",
		}, 400)
	}
	if publicKey == "" || !contains(publicKeys, publicKey) {
		return model.ErrorResponse(model.MessageItem{
			Code:    "request-error",
			Message: "Client is not authorized",
		}, 401)
	}
	workload, err := middleware.GetSecretItems(namespace, middleware.SecretRequest{
		Algorithm:  algo,
		PublicKey:  publicKey,
		SecretType: secretType,
	}, options)
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
