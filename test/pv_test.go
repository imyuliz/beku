package test

import (
	"testing"

	"github.com/yulibaozi/beku"
)

// Test_CreatePV create pv
func Test_CreatePV(t *testing.T) {
	data, err := beku.NewPV().SetName("yulibaozi").SetCapacity(map[beku.ResourceName]string{beku.ResourceStorage: "5Gi"}).SetLabels(map[string]string{"name": "yulibaozi"}).SetAccessMode(beku.ReadWriteOnce).SetNFS(&beku.NFSVolumeSource{Server: "127.0.0.1", Path: "/data"}).Finish()
	if err != nil {
		t.Fatal(err)
	}
	databyts, err := beku.ToYAML(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(databyts))
}
