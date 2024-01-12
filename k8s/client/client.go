package client

import (
	"flag"
	v1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func loadFromKubeConfig() (*rest.Config, error) {
	log.Println("Attempt to load from config")
	var kubeConfigPath *string
	kubeConfigPath = flag.String("kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "kube config file")
	flag.Parse()
	restCfg, restErr := clientcmd.BuildConfigFromFlags("", *kubeConfigPath)
	return restCfg, restErr
}

func assumeServiceAccountAccess() (*rest.Config, error) {
	log.Println("attempt to load from serviceaccount")
	return rest.InClusterConfig()
}

func insideKube() bool {
	apiServerHost := os.Getenv("KUBERNETES_SERVICE_HOST")
	apiServerPort := os.Getenv("KUBERNETES_SERVICE_PORT")
	return len(apiServerHost) > 0 && len(apiServerPort) > 0
}

var clientSet *kubernetes.Clientset
var extClient *v1client.ApiextensionsV1Client
var dynClient dynamic.Interface
var mut = sync.Mutex{}

func GetClient() *kubernetes.Clientset {
	return clientSet
}

func GetExtensionsClient() *v1client.ApiextensionsV1Client {
	return extClient
}

func GetDynamicClient() dynamic.Interface {
	return dynClient
}

func init() {
	mut.Lock()
	defer mut.Unlock()

	var cfg *rest.Config
	var err error
	if insideKube() {
		cfg, err = assumeServiceAccountAccess()
	} else {
		cfg, err = loadFromKubeConfig()
	}

	if err != nil {
		panic(err.Error())
	}

	cset, cerr := kubernetes.NewForConfig(cfg)
	if cerr != nil {
		panic(cerr.Error())
	}
	v1beta1Client, err := v1client.NewForConfig(cfg)
	if err != nil {
		panic(err.Error())
	}
	dClient := dynamic.NewForConfigOrDie(cfg)
	if clientSet == nil {
		clientSet = cset
	}
	if extClient == nil {
		extClient = v1beta1Client
	}
	if dynClient == nil {
		dynClient = dClient
	}
}
