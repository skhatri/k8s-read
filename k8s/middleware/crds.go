package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/skhatri/k8s-read/k8s/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

//GetCrdByName Kubernetes Workload of given Custom Resource Type in a namespace.
func GetCrdByName(namespace string, gvr schema.GroupVersionResource, resourceName string) (*CustomResourceInstance, error) {
	if namespace == "" {
		return nil, errors.New("namespace is required")
	}
	dynamicClient := *(client.GetDynamicClient())

	namespaceResInt := dynamicClient.Resource(gvr).Namespace(namespace)
	resource, err := namespaceResInt.Get(context.TODO(), resourceName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("CRD Resource List error %s", err.Error()))
	}

	resourceData, err := resource.MarshalJSON()
	var cres = CustomResourceInstance{}
	if err == nil {
		buff := bytes.NewBuffer(resourceData)
		json.NewDecoder(buff).Decode(&cres)
	}

	return &cres, nil
}

//GetCrdInstanceList returns custom resource instances for a group
func GetCrdInstanceList(namespace string, gvr schema.GroupVersionResource) ([]CustomResourceInstanceSummary, error) {
	dynamicClient := *(client.GetDynamicClient())

	var namespaceResInt dynamic.ResourceInterface = dynamicClient.Resource(gvr)
	if namespace != "" {
		namespaceResInt = namespaceResInt.(dynamic.NamespaceableResourceInterface).Namespace(namespace)
	}
	resources, err := namespaceResInt.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("CRD Resource List error %s", err.Error()))
	}

	customResources := make([]CustomResourceInstanceSummary, 0)
	for _, res := range resources.Items {
		groupKind := res.GroupVersionKind()
		customResources = append(customResources, CustomResourceInstanceSummary{
			Namespace: res.GetNamespace(),
			Name:    res.GetName(),
			Version: groupKind.Version,
			Group:   groupKind.Group,
			Resource: gvr.Resource,
			Link: fmt.Sprintf("/api/crd-instance?resource-group=%s&resource-type=%s&resource-version=%s&namespace=%s&resource-name=%s",
				groupKind.Group, gvr.Resource, groupKind.Version, res.GetNamespace(), res.GetName()),
		})
	}
	return customResources, nil
}

//GetCrds returns list of CRDs registered against the Api Server
func GetCrds() ([]CrdSummary, error) {
	crds, err := client.GetExtensionsClient().CustomResourceDefinitions().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	crdList := make([]CrdSummary, 0)
	for _, crd := range crds.Items {
		var effectiveVersion = ""
		for _, v := range crd.Spec.Versions {
			if v.Storage && v.Served {
				effectiveVersion = v.Name
				break
			}
		}
		crdList = append(crdList, CrdSummary{
			Name:         crd.Name,
			Group:        crd.Spec.Group,
			Kind:         crd.Spec.Names.Kind,
			Version:      effectiveVersion,
			ResourceType: crd.Spec.Names.Plural,
			Link: fmt.Sprintf("/api/crd-instances?resource-group=%s&resource-type=%s&resource-version=%s",
				crd.Spec.Group, crd.Spec.Names.Plural, effectiveVersion),
		})
	}
	return crdList, nil
}

type CustomResourceInstance struct {
	Spec     map[string]interface{} `json:"spec"`
	Metadata map[string]interface{} `json:"metadata"`
}
type CustomResourceInstanceSummary struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Group     string `json:"group"`
	Version   string `json:"version"`
	Resource   string `json:"resource"`
	Link string `json:"link"`
}

type CrdSummary struct {
	Name         string `json:"name"`
	Group        string `json:"group"`
	ResourceType string `json:"resource-type"`
	Kind         string `json:"kind"`
	Version      string `json:"version"`
	Link string `json:"link"`
}
