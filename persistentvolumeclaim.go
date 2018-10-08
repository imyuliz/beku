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
	v1  *v1.PersistentVolumeClaim
	err error
}

// NewPVC create pvc
func NewPVC() *PersistentVolumeClaim {
	return &PersistentVolumeClaim{
		v1: &v1.PersistentVolumeClaim{},
	}
}

// SetName set pvc name
func (pvc *PersistentVolumeClaim) SetName(name string) *PersistentVolumeClaim {
	pvc.v1.SetName(name)
	return pvc
}

// SetNameSpace set  pvc namespace,default namespace is default
func (pvc *PersistentVolumeClaim) SetNameSpace(namespace string) *PersistentVolumeClaim {
	pvc.v1.SetNamespace(namespace)
	return pvc
}

// SetLabels set pvc label
func (pvc *PersistentVolumeClaim) SetLabels(labels map[string]string) *PersistentVolumeClaim {
	pvc.v1.SetLabels(labels)
	return pvc
}

// SetAnnotations set annotation
func (pvc *PersistentVolumeClaim) SetAnnotations(annotations map[string]string) *PersistentVolumeClaim {
	pvc.v1.SetAnnotations(annotations)
	return pvc
}

// SetPVCAccessMode set pvc accessMode
func (pvc *PersistentVolumeClaim) SetPVCAccessMode(mode core.PersistentVolumeAccessMode) *PersistentVolumeClaim {
	pvc.v1.Spec.AccessModes = []v1.PersistentVolumeAccessMode{mode.ToK8s()}
	return pvc
}

// SetPVCAccessModes set pvc accessModes
func (pvc *PersistentVolumeClaim) SetPVCAccessModes(modes []core.PersistentVolumeAccessMode) *PersistentVolumeClaim {
	var pvcModes []v1.PersistentVolumeAccessMode
	for _, m := range modes {
		pvcModes = append(pvcModes, m.ToK8s())
	}
	pvc.v1.Spec.AccessModes = pvcModes
	return pvc
}

// SetVolumeMode pvc vloume mode,have Block and Filesystem  mode
func (pvc *PersistentVolumeClaim) SetVolumeMode(volumeMode core.PersistentVolumeMode) *PersistentVolumeClaim {
	m := volumeMode.ToK8s()
	if m == nil {
		pvc.err = fmt.Errorf("set volumeMode err: this volumeMode not allow %v", volumeMode)
	}
	pvc.v1.Spec.VolumeMode = m
	return pvc
}

// SetResourceLimit set pvc resource limit
func (pvc *PersistentVolumeClaim) SetResourceLimit(limits map[core.ResourceName]string) *PersistentVolumeClaim {
	data, err := core.ResourceMapsToK8s(limits)

	if err != nil {
		pvc.err = fmt.Errorf("limit set err:%v", err)
		return pvc
	}
	pvc.v1.Spec.Resources.Limits = data
	return pvc
}

// SetResourceRequests set pvc reource requests
func (pvc *PersistentVolumeClaim) SetResourceRequests(requests map[core.ResourceName]string) *PersistentVolumeClaim {
	data, err := core.ResourceMapsToK8s(requests)
	if err != nil {
		pvc.err = fmt.Errorf("request set err:%v", err)
		return pvc
	}
	pvc.v1.Spec.Resources.Requests = data
	return pvc
}

// SetStorageClassName set storageclasss name
func (pvc *PersistentVolumeClaim) SetStorageClassName(classname string) *PersistentVolumeClaim {
	if classname == "" || len(classname) <= 0 {
		pvc.err = errors.New("set StorageClassName is empty,set failed")
		return pvc
	}
	pvc.v1.Spec.StorageClassName = &classname
	return pvc
}

// SetSelector set pvc selector
func (pvc *PersistentVolumeClaim) SetSelector(labels map[string]string) *PersistentVolumeClaim {
	if len(labels) < 1 {
		pvc.err = errors.New("set LabelSelector error, labels is empty")
		return pvc
	}
	if pvc.v1.Spec.Selector == nil {
		selector := &metav1.LabelSelector{
			MatchLabels: labels,
		}
		pvc.v1.Spec.Selector = selector
		return pvc
	}
	pvc.v1.Spec.Selector.MatchLabels = labels
	return pvc
}

// SetMatchExpressions set pvc label selector,have key,operator and values
func (pvc *PersistentVolumeClaim) SetMatchExpressions(ents []core.LabelSelectorRequirement) *PersistentVolumeClaim {
	requirements := make([]metav1.LabelSelectorRequirement, 0)
	err := mapper.AutoMapper(ents, requirements)
	if err != nil {
		pvc.err = fmt.Errorf("SetMatchExpressions error:%v", err)
		return pvc
	}
	if pvc.v1.Spec.Selector == nil {
		selector := &metav1.LabelSelector{
			MatchExpressions: requirements,
		}
		pvc.v1.Spec.Selector = selector
		return pvc
	}
	pvc.v1.Spec.Selector.MatchExpressions = requirements
	return pvc
}

// verify  pvc
func (pvc *PersistentVolumeClaim) verify() {
	if !verifyString(pvc.v1.GetName()) {
		pvc.err = errors.New("pvc name not allow empty")
		return
	}
	if pvc.v1.Spec.AccessModes == nil || len(pvc.v1.Spec.AccessModes) < 1 {
		pvc.err = errors.New("pvc accessModes not allow empty")
		return
	}
	if pvc.v1.Spec.VolumeMode == nil {
		pvc.err = errors.New("pvc volumeMode not allow nil")
		return
	}
	if pvc.v1.Spec.Resources.Limits == nil && pvc.v1.Spec.Resources.Requests == nil {
		pvc.err = errors.New("both limits and requests is nil  not allow")
		return
	}
	pvc.v1.Kind = "PersistentVolumeClaim"
	pvc.v1.APIVersion = "v1"
}

// Finish  the final step,will return kubernetes resource object pv and error
func (pvc *PersistentVolumeClaim) Finish() (*v1.PersistentVolumeClaim, error) {
	pvc.verify()
	if pvc.err != nil {
		return nil, pvc.err
	}
	return pvc.v1, nil
}
