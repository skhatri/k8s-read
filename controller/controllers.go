package controller

import "github.com/skhatri/api-router-go/router"

func Configure(configurer router.ApiConfigurer) {
	configurer.Get("/api/namespaces", namespaceApiHandler)
	configurer.Get("/api/deployments", fetchDeployments)
	configurer.Get("/api/statefulsets", fetchStatefulsets)
	configurer.Get("/api/jobs", fetchJobs)
	configurer.Get("/api/daemonsets", fetchDaemonsets)
}
