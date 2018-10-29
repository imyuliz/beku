package beku

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ghodss/yaml"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PersistentVolumeClaim include kubernetes resource object PersistentVolumeClaim(pvc) and error.
type PersistentVolumeClaim struct {
	pvc *v1.PersistentVolumeClaim
	err error
}

// NewPVC create PersistentVolumeClaim(pvc) and chain function call begin with this function.
func NewPVC() *PersistentVolumeClaim { return &PersistentVolumeClaim{pvc: &v1.PersistentVolumeClaim{}} }

// Finish Chain function call end with this function
// return Kubernetes resource object PersistentVolumeClaim(pvc) and error.
// In the function, it will check necessary parameters?input the default field?
func (obj *PersistentVolumeClaim) Finish() (*v1.PersistentVolumeClaim, error) {
	obj.verify()
	return obj.pvc, obj.err
}

// JSONNew use json data create PersistentVolumeClaim(pvc)
func (obj *PersistentVolumeClaim) JSONNew(jsonbyts []byte) *PersistentVolumeClaim {
	obj.error(json.Unmarshal(jsonbyts, obj.pvc))
	return obj
}

// YAMLNew use yaml data create PersistentVolumeClaim(pvc)
func (obj *PersistentVolumeClaim) YAMLNew(yamlbyts []byte) *PersistentVolumeClaim {
	obj.error(yaml.Unmarshal(yamlbyts, obj.pvc))
	return obj
}

// Replace replace PersistentVolumeClaim by Kubernetes resource object
func (obj *PersistentVolumeClaim) Replace(pvc *v1.PersistentVolumeClaim) *PersistentVolumeClaim {
	if pvc != nil {
		obj.pvc = pvc
	}
	return obj
}

// SetName set PersistentVolumeClaim(pvc) name
func (obj *PersistentVolumeClaim) SetName(name string) *PersistentVolumeClaim {
	obj.pvc.SetName(name)
	return obj
}

// SetNamespaceAndName set Deployment namespace,set Pod namespace,set Deployment name.
func (obj *PersistentVolumeClaim) SetNamespaceAndName(namespace, name string) *PersistentVolumeClaim {
	obj.SetNamespace(namespace)
	obj.SetName(name)
	return obj
}

// GetName get PersistentVolumeClaim(pvc) name
func (obj *PersistentVolumeClaim) GetName() string {
	return obj.pvc.GetName()
}

// GetNamespace get PersistentVolumeClaim(pvc) namespace
func (obj *PersistentVolumeClaim) GetNamespace() string {
	return obj.pvc.GetNamespace()
}

// SetNamespace set PersistentVolumeClaim(pvc) namespace,default namespace is 'default'
func (obj *PersistentVolumeClaim) SetNamespace(namespace string) *PersistentVolumeClaim {
	obj.pvc.SetNamespace(namespace)
	return obj
}

// SetLabels set PersistentVolumeClaim(pvc) labels
func (obj *PersistentVolumeClaim) SetLabels(labels map[string]string) *PersistentVolumeClaim {
	obj.pvc.SetLabels(labels)
	return obj
}

// GetLabels get PersistentVolumeClaim(pvc) labels
func (obj *PersistentVolumeClaim) GetLabels() map[string]string {
	return obj.pvc.GetLabels()
}

// SetAnnotations set PersistentVolumeClaim(pvc) annotations
func (obj *PersistentVolumeClaim) SetAnnotations(annotations map[string]string) *PersistentVolumeClaim {
	obj.pvc.SetAnnotations(annotations)
	return obj
}

// SetAccessMode set PersistentVolumeClaim(pvc) access mode, only one
func (obj *PersistentVolumeClaim) SetAccessMode(mode PersistentVolumeAccessMode) *PersistentVolumeClaim {
	obj.pvc.Spec.AccessModes = []v1.PersistentVolumeAccessMode{mode.ToK8s()}
	return obj
}

// SetAccessModes set PersistentVolumeClaim(pvc) accessModes, many modes
func (obj *PersistentVolumeClaim) SetAccessModes(modes []PersistentVolumeAccessMode) *PersistentVolumeClaim {
	var objModes []v1.PersistentVolumeAccessMode
	for _, m := range modes {
		objModes = append(objModes, m.ToK8s())
	}
	obj.pvc.Spec.AccessModes = objModes
	return obj
}

