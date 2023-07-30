package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/skhatri/k8s-read/k8s/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Workload struct {
	Namespace   string            `json:"namespace,omitempty"`
	Kind        string            `json:"kind"`
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	Replicas    int32             `json:"replicas"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

type DisplayOptions struct {
	Names       []string
	Annotations bool
	Labels      bool
}

//Get Kubernetes Workload of given kind in a namespace.
func GetWorkload(namespace string, kind string, displayOptions DisplayOptions) ([]Workload, error) {
	k8s := client.GetClient()
	if namespace == "" {
		return nil, errors.New("namespace is required")
	}
	if namespace == "any" {
		namespace = ""
	}
	switch kind {
	case "deployment":
		return getDeployments(k8s, namespace, displayOptions)
	case "statefulset":
		return getStatefulSets(k8s, namespace, displayOptions)
	case "daemonset":
		return getDaemonSets(k8s, namespace, displayOptions)
	case "job":
		return getJobs(k8s, namespace, displayOptions)
	}
	return nil, errors.New(fmt.Sprintf("workload retrieval for kind %s is unsupported", kind))
}

func getDeployments(k8s *kubernetes.Clientset, namespace string, displayOptions DisplayOptions) ([]Workload, error) {
	depList, depErr := k8s.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if depErr != nil {
		return nil, depErr
	}
	workloads := make([]Workload, 0)
	for _, deployment := range depList.Items {
		if !shouldDisplay(displayOptions, deployment.Name) {
			continue
		}
		for _, container := range deployment.Spec.Template.Spec.Containers {
			item := Workload{
				Namespace: deployment.Namespace,
				Kind:      "deployment",
				Name:      deployment.Name,
				Image:     container.Image,
				Replicas:  *deployment.Spec.Replicas,
			}
			if displayOptions.Annotations {
				item.Annotations = deployment.Annotations
			}
			if displayOptions.Labels {
				item.Labels = deployment.Labels
			}
			workloads = append(workloads, item)
		}
	}
	return workloads, nil
}

func shouldDisplay(displayOptions DisplayOptions, name string) bool {
	display := true
	if len(displayOptions.Names) > 0 {
		display = false
		for _, d := range displayOptions.Names {
			if d == name {
				display = true
			}
		}
	}
	return display
}

func getStatefulSets(k8s *kubernetes.Clientset, namespace string, displayOptions DisplayOptions) ([]Workload, error) {
	stList, stErr := k8s.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if stErr != nil {
		return nil, stErr
	}
	workloads := make([]Workload, 0)
	for _, statefulset := range stList.Items {
		if !shouldDisplay(displayOptions, statefulset.Name) {
			continue
		}
		for _, container := range statefulset.Spec.Template.Spec.Containers {
			item := Workload{
				Namespace: statefulset.Namespace,
				Kind:      "statefulset",
				Name:      statefulset.Name,
				Image:     container.Image,
				Replicas:  *statefulset.Spec.Replicas,
			}
			if displayOptions.Annotations {
				item.Annotations = statefulset.Annotations
			}
			if displayOptions.Labels {
				item.Labels = statefulset.Labels
			}
			workloads = append(workloads, item)
		}
	}
	return workloads, nil
}

func getDaemonSets(k8s *kubernetes.Clientset, namespace string, displayOptions DisplayOptions) ([]Workload, error) {
	stList, stErr := k8s.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if stErr != nil {
		return nil, stErr
	}
	workloads := make([]Workload, 0)
	for _, daemonset := range stList.Items {
		if !shouldDisplay(displayOptions, daemonset.Name) {
			continue
		}
		for _, container := range daemonset.Spec.Template.Spec.Containers {
			item := Workload{
				Namespace: daemonset.Namespace,
				Kind:      "daemonset",
				Name:      daemonset.Name,
				Image:     container.Image,
				Replicas:  1,
			}
			if displayOptions.Annotations {
				item.Annotations = daemonset.Annotations
			}
			if displayOptions.Labels {
				item.Labels = daemonset.Labels
			}
			workloads = append(workloads, item)
		}
	}
	return workloads, nil
}

func getJobs(k8s *kubernetes.Clientset, namespace string, displayOptions DisplayOptions) ([]Workload, error) {
	jobList, jobErr := k8s.BatchV1().Jobs(namespace).List(context.TODO(), metav1.ListOptions{})
	if jobErr != nil {
		return nil, jobErr
	}
	workloads := make([]Workload, 0)
	for _, job := range jobList.Items {
		if !shouldDisplay(displayOptions, job.Name) {
			continue
		}
		for _, container := range job.Spec.Template.Spec.Containers {
			item := Workload{
				Namespace: job.Namespace,
				Kind:      "job",
				Name:      job.Name,
				Image:     container.Image,
				Replicas:  1,
			}
			if displayOptions.Annotations {
				item.Annotations = job.Annotations
			}
			if displayOptions.Labels {
				item.Labels = job.Labels
			}
			workloads = append(workloads, item)
		}
	}
	return workloads, nil
}
