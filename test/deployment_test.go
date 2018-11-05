package test

import (
	"testing"

	"github.com/yulibaozi/beku"
)

func Test_DeploymentCreate(t *testing.T) {
	dep, err := beku.NewDeployment().SetNamespace("litest").
		SetName("mysql").SetSelector(map[string]string{"app": "mysql"}).
		SetContainer("mysql", "mysql:5.6", 3306).SetEnvs(map[string]string{"MYSQL_ROOT_PASSWORD": "password"}).
		SetPVClaim("yulitest", "yulipv-455d130f").SetPVCMounts("yulitest", "/var/lib/mysql").Finish()
	if err != nil {
		t.Fatal(err)
	}
	data, err := beku.ToYAML(dep)
	if err != nil {
		t.Fatal(err)
	}
	t.Error(string(data))
}
