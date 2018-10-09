package beku

import (
	"errors"
	"fmt"

	"github.com/yulibaozi/beku/core"
	"github.com/yulibaozi/mapper"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PersistentVolumeClaim include kubernetes resource object
type PersistentVolumeClaim struct {
	pvc  *v1.PersistentVolumeClaim
	err error
}

// Newobj create pvc
func NewPVC() *PersistentVolumeClaim {
	return &PersistentVolumeClaim{
		pvc: &v1.PersistentVolumeClaim{},
	}
}

// SetNametae set pvc name
func (obj *PersistentVolumeClaim) SetName(name string) *PersistentVolumeClaim {
	obj.pvc.SetName(name)
	return obj
}

// SetNameSpace set  pvc namespace,default namespace is default
func (obj *PersistentVolumeClaim) SetNameSpace(namespace string) *PersistentVolumeClaim {
	obj.pvc.SetNamespace(namespace)
	return obj
}

// SetLabels set pvc label
func (obj *PersistentVolumeClaim) SetLabels(labels map[string]string) *PersistentVolumeClaim {
	obj.pvc.SetLabels(labels)
	return obj
}

// SetAnnotations set annotation
func (obj *PersistentVolumeClaim) SetAnnotations(annotations map[string]string) *PersistentVolumeClaim {
	obj.pvc.SetAnnotations(annotations)
	return obj
}

// SetobjAccessMode set pvc accessMode
func (obj *PersistentVolumeClaim) SetPVCAccessMode(mode core.PersistentVolumeAccessMode) *PersistentVolumeClaim {
	obj.pvc.Spec.AccessModes = []v1.PersistentVolumeAccessMode{mode.ToK8s()}
	return obj
}

// SetobjAccessModes set pvc accessModes
func (obj *PersistentVolumeClaim) SetPVCAccessModes(modes []core.PersistentVolumeAccessMode) *PersistentVolumeClaim {
	var objModes []v1.PersistentVolumeAccessMode
	for _, m := range modes {
		objModes = append(objModes, m.ToK8s())
	}
	obj.pvc.Spec.AccessModes = objModes
	return obj
}

// SetVolumeMode pvc vloume mode,have Block and Filesystem  mode
func (obj *PersistentVolumeClaim) SetVolumeMode(volumeMode core.PersistentVolumeMode) *PersistentVolumeClaim {
	m := volumeMode.ToK8s()
	if m == nil {
		obj.err = fmt.Errorf("set volumeMode err: this volumeMode not allow %v", volumeMode)
	}
	obj.pvc.Spec.VolumeMode = m
	return obj
}

// SetResourceLimit set pvc resource limit
func (obj *PersistentVolumeClaim) SetResourceLimit(limits map[core.ResourceName]string) *PersistentVolumeClaim {
	data, err := core.ResourceMapsToK8s(limits)

	if err != nil {
		obj.err = fmt.Errorf("limit set err:%v", err)
		return obj
	}
	obj.pvc.Spec.Resources.Limits = data
	return obj
}

// SetResourceRequests set pvc reource requests
func (obj *PersistentVolumeClaim) SetResourceRequests(requests map[core.ResourceName]string) *PersistentVolumeClaim {
	data, err := core.ResourceMapsToK8s(requests)
	if err != nil {
		obj.err = fmt.Errorf("request set err:%v", err)
		return obj
	}
	obj.pvc.Spec.Resources.Requests = data
	return obj
}

// SetStorageClassName set storageclasss name
func (obj *PersistentVolumeClaim) SetStorageClassName(classname string) *PersistentVolumeClaim {
	if classname == "" || len(classname) <= 0 {
		obj.err = errors.New("set StorageClassName is empty,set failed")
		return obj
	}
	obj.pvc.Spec.StorageClassName = &classname
	return obj
}

// SetSelector set pvc selector
func (obj *PersistentVolumeClaim) SetSelector(labels map[string]string) *PersistentVolumeClaim {
	if len(labels) < 1 {
		obj.err = errors.New("set LabelSelector error, labels is empty")
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

// SetMatchExpressions set pvc label selector,have key,operator and values
func (obj *PersistentVolumeClaim) SetMatchExpressions(ents []core.LabelSelectorRequirement) *PersistentVolumeClaim {
	requirements := make([]metav1.LabelSelectorRequirement, 0)
	err := mapper.AutoMapper(ents, requirements)
	if err != nil {
		obj.err = fmt.Errorf("SetMatchExpressions error:%v", err)
		return obj
	}
	if obj.pvc.Spec.Selector == nil {
		selector := &metav1.LabelSelector{
			MatchExpressions: requirements,
		}
		obj.pvc.Spec.Selector = selector
		return obj
	}
	obj.pvc.Spec.Selector.MatchExpressions = requirements
	return obj
}

// verify  pvc
func (obj *PersistentVolumeClaim) verify() {
	if !verifyString(obj.pvc.GetName()) {
		obj.err = errors.New("obj name not allow empty")
		return
	}
	if obj.pvc.Spec.AccessModes == nil || len(obj.pvc.Spec.AccessModes) < 1 {
		obj.err = errors.New("obj accessModes not allow empty")
		return
	}
	if obj.pvc.Spec.VolumeMode == nil {
		obj.err = errors.New("obj volumeMode not allow nil")
		return
	}
	if obj.pvc.Spec.Resources.Limits == nil && obj.pvc.Spec.Resources.Requests == nil {
		obj.err = errors.New("both limits and requests is nil  not allow")
		return
	}
	obj.pvc.Kind = "PersistentVolumeClaim"
	obj.pvc.APIVersion = "v1"
}

// Finish  the final step,will return kubernetes resource object pv and error
func (obj *PersistentVolumeClaim) Finish() (*v1.PersistentVolumeClaim, error) {
	obj.verify()
	if obj.err != nil {
		return nil, obj.err
	}
	return obj.pvc, nil
}
