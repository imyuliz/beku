# Beku

[![GoDoc](https://godoc.org/github.com/imroc/beku?status.svg)](https://godoc.org/github.com/yulibaozi/beku)
[![Go Report Card](https://goreportcard.com/badge/github.com/yulibaozi/beku)](https://goreportcard.com/badge/github.com/yulibaozi/beku)

Beku is an extremely user-friendly Kubernetes API resources building library, extremely easy without any extra intelligence. 

### Installation

```
go get -u github.com/yulibaozi/beku
```

### Features

- Auto release resource object on Kubernetes
- Flexible custom development
- Extremely simple JSON & YAML input / output
- Required Kubernetes API resources fields automatically confirming
- Interrelated Kubernetes API resources announcement which is so user-friendly 
- Rigorous QOS setup
- Precise fileds auto-fillment
- Graceful chain methods and invocation


### Document

- [中文](https://github.com/yulibaozi/beku/blob/master/doc/README-cn.md)
- [More examples](https://github.com/yulibaozi/beku/blob/master/test/example_test.go)
- [Youtube:Deploy your application on Kubernetes with 3 LoC using Beku](https://youtu.be/4CaARsch9ms)
- [Tencent Video:Beku--3行代码发布你的应用到Kubernetes](http://v.qq.com/x/page/d0783vtazs9.html)

### Introduction

Due to the complexities of Kubernetes API resources configuration, miscellaneous fields, diverse hierarchies, rehandling over and over again, Beku was inspirationally born. 

The scenario of Beku is to matching Kubernetes Client-go, and providing json / yaml file for CLI creation. It's very appreciative and helpful that Beku has use Kubernetes codes for reference. 
###  caas-one
    ![caas-one](https://github.com/yulibaozi/beku/blob/master/doc/caas-one.jpeg)
### Beku-style Usage

1. Chain methods starts with `NewXXX()` and end up with `Finish()`, then we could get whole Kubernetes API resource configuration.
2. All setup methods starts with `SetXXX()` and all retrieves starts with `GetXXX()`.
3. Don't use type cast to satisfying the type needed by some functions as far as possible, it may leads to uncertain errors.
4. There are comments of the usage of some function parameters if you don't know how to handle it.
5. There is a PRESUPPOSE that the first container in Pod has higher status, which will have setup priority. The latter in the sequence of containers, the lower status it has. E.g: Beku will only set the first container's environments when we first invoke the setup function, the next time we invoke it will set the next container.
6. If there is **union** in some struct definition, it means two Kubernetes API resource will be created simultaneously. E.g: Deployment, Service union, PersistentVolume, PersistentVolumeClain union.

### Examples

How to create a Service(svc) in few seconds?

As you will see:

```go
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

```yaml
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

```json
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
More examples: [Example.md](https://github.com/yulibaozi/beku/blob/master/doc/example.md)

### Currently supported Kubernetes API resources in Beku

Kubernetes API resources | Abbreviation | Version 
---|---|---|
namespace   | ns| core/v1
service   | svc| core/v1
deployment | - | apps/v1
statefulset | sts | apps/v1
secret | - | core/v1
persistentVolumeClaim | pvc | core/v1
persistentVolume | pv | core/v1
daemonSet | ds | apps/v1
configMap | cm | core/v1

### Beku Implementation Strategy

1. Only one API resource version will be implemented **even it has multi versions**. Since stability instead of diversity is the first place concern, which may lead some lags to the latest version. But it won't be a problem. On the other hand, below are priorities when choosing API resource version:
	* core/v1 
	* apps/v1
	...

2. When implementing Kubernetes API resource objects which lack of stable version, Alpha, Beta versions will not be implemented. Because less stable versions have more probablities to be changed. 
3. Kubernetes API resource versions references:
[Kubernetes API Introduction](http://kubernetes.kansea.com/docs/api/)

### Beku Conception

**In the past**, I found that it's pretty tedious to write Kubernetes API reosurces configuration, besides the complexity of configurable fields, there still need extra intelligence like:
 * It's a problem that configurable fields which is required or optional.
 * Kubernetes API resources have diverse hierarchies, it's a problem the localtions of configurable fields in json/yaml file.
 * Indentation need to taken in to consideration when writing yaml file.
 * It's a problem that those issues above happen over and over again. 

There are drawbacks in **current** implementations, e.g. we implemented the general fill of some fields in one Kubernetes API resource instead of the complex way, this may lead multi advanced fill strategies not available. You could propose a PR or issue for disscution. On the other hand, there are still some capabilities need to be completed, we could work together.

**In the future**, expected to eliminate the 3 drawbacks above and make a progress to the targets below:
- [ ] Invoker provide required fields, non-need fields don't fill.
- [x] Invoker provide non-hierachy fields, extra intelliengcy burden exliminated.
- [x] Invoker provide requreid fields to make complete json/yaml configuration.
- [x] Invoker provide single field, other related fields will be filled automatically.
- [x] Invoker provide incomplete fields, Beku will return lacked fields helping invokers to accomplish fill. 
