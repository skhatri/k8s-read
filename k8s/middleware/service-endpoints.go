package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/skhatri/k8s-read/k8s/client"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ServiceItem struct {
	Namespace      string  `json:"namespace"`
	Kind           string  `json:"kind"`
	Name           string  `json:"name"`
	ServiceType    string  `json:"serviceType"`
	ClusterIP      string  `json:"ip"`
	ExternalRef    *string `json:"externalRef,omitempty"`
	LoadBalancerIP *string `json:"loadBalancerIP,omitempty"`
}

type EndpointItem struct {
	Namespace string    `json:"namespace"`
	Kind      string    `json:"kind"`
	Name      string    `json:"name"`
	Addresses []Address `json:"addresses"`
}
type Address struct {
	IP         string  `json:"ip"`
	HostName   string  `json:"hostname"`
	NodeName   *string `json:"nodename,omitempty"`
	TargetKind *string `json:"kind,omitempty"`
	TargetName *string `json:"targetName,omitempty"`
}

func GetServiceItems(namespace string, kind string) ([]interface{}, error) {

	k8s := client.GetClient()
	if namespace == "" {
		return nil, errors.New("namespace is required")
	}
	if namespace == "any" {
		namespace = ""
	}
	switch kind {
	case "service":
		return getServices(k8s, namespace)
	case "endpoint":
		return getEndpoints(k8s, namespace)
	}
	return nil, errors.New(fmt.Sprintf("workload retrieval for kind %s is unsupported", kind))
}

func getEndpoints(k8s *kubernetes.Clientset, namespace string) ([]interface{}, error) {
	depList, depErr := k8s.CoreV1().Endpoints(namespace).List(context.TODO(), metav1.ListOptions{})
	if depErr != nil {
		return nil, depErr
	}
	workload := make([]interface{}, 0)
	for _, endpoint := range depList.Items {
		addresses := make([]Address, 0)
		for _, subset := range endpoint.Subsets {
			for _, address := range subset.Addresses {
				var targetKind *string
				var targetName *string
				if targetRef := address.TargetRef; targetRef != nil {
					targetKind = &targetRef.Kind
					targetName = &targetRef.Name
				}
				addresses = append(addresses, Address{
					IP:         address.IP,
					HostName:   address.Hostname,
					NodeName:   address.NodeName,
					TargetKind: targetKind,
					TargetName: targetName,
				})

			}
		}
		workload = append(workload, EndpointItem{
			Namespace: endpoint.Namespace,
			Kind:      "endpoint",
			Name:      endpoint.Name,
			Addresses: addresses,
		})
	}
	return workload, nil
}

func getServices(k8s *kubernetes.Clientset, namespace string) ([]interface{}, error) {
	depList, depErr := k8s.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if depErr != nil {
		return nil, depErr
	}
	workload := make([]interface{}, 0)
	for _, service := range depList.Items {
		typeValue := string(service.Spec.Type)
		clusterIP := service.Spec.ClusterIP
		var loadBalancerIP *string
		var externalRef *string
		switch service.Spec.Type {
		case v1.ServiceTypeLoadBalancer:
			loadBalancerIP = &service.Spec.LoadBalancerIP
		case v1.ServiceTypeExternalName:
			externalRef = &service.Spec.ExternalName
		}

		workload = append(workload, ServiceItem{
			Namespace:      service.Namespace,
			Kind:           "service",
			Name:           service.Name,
			ServiceType:    typeValue,
			ClusterIP:      clusterIP,
			ExternalRef:    externalRef,
			LoadBalancerIP: loadBalancerIP,
		})
	}
	return workload, nil
}
