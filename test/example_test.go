package test

import (
	"fmt"

	"github.com/yulibaozi/beku"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func Example_beku_autoRelease() {
	err := beku.RegisterK8sClient("http://192.168.0.183", "", "", "")
	if err != nil {
		panic(err)
	}
	dp, err := beku.NewDeployment().SetNamespaceAndName("roc", "http").
		SetPodLabels(map[string]string{"app": "http"}).SetContainer("http", "wucong60/kube-node-demo1:v1", 8081).Release()
	if err != nil {
		panic(err)
	}
	svc, err := beku.DeploymentToSvc(dp, beku.ServiceTypeNodePort, true)
	if err != nil {
		panic(err)
	}
	byts, err := beku.ToYAML(svc)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n%s", string(byts))
	fmt.Println("over!")

}
func Example_beku_PushApp() {
	client := GetK8sClient()
	// cretate Namespace
	ns, err := beku.NewNs().SetName("roc").Finish()
	if err != nil {
		panic(err)
	}

	// push Namespace
	ns, err = client.CoreV1().Namespaces().Create(ns)
	if err != nil {
		panic(err)
	}
	fmt.Println("===")
	// create Deployment
	dp, err := beku.NewDeployment().SetNamespaceAndName("roc", "http").
		SetPodLabels(map[string]string{"app": "http"}).SetContainer("http", "wucong60/kube-node-demo1:v1", 8081).Finish()
	if err != nil {
		panic(err)
	}

	// push Deployment
	dp, err = client.AppsV1().Deployments(dp.GetNamespace()).Create(dp)
	if err != nil {
		panic(err)
	}
	fmt.Println("===")
	// create Service
	svc, err := beku.DeploymentToSvc(dp, beku.ServiceTypeNodePort)
	if err != nil {
		panic(err)
	}
	//push Service
	svc, err = client.CoreV1().Services(svc.GetNamespace()).Create(svc)
	if err != nil {
		panic(err)
	}
	byts, err := beku.ToYAML(svc)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n%s", string(byts))
	fmt.Println("over!")

}

// Example_beku_SetQos exmple how auto set Qos
func Example_beku_SetQos() {
	beku.RegisterResourceRequest(map[beku.ResourceName]string{beku.ResourceCPU: "200m", beku.ResourceMemory: "5G"})
	beku.RegisterResourceLimit(map[beku.ResourceName]string{beku.ResourceCPU: "200m", beku.ResourceMemory: "5G"})
	dep, err := beku.NewDeployment().SetContainer("http-demo", "wucong60/kube-node-demo1:v1", 8081).
		SetPodLabels(map[string]string{"name": "http-demo"}).SetNamespaceAndName("yulibaozi", "http-demo").SetPodQos("Guaranteed", true).Finish()
	if err != nil {
		panic(err)
	}
	yamlbyts, err := beku.ToYAML(dep)
	if err != nil {
		panic(err)
	}
	fmt.Println("\n" + string(yamlbyts))
	// create svc
	svc, err := beku.DeploymentToSvc(dep, beku.ServiceTypeNodePort)
	if err != nil {
		panic(err)
	}
	yamlbyts, err = beku.ToYAML(svc)
	fmt.Println("\n" + string(yamlbyts))
}

func Example_beku_NewSvc() {
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

// ExampleNewDeployment how to quickly create deployment example
func Example_beku_NewDeployment() {
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

func Example_beku_NewSecret() {
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

func Example_beku_NewCM() {
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

func Example_beku_NewSts() {
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

func Example_beku_NewDS() {
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

func Example_beku_NewUnionPV() {
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

func Example_beku_NewPVC() {
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

func Example_beku_NewPV() {
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

// GetK8sClient 获取k8s的链接
func GetK8sClient() *kubernetes.Clientset {
	cSet, err := kubernetes.NewForConfig(&rest.Config{
		Host: "http://192.168.0.184:6443/",
		TLSClientConfig: rest.TLSClientConfig{
			CAData:   []byte(ca),
			CertData: []byte(cert),
			KeyData:  []byte(key),
		},
	})
	if err != nil {
		panic(err)
	}
	return cSet
}

var ca = `ca`

var cert = `cert`

var key = `keydata`
