package controller

import (
	"github.com/skhatri/api-router-go/router"
	"github.com/skhatri/api-router-go/router/model"
	"github.com/skhatri/k8s-read/k8s/middleware"
	"strings"
)

func fetchDeployments(web *router.WebRequest) *model.Container {
	return getWorkload(web, "deployment")
}

func fetchJobs(web *router.WebRequest) *model.Container {
	return getWorkload(web, "job")
}

func fetchDaemonsets(web *router.WebRequest) *model.Container {
	return getWorkload(web, "daemonset")
}

func fetchStatefulsets(web *router.WebRequest) *model.Container {
	return getWorkload(web, "statefulset")
}

func getWorkload(web *router.WebRequest, kind string) *model.Container {
	var namespace = web.GetQueryParam("namespace")
	filterNames := getfilterList(web.GetQueryParam("names"))
	annotations := web.GetQueryParam("annotations")
	labels := web.GetQueryParam("labels")
	displayOptions := middleware.DisplayOptions{
		Names:       filterNames,
		Annotations: annotations == "true",
		Labels:      labels == "true",
	}
	workload, err := middleware.GetWorkload(namespace, kind, displayOptions)
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

func getfilterList(names string) []string {
	filterNames := make([]string, 0)
	if names != "" {
		for _, n := range strings.Split(names, ",") {
			trimName := strings.TrimSpace(n)
			if trimName != "" {
				filterNames = append(filterNames, trimName)
			}
		}
	}
	return filterNames
}
