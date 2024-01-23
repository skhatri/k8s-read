package middleware

import (
	"errors"
	"fmt"
	"github.com/skhatri/go-crypt/asymmetric"
	"github.com/skhatri/k8s-read/k8s/client"
	"golang.org/x/net/context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Secret struct {
	Namespace string            `json:"namespace"`
	Name      string            `json:"name"`
	Data      map[string][]byte `json:"data"`
	Type      v1.SecretType     `json:"type"`
}

type SecretRequest struct {
	Algorithm  string
	PublicKey  string
	SecretType *string
}

func GetSecretItems(namespace string, secRequest SecretRequest) ([]interface{}, error) {

	k8s := client.GetClient()
	if namespace == "" {
		return nil, errors.New("namespace is required")
	}
	if namespace == "any" {
		namespace = ""
	}
	return getSecrets(k8s, namespace, secRequest)
}

func getSecrets(k8s *kubernetes.Clientset, namespace string, request SecretRequest) ([]interface{}, error) {
	depList, depErr := k8s.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{})
	if depErr != nil {
		return nil, depErr
	}
	workload := make([]interface{}, 0)
	secType := "kubernetes.io/tls"
	if request.SecretType != nil && *request.SecretType == "Opaque" {
		secType = "Opaque"
	}
	for _, secret := range depList.Items {
		if fmt.Sprintf("%v", secret.Type) == secType {
			encryptData := make(map[string][]byte)
			for k, v := range secret.Data {
				encrypted, encErr := asymmetric.AgeEncrypt(request.PublicKey, string(v))
				if encErr == nil {
					encryptData[k] = []byte(encrypted)
				}
			}

			workload = append(workload, Secret{
				Namespace: secret.Namespace,
				Name:      secret.Name,
				Data:      encryptData,
				Type:      secret.Type,
			})
		}
	}
	return workload, nil
}
