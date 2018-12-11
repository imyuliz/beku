package test

import (
	"fmt"

	"github.com/yulibaozi/beku"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	// certificate-authority-data
	caBase64 = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUIwakNDQVRPZ0F3SUJBZ0lKQUlBKzVnYVR4ZjZOTUFvR0NDcUdTTTQ5QkFNQ01CZ3hGakFVQmdOVkJBTU0KRFd0MVltVnlibVYwWlhNdFkyRXdIaGNOTVRnd05USTFNVEV3T1RFeldoY05Namd3TlRJeU1URXdPVEV6V2pBWQpNUll3RkFZRFZRUUREQTFyZFdKbGNtNWxkR1Z6TFdOaE1JR2JNQkFHQnlxR1NNNDlBZ0VHQlN1QkJBQWpBNEdHCkFBUUFnTzVvZE1IUGN4ZmVaRVgzZTE4NExpaUVSUjkyRWRvV0s4NUcwMWdTVEowR09wakFVbjd1bWVnUVZ0UVoKMWV2cXhQSCtuNFVMcGZSTUdBSjRxTWlzTmFzQUZaZ1dWQVo0eUxHWklMZEFxb2I0ZU5WTjVsZHRIUTE3blpHMApTUTBoVGpUa1ZPK1Z6SWljN0kvV1cwdVVyRnpqakw3MVJxNVdUc1g0akZXeEo0azFsVnlqSXpBaE1BOEdBMVVkCkV3RUIvd1FGTUFNQkFmOHdEZ1lEVlIwUEFRSC9CQVFEQWdLa01Bb0dDQ3FHU000OUJBTUNBNEdNQURDQmlBSkMKQVc0SnB5RGZUajJWbW5lMTJGWEZLUHF0SnBjamhHNm1nS1JNeGlQVWtMZ0U4SHEzbjFvNHJ4S1JOVzVHNVQ2WApBZitYdmZuRysramRMNkRXMmxJZjVYUk9Ba0lBNlo2bXlvZ1pwT29xUEdDTEZOb0xJWWJaRGlmWEZVaE9jVWVyClFKR1lIbjl4MStBR0hONWh5WExVQXM5M3ZkeFJ2TE5vajFXWDdNN0Q4R1NKL1pCUUM4Zz0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo="
	// client-certificate-data
	certBase64 = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNDVENDQVd1Z0F3SUJBZ0lKQUtLYXZVN3FrRkQyTUFvR0NDcUdTTTQ5QkFNQ01CZ3hGakFVQmdOVkJBTU0KRFd0MVltVnlibVYwWlhNdFkyRXdIaGNOTVRnd05USTFNVEV3T1RJMldoY05Namd3TlRJeU1URXdPVEkyV2pCQgpNU1l3SkFZRFZRUUREQjFyZFdKbExXRndhWE5sY25abGNpMXJkV0psYkdWMExXTnNhV1Z1ZERFWE1CVUdBMVVFCkNnd09jM2x6ZEdWdE9tMWhjM1JsY25Nd2dac3dFQVlIS29aSXpqMENBUVlGSzRFRUFDTURnWVlBQkFHdlhSTTUKMzRkUkdJMFlyV0hxcHZxaGVaeEJsbjVpTFd4azNWaFhmRjlqWkZucFBRdzBMclRNRXB6bEtXdzB1dnJsaWNmKwp6c1lTRHp4K1IvRXkwUVo2NXdBdkV2NU9KYzN1MHNINXArcjE4TSs0dlJ1WlUzOS9hRTB0U21OVDNmdTZSOGhpCmdVamlOS3VrMkErTHVxc25zbFlXV2ZxdGw2dU5mQW5lL0UvZmRiWW5XYU15TURBd0NRWURWUjBUQkFJd0FEQU8KQmdOVkhROEJBZjhFQkFNQ0JhQXdFd1lEVlIwbEJBd3dDZ1lJS3dZQkJRVUhBd0l3Q2dZSUtvWkl6ajBFQXdJRApnWXNBTUlHSEFrRm1INEJvS1d2NHBUY1h5V3JNQzBMaDlYbGNWQjloZ3dxd094V2RqLzdKNW9QbDRuaGlJb0RiCkkzZjYweFNZNTBTaStKTWtwSHFXbVlwTkxLZUh3TGY5MHdKQ0FMeTNtb2xaaEVBc24rTjFTZzVYMGxVYkxCV2IKTU1PN1BKZWdOUjBLaVQvQTN2QjA3UVRCcHJScTRJUnZ6MDBDZXJwUFkxeDArU1FCbTdHdE95YUgrWGtvCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"
	// client-key
	keyBase64 = "LS0tLS1CRUdJTiBFQyBQUklWQVRFIEtFWS0tLS0tCk1JSGNBZ0VCQkVJQVFMYTVLTUFmVDVicjdsc09MVlFtTlJ5RW96YTU1M2hRTGt2bGhhYXJLRVhQczQxVFhFblgKT3F3S1FiK01VaDhzL005Y1hmdFhrVjRmSHdHcTZFdUd5ZmFnQndZRks0RUVBQ09oZ1lrRGdZWUFCQUd2WFJNNQozNGRSR0kwWXJXSHFwdnFoZVp4QmxuNWlMV3hrM1ZoWGZGOWpaRm5wUFF3MExyVE1FcHpsS1d3MHV2cmxpY2YrCnpzWVNEengrUi9FeTBRWjY1d0F2RXY1T0pjM3Uwc0g1cCtyMThNKzR2UnVaVTM5L2FFMHRTbU5UM2Z1NlI4aGkKZ1VqaU5LdWsyQStMdXFzbnNsWVdXZnF0bDZ1TmZBbmUvRS9mZGJZbldRPT0KLS0tLS1FTkQgRUMgUFJJVkFURSBLRVktLS0tLQo="
)

