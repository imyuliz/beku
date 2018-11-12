package main

import (
	"fmt"

	"github.com/yulibaozi/beku"
	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type DepAndSvc struct {
	Namespace   string
	Name        string
	Labels      map[string]string
	Port        int32
	Image       string
	ImageName   string
	ServiceType beku.ServiceType
}

func main() {
	depAndSvcInfo := &DepAndSvc{
		Namespace:   "roc",
		Name:        "http",
		Labels:      map[string]string{"app": "http"},
		Port:        8081,
		Image:       "wucong60/kube-node-demo1:v1",
		ImageName:   "http",
		ServiceType: "NodePort",
	}
	//注册Kubernetes的ApiServer
	err := beku.RegisterK8sClient("http://192.168.0.183", "", "", "")
	if err != nil {
		panic(err)
	}
	//发布应用
	dp, svc, err := UnionDepAndSvc(depAndSvcInfo)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", svc)
	fmt.Printf("%+v", dp)
}

func UnionDepAndSvc(info *DepAndSvc) (dp *v1.Deployment, svc *corev1.Service, err error) {
	dp, err = beku.NewDeployment().SetNamespaceAndName(info.Namespace, info.Name).
		SetPodLabels(info.Labels).SetContainer(info.ImageName, info.Image, info.Port).Release()
	if err != nil {
		return
	}
	svc, err = beku.DeploymentToSvc(dp, info.ServiceType, true)
	if err != nil {
		return
	}
	return
}
