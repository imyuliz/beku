package test

import (
	"testing"

	"github.com/yulibaozi/beku"
	"github.com/yulibaozi/beku/core"
)

// Test_PVCCreate create pvc
func Test_PVCCreate(t *testing.T) {
	data, err := beku.NewPVC().SetName("yulibaozi-pv").SetAccessMode(core.ReadWriteMany).SetVolumeMode(core.PersistentVolumeBlock).SetResourceRequests(map[core.ResourceName]string{core.ResourceStorage: "5Gi"}).Finish()
	if err != nil {
		t.Fatal(err)
	}
	databyts, err := core.ToJSON(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(databyts))
}
