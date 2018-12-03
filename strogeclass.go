package beku

import (
	"encoding/json"
	"errors"

	"github.com/ghodss/yaml"
	"k8s.io/api/storage/v1"
)

// StorageClass include Kubernetes resource object StorageClass and error.
type StorageClass struct {
	sc  *v1.StorageClass
	err error
}

// NewStorageClass create StorageClass and chain function call begin with this function.
func NewStorageClass() *StorageClass { return &StorageClass{sc: &v1.StorageClass{}} }

// Finish chain function call end with this function
// return Kubernetes resource object StorageClass and error.
// In the function, it will check necessary parameters,input the default field.
func (obj *StorageClass) Finish() (*v1.StorageClass, error) {
	obj.verify()
	return obj.sc, obj.err
}

// JSONNew use json data create StorageClass
func (obj *StorageClass) JSONNew(jsonbyte []byte) *StorageClass {
	obj.error(json.Unmarshal(jsonbyte, obj.sc))
	return obj
}

// YAMLNew use yaml data create StorageClass
func (obj *StorageClass) YAMLNew(yamlbyts []byte) *StorageClass {
	obj.error(yaml.Unmarshal(yamlbyts, obj.sc))
	return obj
}

// SetName set storageCLASS name
func (obj *StorageClass) SetName() *StorageClass { return obj }

// SetProvisioner set storageClass privisioner
func (obj *StorageClass) SetProvisioner(provisioner string) *StorageClass {
	obj.sc.Provisioner = provisioner
	return obj
}

// SetParameters set storageClass parameters
func (obj *StorageClass) SetParameters(parameters map[string]string) *StorageClass {
	obj.sc.Parameters = parameters
	return obj
}

// SetReclaimPolicy set setReclaim policy
func (obj *StorageClass) SetReclaimPolicy(reclaimPolicy PersistentVolumeReclaimPolicy) *StorageClass {
	policy := reclaimPolicy.ToK8s()
	obj.sc.ReclaimPolicy = &policy
	return obj
}

// SetMountOptions set sorageClass mount Options
func (obj *StorageClass) SetMountOptions(opts []string) *StorageClass {
	obj.sc.MountOptions = opts
	return obj
}

// SetVolumeBindingMode set storageClass Volume BindingMode,value only:WaitForFirstConsumer,Immediate
func (obj *StorageClass) SetVolumeBindingMode(bindingMode VolumeBindingMode) *StorageClass {
	obj.sc.VolumeBindingMode = bindingMode.ToK8s()
	return obj
}

// SetAnnotations set storageClass annotations
func (obj *StorageClass) SetAnnotations(annotations map[string]string) *StorageClass {
	obj.sc.SetAnnotations(annotations)
	return obj
}

// SetLabels set StorageClass labels
func (obj *StorageClass) SetLabels(labels map[string]string) *StorageClass {
	obj.sc.SetLabels(labels)
	return obj
}

func (obj *StorageClass) verify() {
	if obj.err != nil {
		return
	}
	if obj.sc.Provisioner == "" {
		obj.err = errors.New("StorageClass.spec.Provisioner is not allowed to be empty")
		return
	}
	obj.sc.APIVersion = "storage.k8s.io/v1"
	obj.sc.Kind = "StorageClass"
	return
}
func (obj *StorageClass) error(err error) {
	if obj.err != nil {
		return
	}
	obj.err = err
}
