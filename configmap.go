package beku

import (
	"errors"

	"k8s.io/api/core/v1"
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
