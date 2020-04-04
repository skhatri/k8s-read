### k8s-read
A list of HTTP API to enquire Kubernetes Cluster about the active workload.

It's goal is to be the kubectl for read purpose. 

#### Running App

```
go mod vendor
go build
./k8s-read --port=6100
```

#### List Namespaces

GET /api/namespaces

Output:
```
{
  "data": {
    "namespaces": [
      "default", "kube-system", "kube-public"
    ]
  }
}
```

#### List Workloads in Namespace

GET /api/deployments?namespace=default

Output:
```
{
  data: [{
    "namespace": "default",
    "kind": "deployment",
    "name": "nginx-app",
    "image": "nginx:latest",
    "replicas": 2 
  }]
}
```