func Example_beku_RegisterK8sClientBase64() {
	err := beku.RegisterK8sClientBase64("https://192.168.0.183:8080", caBase64, certBase64, keyBase64)
	if err != nil {
		panic(err)
	}
}

// certificate-authority-data
var ca = `-----BEGIN CERTIFICATE-----
MIIB0jCCATOgAwIBAgIJAIA+5gaTxf6NMAoGCCqGSM49BAMCMBgxFjAUBgNVBAMM
DWt1YmVybmV0ZXMtY2EwHhcNMTgwNTI1MTEwOTEzWhcNMjgwNTIyMTEwOTEzWjAY
MRYwFAYDVQQDDA1rdWJlcm5ldGVzLWNhMIGbMBAGByqGSM49AgEGBSuBBAAjA4GG
AAQAgO5odMHPcxfeZEX3e184LiiERR92EdoWK85G01gSTJ0GOpjAUn7umegQVtQZ
1evqxPH+n4ULpfRMGAJ4qMisNasAFZgWVAZ4yLGZILdAqob4eNVN5ldtHQ17nZG0
SQ0hTjTkVO+VzIic7I/WW0uUrFzjjL71Rq5WTsX4jFWxJ4k1lVyjIzAhMA8GA1Ud
EwEB/wQFMAMBAf8wDgYDVR0PAQH/BAQDAgKkMAoGCCqGSM49BAMCA4GMADCBiAJC
AW4JpyDfTj2Vmne12FXFKPqtJpcjhG6mgKRMxiPUkLgE8Hq3n1o4rxKRNW5G5T6X
Af+XvfnG++jdL6DW2lIf5XROAkIA6Z6myogZpOoqPGCLFNoLIYbZDifXFUhOcUer
QJGYHn9x1+AGHN5hyXLUAs93vdxRvLNoj1WX7M7D8GSJ/ZBQC8g=
-----END CERTIFICATE-----`

// client-certificate-data
var cert = `-----BEGIN CERTIFICATE-----
MIICCTCCAWugAwIBAgIJAKKavU7qkFD2MAoGCCqGSM49BAMCMBgxFjAUBgNVBAMM
DWt1YmVybmV0ZXMtY2EwHhcNMTgwNTI1MTEwOTI2WhcNMjgwNTIyMTEwOTI2WjBB
MSYwJAYDVQQDDB1rdWJlLWFwaXNlcnZlci1rdWJlbGV0LWNsaWVudDEXMBUGA1UE
CgwOc3lzdGVtOm1hc3RlcnMwgZswEAYHKoZIzj0CAQYFK4EEACMDgYYABAGvXRM5
34dRGI0YrWHqpvqheZxBln5iLWxk3VhXfF9jZFnpPQw0LrTMEpzlKWw0uvrlicf+
zsYSDzx+R/Ey0QZ65wAvEv5OJc3u0sH5p+r18M+4vRuZU39/aE0tSmNT3fu6R8hi
gUjiNKuk2A+LuqsnslYWWfqtl6uNfAne/E/fdbYnWaMyMDAwCQYDVR0TBAIwADAO
BgNVHQ8BAf8EBAMCBaAwEwYDVR0lBAwwCgYIKwYBBQUHAwIwCgYIKoZIzj0EAwID
gYsAMIGHAkFmH4BoKWv4pTcXyWrMC0Lh9XlcVB9hgwqwOxWdj/7J5oPl4nhiIoDb
I3f60xSY50Si+JMkpHqWmYpNLKeHwLf90wJCALy3molZhEAsn+N1Sg5X0lUbLBWb
MMO7PJegNR0KiT/A3vB07QTBprRq4IRvz00CerpPY1x0+SQBm7GtOyaH+Xko
-----END CERTIFICATE-----`

// client-key
var key = `-----BEGIN EC PRIVATE KEY-----
MIHcAgEBBEIAQLa5KMAfT5br7lsOLVQmNRyEoza553hQLkvlhaarKEXPs41TXEnX
OqwKQb+MUh8s/M9cXftXkV4fHwGq6EuGyfagBwYFK4EEACOhgYkDgYYABAGvXRM5
34dRGI0YrWHqpvqheZxBln5iLWxk3VhXfF9jZFnpPQw0LrTMEpzlKWw0uvrlicf+
zsYSDzx+R/Ey0QZ65wAvEv5OJc3u0sH5p+r18M+4vRuZU39/aE0tSmNT3fu6R8hi
gUjiNKuk2A+LuqsnslYWWfqtl6uNfAne/E/fdbYnWQ==
-----END EC PRIVATE KEY-----`

func Example_beku_RegisterK8sClient() {
	err := beku.RegisterK8sClient("https://192.168.0.183:8080", ca, cert, key)
	if err != nil {
		panic(err)
	}

}

func Example_beku_autoRelease() {
	err := beku.RegisterK8sClient("http://192.168.0.183:8080", "", "", "")
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
