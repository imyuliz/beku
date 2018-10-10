package test

import (
	"testing"

	"github.com/yulibaozi/beku"
	"github.com/yulibaozi/beku/core"
)

// Test_CreatePV create pv
func Test_CreatePV(t *testing.T) {
	data, err := beku.NewPV().SetName("yulibaozi").SetCapacity(map[core.ResourceName]string{core.ResourceStorage: "5Gi"}).SetLabels(map[string]string{"name": "yulibaozi"}).SetAccessMode(core.ReadWriteOnce).SetNFS(&core.NFSVolumeSource{Server: "127.0.0.1", Path: "/data"}).Finish()
	if err != nil {
		t.Fatal(err)
	}
	databyts, err := core.ToYAML(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(databyts))
}
