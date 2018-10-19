package test

import (
	"testing"

	"github.com/yulibaozi/beku"
)

func Test_CreateSecret(t *testing.T) {
	data, err := beku.NewSecret().SetNamespaceAndName("yulibaozi", "rbd-secret").SetType(beku.SecretTypeOpaque).SetDataString(map[string]string{"key": beku.Base64Encode([]byte("AQAqfzNaofxHLBAAS7qY64uE/ddqWLOMVDhkYQ=="))}).Finish()
	if err != nil {
		panic(err)
	}

	result, err := beku.ToJSON(data)
	if err != nil {
		panic(err)
	}

	t.Error(string(result))

}
