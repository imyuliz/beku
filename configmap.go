package beku

import (
	"encoding/json"
	"errors"

	"github.com/ghodss/yaml"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigMap include Kubernetes resource object ConfigMap(cm) and error.
type ConfigMap struct {
	cm  *v1.ConfigMap
	err error
}

// NewCM create ConfigMap(cm) and chain function call begin with this function.
func NewCM() *ConfigMap { return &ConfigMap{cm: &v1.ConfigMap{}} }

// Finish chain function call end with this function
// return real ConfigMap(really ConfigMap is Kubernetes resource object ConfigMap(cm) and error)
// In the function, it will check necessary parameters、input the default field。
func (obj *ConfigMap) Finish() (cm *v1.ConfigMap, err error) {
	obj.verify()
	return obj.cm, obj.err
}

// JSONNew use json data create ConfigMap
func (obj *ConfigMap) JSONNew(jsonbyts []byte) *ConfigMap {
	obj.error(json.Unmarshal(jsonbyts, obj.cm))
	return obj
}

// YAMLNew use yaml data create ConfigMap
func (obj *ConfigMap) YAMLNew(yamlbyts []byte) *ConfigMap {
	obj.error(yaml.Unmarshal(yamlbyts, obj.cm))
	return obj
}

// Replace replace cm by Kubernetes resource object
func (obj *ConfigMap) Replace(cm *v1.ConfigMap) *ConfigMap {
	if cm != nil {
		obj.cm = cm
	}
	return obj
}

// SetName set ConfigMap(cm) name
func (obj *ConfigMap) SetName(name string) *ConfigMap {
	obj.cm.SetName(name)
	return obj
}

// SetNamespace set CofigMap(cm) namespace, default namespace value is 'dafault'
func (obj *ConfigMap) SetNamespace(namespace string) *ConfigMap {
	obj.cm.SetNamespace(namespace)
	return obj
}

// SetNamespaceAndName set ConfigMap(cm) namespace and name
func (obj *ConfigMap) SetNamespaceAndName(namespace, name string) *ConfigMap {
	obj.cm.SetName(name)
	obj.cm.SetNamespace(namespace)
	return obj
}

// SetLabels set ConfigMap(cm) labels
func (obj *ConfigMap) SetLabels(labels map[string]string) *ConfigMap {
	obj.cm.SetLabels(labels)
	return obj
}

// SetData set ConfigMap(cm) data, map[key]value
func (obj *ConfigMap) SetData(data map[string]string) *ConfigMap {
	obj.cm.Data = data
	return obj
}

// Release release ConfigMap on Kubernetes
func (obj *ConfigMap) Release() (*v1.ConfigMap, error) {
	cm, err := obj.Finish()
	if err != nil {
		return nil, err
	}
	client, err := GetKubeClient()
	if err != nil {
		return nil, err
	}
	return client.CoreV1().ConfigMaps(cm.GetNamespace()).Create(cm)
}

// Apply  it will be updated when this resource object exists in K8s,
// it will be created when it does not exist.
func (obj *ConfigMap) Apply() (*v1.ConfigMap, error) {
	cm, err := obj.Finish()
	if err != nil {
		return nil, err
	}
	client, err := GetKubeClient()
	if err != nil {
		return nil, err
	}

	_, err = client.CoreV1().ConfigMaps(cm.GetNamespace()).Get(cm.GetName(), metav1.GetOptions{})
	if err != nil {
		return client.CoreV1().ConfigMaps(cm.GetNamespace()).Create(cm)
	}
	return client.CoreV1().ConfigMaps(cm.GetNamespace()).Update(cm)
}

func (obj *ConfigMap) error(err error) {
	if obj.err != nil {
		return
	}
	obj.err = err
}

// verify check ConfigMap necessary value,input default field.
func (obj *ConfigMap) verify() {
	if obj.err != nil {
		return
	}
	if !verifyString(obj.cm.Name) {
		obj.err = errors.New("ConfigMap name is not allowed to be empty")
		return
	}
	if len(obj.cm.Data) <= 0 {
		obj.err = errors.New("ConfigMap.Data is not allowed to be empty")
	}
	obj.cm.APIVersion = "v1"
	obj.cm.Kind = "ConfigMap"
}
