package test

import (
	"testing"

	"github.com/yulibaozi/beku"
)

// Test_PVCCreate create pvc
func Test_PVCCreate(t *testing.T) {
	data, err := beku.NewPVC().SetName("yulibaozi-pv").SetAccessMode(beku.ReadWriteMany).SetVolumeMode(beku.PersistentVolumeBlock).SetResourceRequests(map[beku.ResourceName]string{beku.ResourceStorage: "5Gi"}).Finish()
	if err != nil {
		t.Fatal(err)
	}
	databyts, err := beku.ToJSON(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(databyts))

}
