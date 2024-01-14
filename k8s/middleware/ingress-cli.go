package middleware

import (
	"context"
	"errors"
	"github.com/skhatri/k8s-read/k8s/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type IngressItem struct {
	Namespace    string            `json:"namespace"`
	Kind         string            `json:"kind"`
	Name         string            `json:"name"`
	IngressClass *string           `json:"ingressClass"`
	Hosts        []HostType        `json:"hosts"`
	IP           []string          `json:"ip"`
	Annotations  map[string]string `json:"annotations,omitempty"`
	Labels       map[string]string `json:"labels,omitempty"`
}

type HostType struct {
	Name  string     `json:"name"`
	Tls   bool       `json:"tls"`
	Paths []HostPath `json:"paths"`
}

type HostPath struct {
	Path     string    `json:"path"`
	PathType *string   `json:"pathType"`
	Resource string    `json:"resource"`
	Port     *PortType `json:"port,omitempty"`
	Kind     string    `json:"kind"`
}
type PortType struct {
	Name   string `json:"name"`
	Number int32  `json:"number"`
}

func GetIngress(namespace string, options DisplayOptions) ([]interface{}, error) {

	k8s := client.GetClient()
	if namespace == "" {
		return nil, errors.New("namespace is required")
	}
	if namespace == "any" {
		namespace = ""
	}
	return getIngresses(k8s, namespace, options)
}

func getIngresses(k8s *kubernetes.Clientset, namespace string, options DisplayOptions) ([]interface{}, error) {
	depList, depErr := k8s.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
	if depErr != nil {
		return nil, depErr
	}
	workload := make([]interface{}, 0)
	for _, ingress := range depList.Items {
		if len(options.Names) > 0 {
			found := false
			for _, name := range options.Names {
				if name == ingress.Name {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		spec := ingress.Spec
		tlsMap := make(map[string]string)

		for _, tlsRule := range spec.TLS {
			for _, host := range tlsRule.Hosts {
				tlsMap[host] = tlsRule.SecretName
			}
		}
		hosts := make([]HostType, 0)
		for _, rule := range spec.Rules {
			host := rule.Host
			hostPaths := make([]HostPath, 0)
			for _, httpPath := range rule.HTTP.Paths {
				hostPath := HostPath{
					Path:     httpPath.Path,
					PathType: (*string)(httpPath.PathType),
				}
				if httpPath.Backend.Service != nil {
					serviceName := httpPath.Backend.Service.Name
					port := httpPath.Backend.Service.Port
					hostPath.Resource = serviceName
					hostPath.Kind = "service"
					hostPath.Port = &PortType{
						port.Name,
						port.Number,
					}
				} else if httpPath.Backend.Resource != nil {
					hostPath.Resource = httpPath.Backend.Resource.Name
					hostPath.Kind = httpPath.Backend.Resource.Kind
				}
				hostPaths = append(hostPaths, hostPath)
			}
			tlsEnabled := false
			if _, ok := tlsMap[host]; ok {
				tlsEnabled = true
			}
			hostType := HostType{
				Name:  host,
				Tls:   tlsEnabled,
				Paths: hostPaths,
			}
			hosts = append(hosts, hostType)
		}
		addresses := make([]string, 0)
		for _, ing := range ingress.Status.LoadBalancer.Ingress {
			addresses = append(addresses, ing.IP)
		}
		ingressItem := IngressItem{
			Namespace:    ingress.Namespace,
			Kind:         "Ingress",
			Name:         ingress.Name,
			IngressClass: ingress.Spec.IngressClassName,
			Hosts:        hosts,
			IP:           addresses,
		}
		if options.Labels {
			ingressItem.Labels = ingress.Labels
		}
		if options.Annotations {
			ingressItem.Annotations = ingress.Annotations
		}
		workload = append(workload, ingressItem)
	}
	return workload, nil
}
