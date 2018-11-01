# Beku

[![GoDoc](https://godoc.org/github.com/imroc/beku?status.svg)](https://godoc.org/github.com/yulibaozi/beku)
[![Go Report Card](https://goreportcard.com/badge/github.com/yulibaozi/beku)](https://goreportcard.com/badge/github.com/yulibaozi/beku)

Beku is an extremely user-friendly Kubernetes API resources building library, extremely easy without any extra intelligence. 

### Installation

```
go get -u github.com/yulibaozi/beku
```

### Features

- Extremely simple JSON & YAML input / output
- Required Kubernetes API resources fields automatically confirming
- Interrelated Kubernetes API resources announcement which is so user-friendly 
- Rigorous QOS setup
- Precise fileds auto-fillment
- Graceful chain methods and invocation

### Introduction

Due to the complexities of Kubernetes API resources configuration, miscellaneous fields, diverse hierarchies, rehandling over and over again, Beku was inspirationally born. 

The scenario of Beku is to matching Kubernetes Client-go, and providing json / yaml file for CLI creation. It's very appreciative and helpful that Beku has use Kubernetes codes for reference. 

### Beku-style Usage

1. Chain methods starts with `NewXXX()` and end up with `Finish()`, then we could get whole Kubernetes API resource configuration.
2. All setup methods starts with `SetXXX()` and all retrieves starts with `GetXXX()`.
3. Don't use type cast to satisfying the type needed by some functions as far as possible, it may leads to uncertain errors.
4. There are comments of the usage of some function parameters if you don't know how to handle it.
5. There is a PRESUPPOSE that the first container in Pod has higher status, which will have setup priority. The latter in the sequence of containers, the lower status it has. E.g: Beku will only set the first container's environments when we first invoke the setup function, the next time we invoke it will set the next container.
6. If there is **union** in some struct definition, it means two Kubernetes API resource will be created simultaneously. E.g: Deployment, Service union, PersistentVolume, PersistentVolumeClain union.




How to use beku?
---

### How to quickly Create Service(svc)?

If you want to quicky create service, only input necessary Fields.

As shown in the code:
```
func howToNewSvc() {
	svc, err := beku.NewSvc().SetNamespaceAndName("roc", "mysql-svc").
		SetSelector(map[string]string{"app": "mysql"}).SetServiceType(beku.ServiceTypeNodePort).
		SetPort(beku.ServicePort{Port: 3306, TargetPort: 3306}).Finish()
	if err != nil {
		panic(err)
	}
	yamlbyts, err := beku.ToYAML(svc)
	jsonbyts, err := beku.ToJSON(svc)
	if err != nil {
		panic(err)
	}
}
```

ToYAML
```
apiVersion: v1
kind: Service
metadata:
  creationTimestamp: null
  name: mysql-svc
  namespace: roc
spec:
  ports:
  - port: 3306
    protocol: TCP
    targetPort: 3306
  selector:
    app: mysql
  type: NodePort
status:
  loadBalancer: {}
```
ToJSON
```
{
    "kind":"Service",
    "apiVersion":"v1",
    "metadata":
    {
        "name":"mysql-svc",
        "namespace":"roc",
        "creationTimestamp":null
    },
    "spec":
    {
        "ports":
        [
            {
                "protocol":"TCP",
                "port":3306,
                "targetPort":3306
            }
        ],
        "selector":
        {
            "app":"mysql"
        },
        "type":"NodePort"
    },
    "status":
    {
        "loadBalancer":{}
    }
}
```

### How to quickly Create Deployment?

If you want to quicky create Deployment, only input necessary Fields.

As shown in the code:
```
func howToNewDeployment() {
	dep, err := beku.NewDeployment().
		SetName("yulibaozi").SetLabels(map[string]string{"name": "yulibaozi"}).
		SetContainer("first", "mysql", 3306).Finish()
	if err != nil {
		panic(err)
	}
	yamlByts, err := beku.ToYAML(dep)
	if err != nil {
		panic(err)
	}
	fmt.Println("\n" + string(yamlByts))

```

ToYAML:
```
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    name: yulibaozi
  name: yulibaozi
spec:
  selector:
    matchLabels:
      name: yulibaozi
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        name: yulibaozi
    spec:
      containers:
      - image: mysql
        imagePullPolicy: IfNotPresent
        name: first
        ports:
        - containerPort: 3306
        resources: {}
status: {}
```

