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