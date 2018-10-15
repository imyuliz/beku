package test

import (
	"testing"

	"github.com/yulibaozi/beku"
	"github.com/yulibaozi/beku/core"
)

func Test_CreateSecret(t *testing.T) {
	data, err := beku.NewSecret().SetNameSpaceAndName("yulibaozi", "rbd-secret").SetType(core.SecretTypeOpaque).SetDataString(map[string]string{"key": core.Base64Encode([]byte("AQAqfzNaofxHLBAAS7qY64uE/ddqWLOMVDhkYQ=="))}).Finish()
	if err != nil {
		panic(err)
	}

	result, err := core.ToJSON(data)
	if err != nil {
		panic(err)
	}

	t.Error(string(result))

}
