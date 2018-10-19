package test

import (
	"fmt"
	"testing"

	"github.com/yulibaozi/beku"
)

// Test_Unionpv create unionpv, include pv and pvc
func Test_Unionpv(t *testing.T) {
	pv, pvc, err := beku.NewUnionPV().SetName("yulibaozi-test").SetVolumeMode(beku.PersistentVolumeBlock).
		SetAccessMode(beku.ReadWriteMany).SetNamespace("yulibaozi").SetCapacity(map[beku.ResourceName]string{beku.ResourceStorage: "5Gi"}).
		SetNFS(&beku.NFSVolumeSource{Server: "xx.xx.xx.xxx", Path: "/data"}).Finish()
	if err != nil {
		t.Fatal(err)
	}
	data, err := beku.ToJSON(pv)
	if err != nil {
		t.Fatal(err)
	}
	data, err = beku.ToJSON(pvc)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))

	pv, pvc, err = beku.NewUnionPV().SetName("yulibaozi-test").SetAccessMode(beku.ReadWriteMany).SetCapacity(map[beku.ResourceName]string{beku.ResourceStorage: "5G"}).
		SetRBD(&beku.RBDPersistentVolumeSource{}).Finish()
	if err != nil {
		panic(err)
	}
	data1, err := beku.ToJSON(pv)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data1))
	fmt.Println("===")
	data2, err := beku.ToJSON(pvc)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data2))

	// beku

	// fmt.Println("hello world")
	// pvc, err := beku.NewPVC().SetName("test-yulibaozi").SetLabels(map[string]string{"name": "test-yulibaozi"}).SetAccessMode(ReadWriteOnce).
	// 	SetVolumeMode(PersistentVolumeFilesystem).SetResourceLimit(map[ResourceName]string{ResourceMemory: "1Gi"}).
	// 	Finish()
	// if err != nil {

	// 	panic(err)
	// }
	// data, err := ToYAML(pvc)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(string(data))

	// fmt.Println("=========pv=========")

	// pv, err := beku.NewPV().SetName("test-pv").SetLabels(map[string]string{"name": "test-pv"}).SetNFS(&NFSVolumeSource{Server: "10.141.40.141", Path: "/data"}).SetAccessMode(RWX).SetCapacity(map[ResourceName]string{ResourceMemory: "1Gi"}).Finish()
	// if err != nil {
	// 	panic(err)
	// }
	// // &RBDPersistentVolumeSource{}
	// data, err = ToYAML(pv)
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
