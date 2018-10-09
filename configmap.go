package beku

import (
	"errors"

	"k8s.io/api/core/v1"
)

// ConfigMap include k8s resource object ConfigMap and error
type ConfigMap struct {
	cm  *v1.ConfigMap
	err error
}

// NewCM create configMap
func NewCM() *ConfigMap {
	return &ConfigMap{
		cm: &v1.ConfigMap{},
	}
}

// SetName set configMap name
func (obj *ConfigMap) SetName(name string) *ConfigMap {
	obj.cm.SetName(name)
	return obj
}

// SetNamespace  set cm namespace ,default namespace dafault
func (obj *ConfigMap) SetNamespace(namespace string) *ConfigMap {
	obj.cm.SetNamespace(namespace)
	return obj
}

// SetNamespaceAndName set configMap namespace and name
func (obj *ConfigMap) SetNamespaceAndName(namespace, name string) *ConfigMap {
	obj.cm.SetName(name)
	obj.cm.SetNamespace(namespace)
	return obj
}

// SetLabels set configMap labels,only map
func (obj *ConfigMap) SetLabels(labels map[string]string) *ConfigMap {
	obj.cm.SetLabels(labels)
	return obj
}

// SetData set configMap data, map[key]value
func (obj *ConfigMap) SetData(data map[string]string) *ConfigMap {
	obj.cm.Data = data
	return obj
}

// verify check configMap integrity
func (obj *ConfigMap) verify() {
	if obj.err != nil {
		return
	}
	if !verifyString(obj.cm.Name) {
		obj.err = errors.New("configmap name not allow empty")
		return
	}
	obj.cm.APIVersion = "v1"
	obj.cm.Kind = "ConfigMap"
}

// Finish the final step, verify configMap and return real configMap(cm)
func (obj *ConfigMap) Finish() (*v1.ConfigMap, error) {
	obj.verify()
	if obj.err != nil {
		return nil, obj.err
	}
	return obj.cm, nil
}
