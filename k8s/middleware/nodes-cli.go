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
	Annotations map[string]string `json:"annotations"`
	Labels      map[string]string `json:"labels"`
	Taints    []v1.Taint `json:"taints"`
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
			Annotations: item.Annotations,
			Labels: item.Labels,
			Taints: item.Spec.Taints,
		})
	}
	return &NodeList{Nodes: nodeItems}, nil
}
