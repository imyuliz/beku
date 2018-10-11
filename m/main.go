package main

import (
	"fmt"

	"github.com/yulibaozi/beku"
	"github.com/yulibaozi/beku/core"
)

func main() {
	pv, pvc, err := beku.NewUnionPV().SetName("yulibaozi-test").SetAccessMode(core.ReadWriteMany).SetCapacity(map[core.ResourceName]string{core.ResourceStorage: "5G"}).SetRBD(&core.RBDPersistentVolumeSource{
		CephMonitors: []string{
			"10.151.21.11:6789",
			"10.151.21.12:6789",
			"10.151.21.13:6789",
		},
		FSType:    "xfs",
		RBDPool:   "pool",
		RBDImage:  "xxx",
		RadosUser: "admin",
		Keyring:   "/etc/ceph/keyring",
		SecretRef: &core.SecretReference{
			Name:      "rbd-secret",
			Namespace: "xxx",
		},
	}).Finish()
	if err != nil {
		panic(err)
	}
	data, err := core.ToJSON(pv)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
	fmt.Println("===")
	data, err = core.ToJSON(pvc)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	// beku

	// fmt.Println("hello world")
	// pvc, err := beku.NewPVC().SetName("test-yulibaozi").SetLabels(map[string]string{"name": "test-yulibaozi"}).SetAccessMode(core.ReadWriteOnce).
	// 	SetVolumeMode(core.PersistentVolumeFilesystem).SetResourceLimit(map[core.ResourceName]string{core.ResourceMemory: "1Gi"}).
	// 	Finish()
	// if err != nil {

	// 	panic(err)
	// }
	// data, err := core.ToYAML(pvc)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(string(data))

	// fmt.Println("=========pv=========")

	// pv, err := beku.NewPV().SetName("test-pv").SetLabels(map[string]string{"name": "test-pv"}).SetNFS(&core.NFSVolumeSource{Server: "10.141.40.141", Path: "/data"}).SetAccessMode(core.RWX).SetCapacity(map[core.ResourceName]string{core.ResourceMemory: "1Gi"}).Finish()
	// if err != nil {
	// 	panic(err)
	// }
	// // &core.RBDPersistentVolumeSource{}
	// data, err = core.ToYAML(pv)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(string(data))
	//PV一定需要的字段
	// 1.name,labels 可以指定为name:xxx
	// 2.需要限制accessModes可以是多选
	// 3. 相关资源限制,比如storage 比如5GB
	// 4. 后端支持的存储方式  nfs,rbd
	// 	if nfs
	// 		需要服务器地址（Server）
	// 		Path:服务器的路径
	// 		是否只读

	// 	if rbd
	// 		CephMonitors
	// 		RBDImage
	// 		FSType ("ext4", "xfs", "ntfs".)
	// 		RBDPool（rbd）
	// 		User

	// 		Keyring string `json:"keyring,omitempty" protobuf:"bytes,6,opt,name=keyring"`

	// m1 := map[string]string{"name": "yulibaozi"}
	// m2 := map[string]string{"name": "yulibaozi"}
	// // var m1, m2 map[string]string
	// fmt.Println(reflect.DeepEqual(m1, m2))

}
