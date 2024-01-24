[![Build](https://github.com/skhatri/k8s-read/actions/workflows/build.yml/badge.svg)](https://github.com/skhatri/k8s-read/actions/workflows/build.yml)
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

### List Ingresses
GET /api/ingresses

```json
{
  "data": [
    {
      "namespace": "default",
      "kind": "Ingress",
      "name": "httpserver-ingress",
      "ingressClass": "nginx",
      "hosts": [
        {
          "name": "www.example.com",
          "tls": true,
          "paths": [
            {
              "path": "/",
              "pathType": "Prefix",
              "resource": "k8s-read",
              "port": {
                "name": "",
                "number": 6100
              },
              "kind": "service"
            }
          ]
        }
      ],
      "ip": ["192.168.0.3"]
    }
  ]
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

### Gateway API
You can query for objects like TCPRoute, HTTPRoute using the CRD API

#### Get HTTPRoutes
GET /api/crd-instances?resource-group=gateway.networking.k8s.io&resource-type=httproutes&resource-version=v1beta1
```json
{
  "data": [
    {
      "namespace": "default",
      "name": "airflow-http",
      "group": "gateway.networking.k8s.io",
      "version": "v1beta1",
      "resource": "httproutes",
      "link": "/api/crd-instance?resource-group=gateway.networking.k8s.io&resource-type=httproutes&resource-version=v1beta1&namespace=default&resource-name=airflow-http"
    },
    {
      "namespace": "default",
      "name": "k8s-read-http",
      "group": "gateway.networking.k8s.io",
      "version": "v1beta1",
      "resource": "httproutes",
      "link": "/api/crd-instance?resource-group=gateway.networking.k8s.io&resource-type=httproutes&resource-version=v1beta1&namespace=default&resource-name=k8s-read-http"
    }
  ]
}
```

#### Get TCPRoutes
GET /api/crd-instances?resource-group=gateway.networking.k8s.io&resource-type=tcproutes&resource-version=v1alpha2

```json 
{
  "data": [
    {
      "namespace": "default",
      "name": "postgres-endpoint",
      "group": "gateway.networking.k8s.io",
      "version": "v1alpha2",
      "resource": "tcproutes",
      "link": "/api/crd-instance?resource-group=gateway.networking.k8s.io&resource-type=tcproutes&resource-version=v1alpha2&namespace=default&resource-name=postgres-endpoint"
    },
    {
      "namespace": "default",
      "name": "cassandra-endpoint",
      "group": "gateway.networking.k8s.io",
      "version": "v1alpha2",
      "resource": "tcproutes",
      "link": "/api/crd-instance?resource-group=gateway.networking.k8s.io&resource-type=tcproutes&resource-version=v1alpha2&namespace=default&resource-name=cassandra-endpoint"
    }
  ]
}
```


#### Get Secrets
Since secrets are not meant to be intercepted in transit, we would like to encrypt each entry with provided public key. For this we use ```age```

We also have additional settings we need to configure. Secret endpoint is disabled by default and you need to set ```secret_endpoint``` to true.

Relevant snippet from router.json is presented below.
```json
{
  "toggles": {
    "daemonset_endpoint": true,
    "secret_endpoint": true
  }
}
```

Similarly, we need to whitelist public keys of clients who will be calling the secrets endpoint. This is done with the assumpption that the data can only
be decrypted with expected private keys. It is client's responsibility to keep the private key secure.

The list of public keys can be provided in a comma separate list under variable ```public-keys``` in router.json

```
{
    "variables": {
        "public-keys": "age10qq6fyrurpkhg7nnt98ccewnvy6utpaf54rmjesq68c6qp9s99rsgamn0z,age1gn26zalgf5xn5dn04lxemu4x4uapvkgh3jf4ajqwxklxdtdtdd3sy83wcx"
    }
}
```

Once the application is configured with toggle and public key list, you may call the endpoint to retrieve the secrets.

curl -H "x-request-encrypt-algorithm: age" \
    -H"x-request-public-key: age1gn26zalgf5xn5dn04lxemu4x4uapvkgh3jf4ajqwxklxdtdtdd3sy83wcx" \
    "https://localhost:6100/api/secrets?namespace=default&type="

```json
{
  "data": [
    {
      "namespace": "default",
      "name": "k8s-read",
      "data": {
        "tls.crt": "encrypted cert",
        "tls.key": "encrypted key",
        "type": "kubernetes.io/tls"
      }
    }
  ]
}
```
when type parameter is empty, it defaults to tls. To retrieve, Opaque secrets, use type=Opaque.

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

#### TLS
A self-signed key/cert is provided to run k8s-read with TLS enabled.

Create certificate like so

```
openssl genrsa -out private.key 2048
openssl req -new -x509 -sha256 -key private.key -out cert.pem -days 730 -subj "/C=AU/ST=NSW/L=SYD/O=OSS/OU=IT/CN=k8s-read"
```

Update router.json to enable or disable TLS

```
  "transport": {
    "port": 6100,
    "tls": {
      "enabled": true,
      "private-key": "private.key",
      "public-key": "cert.pem"
    }
  }

```

