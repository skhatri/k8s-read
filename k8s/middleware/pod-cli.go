package middleware

import (
	"context"
	"errors"
	"github.com/skhatri/k8s-read/k8s/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type ResourceRequirement struct {
	Cpu    int64 `json:"cpu"`
	Memory int64 `json:"memory"`
}

type Pod struct {
	Namespace string              `json:"namespace"`
	Kind      string              `json:"kind"`
	Name      string              `json:"name"`
	Image     string              `json:"image"`
	Request   ResourceRequirement `json:"request"`
	Limit     ResourceRequirement `json:"limit"`
	Node      string              `json:"node"`
}

type PodSummary struct {
	NodeWorkloads []NodeWorkload      `json:"nodes"`
	Request       ResourceRequirement `json:"request"`
	Limit         ResourceRequirement `json:"limit"`
}

type NodeWorkload struct {
	Node             string              `json:"node"`
	Pods             []Pod               `json:"pods"`
	Request          ResourceRequirement `json:"request"`
	Limit            ResourceRequirement `json:"limit"`
	Allocatable      ResourceRequirement `json:"allocatable"`
	Capacity         ResourceRequirement `json:"capacity"`
	Labels           map[string]string   `json:"labels"`
	Annotations      map[string]string   `json:"annotations"`
	CreatedTimestamp time.Time           `json:"createdTimestamp"`
}

func GetPods(namespace string, nodeName string) (*PodSummary, error) {
	k8s := client.GetClient()
	if namespace == "" {
		return nil, errors.New("namespace is required")
	}
	if namespace == "any" {
		namespace = ""
	}
	podList, podErr := k8s.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if podErr != nil {
		return nil, podErr
	}
	workload := make(map[string][]Pod, 0)
	for _, pod := range podList.Items {
		if nodeName != "" && nodeName != pod.Spec.NodeName {
			continue
		}
		hostNode := pod.Spec.NodeName
		podData := Pod{
			Namespace: pod.Namespace,
			Kind:      "pod",
			Name:      pod.Name,
			Request: ResourceRequirement{
				Cpu:    0,
				Memory: 0,
			},
			Node: hostNode,
		}
		for _, container := range pod.Spec.Containers {
			podData.Request.Cpu += container.Resources.Requests.Cpu().AsDec().UnscaledBig().Int64()
			podData.Request.Memory += container.Resources.Requests.Memory().AsDec().UnscaledBig().Int64()

			podData.Limit.Cpu += container.Resources.Limits.Cpu().AsDec().UnscaledBig().Int64()
			podData.Limit.Memory += container.Resources.Limits.Memory().AsDec().UnscaledBig().Int64()
		}
		existingPods, ok := workload[hostNode]
		if !ok {
			existingPods = make([]Pod, 0)
		}
		existingPods = append(existingPods, podData)
		workload[hostNode] = existingPods
	}

	nodes, err := GetNodes()
	if err != nil {
		return nil, err
	}
	nodeSummary := make(map[string]Node, 0)
	for _, node := range nodes.Nodes {
		nodeSummary[node.Name] = node
	}

	hostWorkloads := make([]NodeWorkload, 0)
	overallRequest := ResourceRequirement{Cpu: 0, Memory: 0}
	overallLimit := ResourceRequirement{Cpu: 0, Memory: 0}
	for k, v := range workload {
		if k == "" {
			continue
		}
		request := ResourceRequirement{Cpu: 0, Memory: 0}
		limit := ResourceRequirement{Cpu: 0, Memory: 0}
		for _, item := range v {
			if item.Request.Cpu <= item.Limit.Cpu {
				request.Cpu += item.Request.Cpu
			} else {
				request.Cpu += item.Limit.Cpu
			}
			if item.Request.Memory <= item.Limit.Memory {
				request.Memory += item.Request.Memory
			} else {
				request.Memory += item.Limit.Memory
			}
			limit.Cpu += item.Limit.Cpu
			limit.Memory += item.Limit.Memory
		}
		nodeSummaryData := nodeSummary[k]
		hostWorkloads = append(hostWorkloads, NodeWorkload{
			Node:             k,
			Pods:             v,
			Request:          request,
			Limit:            limit,
			Allocatable:      nodeSummaryData.Allocatable,
			Capacity:         nodeSummaryData.Capacity,
			Labels:           nodeSummaryData.Labels,
			Annotations:      nodeSummaryData.Annotations,
			CreatedTimestamp: nodeSummaryData.CreatedTimestamp,
		})
		overallRequest.Cpu += request.Cpu
		overallRequest.Memory += request.Memory

		overallLimit.Cpu += limit.Cpu
		overallLimit.Memory += limit.Memory

	}

	podSummary := PodSummary{
		NodeWorkloads: hostWorkloads,
		Request:       overallRequest,
		Limit:         overallLimit,
	}
	return &podSummary, nil
}
