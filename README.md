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
### List Custom Resource Definitions
GET /api/crds
```
{
    "data": [
        {
            "name": "cities.world.io",
            "group": "world.io",
            "resource-type": "cities",
            "kind": "City",
            "version": "v1alpha1",
            "link": "/api/crd-instances?resource-group=world.io&resource-type=cities&resource-version=v1alpha1"
        }
    ]
}
```

### List Custom Resource Instances
GET /api/crd-instances?resource-group=world.io&resource-type=cities&resource-version=v1alpha1

```
{
    "data": [
        {
            "namespace": "default",
            "name": "s-cities",
            "group": "world.io",
            "version": "v1alpha1",
            "kind": "City",
            "link": "/api/crd-instance?resource-group=world.io&resource-type=cities&resource-version=v1alpha1&namespace=default&resource-name=s-cities"
        }
    ]
}   
```

### Get Custom Resource Instance
GET /api/crds?namespace=default&resource-type=cities&resource-group=world.io&resource-version=v1alpha1&resource-name=s-cities

```
{
    "data": {
        "spec": {
            "apps": [
                {
                    "country": "Australia",
                    "name": "sydney"
                },
                {
                    "country": "USA",
                    "name": "san francisco"
                }
            ]
        },
        "metadata": {
            "annotations": {
                "living-expense": "high",
                "weather": "great"
            },
            "creationTimestamp": "2020-04-19T00:22:21Z",
            "generation": 1,
            "labels": {
                "starts-with": "S"
            },
            "name": "s-cities",
            "namespace": "default",
            "resourceVersion": "176136",
            "selfLink": "/apis/world.io/v1alpha1/namespaces/default/cities/s-cities",
            "uid": "d5a9c026-81d3-11ea-b6f1-02430e0005fc"
        }
    }
}
```
#### Filtering
The data can be filtered by additionally providing the following three parameters.

|Parameter|Description|
|---|---|
|annotations|Whether to display annotation. default is false.|
|labels|Whether to display labels. default is false.|
|names|Object names to filter. comma separated names|


#### Docker
```
docker build --no-cache -t k8s-read .
```
#### Deploy
```
kubectl apply -f deploy/
```

#### Additional Roles
```
kubectl apply -f deploy/rbac/
```

