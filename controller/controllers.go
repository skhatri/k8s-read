package controller

import (
	"github.com/skhatri/api-router-go/router"
	"github.com/skhatri/api-router-go/router/functions"
	"github.com/skhatri/api-router-go/router/model"
	"github.com/skhatri/api-router-go/router/settings"
)

func Configure(configurer router.ApiConfigurer) {
	var _settings = settings.GetSettings()
	configurer.Get("/api/namespaces", namespaceApiHandler).
		Get("/status", functions.StatusFunc).
		Get("/api/deployments", fetchDeployments).
		Get("/api/statefulsets", fetchStatefulsets).
		Get("/api/pods", fetchPods).
		Get("/api/jobs", fetchJobs).
		Get("/api/nodes", nodeHandler).
		Get("/api/crd-instances", getCrdInstanceList).
		Get("/api/crd-instance", getCrdInstance).
		Get("/api/crds", getCrds).
		Get("/api/endpoints", getEndpoints).
		Get("/api/services", getServices).
		Get("/api/ingresses", getIngresses).
		GetIf(_settings.IsToggleOn("secret_endpoint")).
		Register("/api/secrets",
			func() func(request *router.WebRequest) *model.Container {
				keys := _settings.Variable("public-keys")
				if keys == nil || len(*keys) == 0 {
					return disabled
				}
				return getSecretsFunc(keys)
			}()).
		GetIf(_settings.IsToggleOn("daemonset_endpoint")).Register("/api/daemonsets", fetchDaemonsets)
}
