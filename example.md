How to quickly use beku?
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
	fmt.Println("\n" + string(yamlbyts))
	fmt.Println("\n" + string(jsonbyts))
}
```

### How to quickly Create StatefulSet(sts)?

If you want to quicky create sts, only input necessary Fields.

As shown in the code:

```
func howToNewSts() {
	sts, err := beku.NewSts().SetNamespaceAndName("yulibaozi", "mysql").SetSelector(map[string]string{"name": "mysql"}).SetContainer("first", "mysql", 3306).Finish()
	if err != nil {
		panic(err)
	}
	yamlByts, err := beku.ToYAML(sts)
	if err != nil {
		panic(err)
	}
	fmt.Println("\n" + string(yamlByts))
}
```


### How to quickly Create Deployment?

If you want to quicky create Deployment, only input necessary Fields.

As shown in the code:

```
func howToNewDeployment() {
	dep, err := beku.NewDeployment().
		SetName("yulibaozi").SetLabels(map[string]string{"name": "yulibaozi"}).
		SetContainer("first", "mysql", 3307).SetContainer("second", "redis", 6379).SetHTTPLiveness(8080, "/metic", 30, 10, 10).Finish()
	if err != nil {
		panic(err)
	}
	yamlByts, err := beku.ToYAML(dep)
	if err != nil {
		panic(err)
	}
	fmt.Println("\n" + string(yamlByts))
}
```

### How to quickly Create Secret?

If you want to quicky create Secret, only input necessary Fields.

As shown in the code:
```
func howToNewSecret() {
	sec, err := beku.NewSecret().SetDataString(map[string]string{"key": "beku is very good!"}).SetNamespaceAndName("yulibaozi", "beku").Finish()
	if err != nil {
		panic(err)
	}
	yamlByts, err := beku.ToYAML(sec)
	if err != nil {
		panic(err)
	}
	fmt.Println("\n" + string(yamlByts))
}
```

### How to quickly Create ConfigMap(cm)?

If you want to quicky create cm, only input necessary Fields.

As shown in the code:
```
func howToNewCM() {
	sec, err := beku.NewCM().SetNamespaceAndName("yulibaozi", "beku").SetData(map[string]string{"key": "beku is very good!"}).Finish()
	if err != nil {
		panic(err)
	}
	yamlByts, err := beku.ToYAML(sec)
	if err != nil {
		panic(err)
	}
	fmt.Println("\n" + string(yamlByts))
}

```

### How to quickly Create DaemonSet(ds)?

If you want to quicky create ds, only input necessary Fields.

As shown in the code:
```
func howToNewDS() {
	ds, err := beku.NewDS().SetNamespaceAndName("roc", "api").SetContainer("first", "go-sdk", 8081).
		SetPodLabels(map[string]string{"name": "go-sdk"}).SetMinReadySeconds(10).
		SetPodQos("Guaranteed", true).SetHTTPLiveness(8081, "/health", 35, 10, 10, map[string]string{"userid": "1"}).SetContainer("second", "check", 8082).
		Finish()
	if err != nil {
		panic(err)
	}
	yamlByts, err := beku.ToYAML(ds)
	if err != nil {
		panic(err)
	}
	fmt.Println("\n" + string(yamlByts))
}
```

### How to quickly Create UnionPV?

If you want to quicky create UnionPV, only input necessary Fields.

As shown in the code:
```
func howToNewUnionPV() {
	pv, pvc, err := beku.NewUnionPV().SetNamespaceAndName("yulibaozi", "mysql-persistent").
		SetAccessMode(beku.ReadOnlyMany).SetCapacity(map[beku.ResourceName]string{beku.ResourceMemory: "5Gi"}).
		SetNFS(&beku.NFSVolumeSource{Server: "192.168.0.165:2314", Path: "/data"}).
		Finish()
	if err != nil {
		panic(err)
	}
	pvYamlByts, err := beku.ToYAML(pv)
	if err != nil {
		panic(err)
	}
	pvcYamlByts, err := beku.ToYAML(pvc)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n" + string(pvYamlByts))
	fmt.Println("---")
	fmt.Println("\n" + string(pvcYamlByts))
}
```

### How to quickly Create PersistentVolume(pv)?

If you want to quicky create pv, only input necessary Fields.

As shown in the code:
```
func howToNewPV() {
	pv, err := beku.NewPV().SetName("redis-pv").SetAccessMode(beku.ReadWriteMany).SetNFS(&beku.NFSVolumeSource{Server: "192.168.0.165:2132", Path: "/data"}).
		SetCapacity(map[beku.ResourceName]string{beku.ResourceEphemeralStorage: "50G"}).Finish()
	if err != nil {
		panic(err)
	}
	yamlByts, err := beku.ToYAML(pv)
	if err != nil {
		panic(err)
	}
	fmt.Println("\n" + string(yamlByts))
}
```

### How to quickly Create PersistentVolumeClaim(pvc)?

If you want to quicky create pvc, only input necessary Fields.

As shown in the code:
```
func howToNewPVC() {
	pvc, err := beku.NewPVC().SetNamespaceAndName("yulibaozi", "redis").SetAccessMode(beku.ReadOnlyMany).
		SetResourceRequests(map[beku.ResourceName]string{beku.ResourceEphemeralStorage: "5Gi"}).
		Finish()
	if err != nil {
		panic(err)
	}
	yamlByts, err := beku.ToYAML(pvc)
	if err != nil {
		panic(err)
	}
	fmt.Println("\n" + string(yamlByts))
}
```
