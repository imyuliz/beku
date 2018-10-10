package test

import (
	"testing"

	"github.com/yulibaozi/beku"
	"github.com/yulibaozi/beku/core"
)

// Test_Unionpv create unionpv, include pv and pvc
func Test_Unionpv(t *testing.T) {
	pv, pvc, err := beku.NewUnionPV().SetName("yulibaozi-test").SetVolumeMode(core.PersistentVolumeBlock).
		SetAccessMode(core.ReadWriteMany).SetNamespace("yulibaozi").SetCapacity(map[core.ResourceName]string{core.ResourceStorage: "5Gi"}).
		SetNFS(&core.NFSVolumeSource{Server: "10.141.40.141", Path: "/data"}).Finish()
	if err != nil {
		t.Fatal(err)
	}
	data, err := core.ToJSON(pv)
	if err != nil {
		t.Fatal(err)
	}
	data, err = core.ToJSON(pvc)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}

需要的字段
name,namespace,
