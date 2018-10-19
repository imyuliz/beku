package beku

import (
	"encoding/json"
	"errors"

	"github.com/ghodss/yaml"
	"k8s.io/api/core/v1"
)

// Secret include Kuebernetes resource object Secret and error.
type Secret struct {
	sc  *v1.Secret
	err error
}

// NewSecret create Secret and chain function call begin with this function.
func NewSecret() *Secret { return &Secret{sc: &v1.Secret{}} }

// Finish chain function call end with this function.
// return obj(Kubernetes resource object) and error
// In the function, it will check necessary parameters、input the default field。
func (obj *Secret) Finish() (*v1.Secret, error) {
	obj.verify()
	return obj.sc, obj.err
}

// JSONNew use json data create Secret
func (obj *Secret) JSONNew(jsonbyts []byte) *Secret {
	obj.err = json.Unmarshal(jsonbyts, obj.sc)
	return obj
}

// YAMLNew use yaml data create Secret
func (obj *Secret) YAMLNew(yamlbyts []byte) *Secret {
	obj.err = yaml.Unmarshal(yamlbyts, obj.sc)
	return obj
}

// SetName set Secret name
func (obj *Secret) SetName(name string) *Secret {
	obj.sc.SetName(name)
	return obj
}

// SetNamespace set Secret namespace ,default namespace is 'default'
func (obj *Secret) SetNamespace(namespace string) *Secret {
	obj.sc.SetNamespace(namespace)
	return obj
}

// SetNamespaceAndName set Secret namespace and name
func (obj *Secret) SetNamespaceAndName(namespace, name string) *Secret {
	obj.sc.SetNamespace(namespace)
	obj.sc.SetName(name)
	return obj
}

// SetLabels set Secret labels
func (obj *Secret) SetLabels(labels map[string]string) *Secret {
	obj.sc.SetLabels(labels)
	return obj
}

// SetDataString set Secret data, and Don't need to encode base64
func (obj *Secret) SetDataString(datas map[string]string) *Secret {
	obj.sc.StringData = datas
	return obj
}

// SetDataBytes set Secret data for byte,and Don't need to encode base64
func (obj *Secret) SetDataBytes(bytes map[string][]byte) *Secret {
	obj.sc.Data = bytes
	return obj
}

// SetType set Secret type,have Opaque and kubernetes.io/service-account-token
// Opaque user-defined data
// kubernetes.io/service-account-token is used to kubernetes apiserver,because apiserver need to auth
func (obj *Secret) SetType(secType SecretType) *Secret {
	obj.sc.Type = secType.ToK8s()
	return obj
}

// verify check Secret necessary value, input the default field and input related data.
func (obj *Secret) verify() {
	if obj.err != nil {
		return
	}
	if !verifyString(obj.sc.Name) {
		obj.err = errors.New("secret name not allow empty")
		return
	}
	obj.sc.APIVersion = "v1"
	obj.sc.Kind = "Secret"
}
