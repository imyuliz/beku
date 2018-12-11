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
	t.Log(string(data))
}

func Test_DeploymentPreStop(t *testing.T) {
	dep, err := beku.NewDeployment().SetSelector(map[string]string{"name": "live"}).SetNamespaceAndName("live", "mysql").SetContainer("mysql", "mysql", 8080).SetPreStopHTTP(beku.URISchemeHTTP, "127.0.0.1", 8080, "/init").Finish()
	if err != nil {
		t.Fatal(err)
	}
	data, err := beku.ToYAML(dep)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}

func Test_DeploymentPostStart(t *testing.T) {
	dep, err := beku.NewDeployment().SetSelector(map[string]string{"name": "devfeel"}).SetNamespaceAndName("devfeel", "mysql").SetContainer("mysql", "mysql", 8080).SetPostStartExec([]string{"bash", "c", `echo "hello world"`}).Finish()
	if err != nil {
		t.Fatal(err)
	}
	data, err := beku.ToYAML(dep)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}
