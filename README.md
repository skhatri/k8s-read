[![Build](https://travis-ci.com/skhatri/k8s-read.svg?branch=master)](https://travis-ci.com/github/skhatri/k8s-read)
[![Code Coverage](https://img.shields.io/codecov/c/github/skhatri/k8s-read/master.svg)](https://codecov.io/github/skhatri/k8s-read?branch=master)
[![Maintainability](https://api.codeclimate.com/v1/badges/8efb53366c803ff32aff/maintainability)](https://codeclimate.com/github/skhatri/k8s-read/maintainability)

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

