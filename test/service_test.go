package test

import (
	"testing"

	"github.com/yulibaozi/beku"
)

func Test_CreateSvc(t *testing.T) {
	svc, err := beku.NewSvc().SetNamespaceAndName("yulibaozi", "mysql-svc").SetSelector(map[string]string{"app": "mysql"}).
		SetServiceType(beku.ServiceTypeNodePort).SetPorts([]beku.ServicePort{beku.ServicePort{
		Name: "mysql",
		Port: 3306,
	}}).Finish()
	if err != nil {
		t.Fatal(err)
	}
	data, err := beku.ToJSON(svc)
	if err != nil {
		t.Fatal(err)
	}
	t.Error(string(data))
}
