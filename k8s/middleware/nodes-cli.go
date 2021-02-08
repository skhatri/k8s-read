package middleware

import (
	"context"
	"github.com/skhatri/k8s-read/k8s/client"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NodeList struct {
	Nodes []Node `json:"nodes"`
}

type Node struct {
	Name        string              `json:"name"`
	Annotations map[string]string   `json:"annotations"`
	Labels      map[string]string   `json:"labels"`
	Taints      []v1.Taint          `json:"taints"`
	Capacity    ResourceRequirement `json:"capacity"`
	Allocatable ResourceRequirement `json:"allocatable"`
}

//GetNodes returns list of namespaces in the current cluster
func GetNodes() (*NodeList, error) {
	client := client.GetClient()
	nodeList, nameErr := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{
	})
	if nameErr != nil {
		return nil, nameErr
	}
	var nodeItems = make([]Node, 0, len(nodeList.Items))
	for _, item := range nodeList.Items {
		nodeItems = append(nodeItems, Node{
			Name:        item.Name,
			Annotations: item.Annotations,
			Labels:      item.Labels,
			Taints:      item.Spec.Taints,
			Capacity: ResourceRequirement{
				Cpu:    1000 * item.Status.Capacity.Cpu().AsDec().UnscaledBig().Int64(),
				Memory: item.Status.Capacity.Memory().AsDec().UnscaledBig().Int64(),
			},
			Allocatable: ResourceRequirement{
				Cpu:    1000 * item.Status.Allocatable.Cpu().AsDec().UnscaledBig().Int64(),
				Memory: item.Status.Allocatable.Memory().AsDec().UnscaledBig().Int64(),
			},
		})
	}
	return &NodeList{Nodes: nodeItems}, nil
}
