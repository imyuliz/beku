package core

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/ghodss/yaml"
	"k8s.io/api/core/v1"
	apiresource "k8s.io/apimachinery/pkg/api/resource"
)

// ToYAML struct translation into yaml
func ToYAML(o interface{}) (byts []byte, err error) {
	byts, err = yaml.Marshal(o)
	return
}

// ToJSON struct translation into json
func ToJSON(v interface{}) (byts []byte, err error) {
	byts, err = json.Marshal(v)
	return
}

// JSONToYAML json data translation into yaml
func JSONToYAML(jbyts []byte) (ybyts []byte, err error) {
	ybyts, err = yaml.JSONToYAML(jbyts)
	return
}

// YAMLToJSON yaml data translation into json
func YAMLToJSON(ybyts []byte) (jbyts []byte, err error) {
	jbyts, err = yaml.YAMLToJSON(ybyts)
	return
}

// Base64Encode base64 编码
func Base64Encode(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

// Base64Decode base64 解码
func Base64Decode(src string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(src)
}

// ResourceMapsToK8s to K8s resourceList
func ResourceMapsToK8s(maps map[ResourceName]string) (v1.ResourceList, error) {
	data := make(v1.ResourceList, 0)
	for k, v := range maps {
		q, err := apiresource.ParseQuantity(v)
		if err != nil {
			return nil, err
		}
		reName := k.ToK8s()
		if reName == "" {
			return nil, errors.New("resource name not allow")
		}
		data[reName] = q
	}
	if len(data) < 1 {
		return nil, errors.New("source cann't allow empty")
	}
	return data, nil
}