// SetVolumeMode PersistentVolumeClaim(pvc) vloume mode,have Block and Filesystem mode
func (obj *PersistentVolumeClaim) SetVolumeMode(volumeMode PersistentVolumeMode) *PersistentVolumeClaim {
	m := volumeMode.ToK8s()
	if m == nil {
		obj.error(fmt.Errorf("SetVolumeMode err: the volumeMode: %v is not allowed", volumeMode))
		return obj
	}
	obj.pvc.Spec.VolumeMode = m
	return obj
}

// SetResourceLimits set PersistentVolumeClaim(pvc) resource limits
func (obj *PersistentVolumeClaim) SetResourceLimits(limits map[ResourceName]string) *PersistentVolumeClaim {
	data, err := ResourceMapsToK8s(limits)
	if err != nil {
		obj.error(fmt.Errorf("SetResourceLimit err:%v", err))
		return obj
	}
	obj.pvc.Spec.Resources.Limits = data
	return obj
}

// SetResourceRequests set PersistentVolumeClaim(pvc) reource requests
func (obj *PersistentVolumeClaim) SetResourceRequests(requests map[ResourceName]string) *PersistentVolumeClaim {
	data, err := ResourceMapsToK8s(requests)
	if err != nil {
		obj.error(fmt.Errorf("SetResourceRequests err:%v", err))
		return obj
	}
	obj.pvc.Spec.Resources.Requests = data
	return obj
}

// SetStorageClassName set PersistentVolumeClaim(pvc) storageclasss name
func (obj *PersistentVolumeClaim) SetStorageClassName(classname string) *PersistentVolumeClaim {
	if classname == "" || len(classname) <= 0 {
		obj.error(errors.New("SetStorageClassName err, StorageClassName is not allowed to be empty"))
		return obj
	}
	obj.pvc.Spec.StorageClassName = &classname
	return obj
}

// SetSelector set PersistentVolumeClaim(pvc) selector
func (obj *PersistentVolumeClaim) SetSelector(labels map[string]string) *PersistentVolumeClaim {
	if len(labels) < 1 {
		obj.error(errors.New("SetSelector error, labels is not allowed to be empty"))
		return obj
	}
	if obj.pvc.Spec.Selector == nil {
		selector := &metav1.LabelSelector{
			MatchLabels: labels,
		}
		obj.pvc.Spec.Selector = selector
		return obj
	}
	obj.pvc.Spec.Selector.MatchLabels = labels
	return obj
}

// GetSelector get PersistentVolumeClaim(pvc) selector
func (obj *PersistentVolumeClaim) GetSelector() map[string]string {
	if obj.pvc.Spec.Selector == nil {
		return nil
	}
	return obj.pvc.Spec.Selector.MatchLabels
}

// SetMatchExpressions set Deployment match expressions
// the field is used to set complicated Label.
func (obj *PersistentVolumeClaim) SetMatchExpressions(ents []LabelSelectorRequirement) *PersistentVolumeClaim {
	requirements := make([]metav1.LabelSelectorRequirement, 0)
	for index := range ents {
		requirements = append(requirements, metav1.LabelSelectorRequirement{
			Key:      ents[index].Key,
			Operator: metav1.LabelSelectorOperator(ents[index].Operator),
			Values:   ents[index].Values,
		})
	}
	if obj.pvc.Spec.Selector == nil {
		obj.pvc.Spec.Selector = &metav1.LabelSelector{
			MatchExpressions: requirements,
		}
		return obj
	}
	obj.pvc.Spec.Selector.MatchExpressions = requirements
	return obj
}

func (obj *PersistentVolumeClaim) error(err error) {
	if obj.err != nil {
		return
	}
	obj.err = err
	return
}

// verify check service necessary value, input the default field and input related data.
func (obj *PersistentVolumeClaim) verify() {
	if obj.err != nil {
		return
	}
	if !verifyString(obj.pvc.GetName()) {
		obj.err = errors.New("pvc name is not allowed to be empty")
		return
	}
	if obj.pvc.Spec.AccessModes == nil || len(obj.pvc.Spec.AccessModes) < 1 {
		obj.err = errors.New("pvc accessModes is not allowed to be empty")
		return
	}
	if obj.pvc.Spec.Resources.Limits == nil && obj.pvc.Spec.Resources.Requests == nil {
		obj.err = errors.New("both limits and requests is empty not allowed")
		return
	}
	obj.pvc.Kind = "PersistentVolumeClaim"
	obj.pvc.APIVersion = "v1"
}
