package beku

import (
	"errors"
	"reflect"

	"github.com/yulibaozi/beku/core"
	"k8s.io/api/core/v1"
)

/*
	pvc and pv union release,will create kuberentes resource object  pvc and pv.
*/

// UnionPV output pvc and pv
type UnionPV struct {
	pv  *PersistentVolume
	pvc *PersistentVolumeClaim
	err error
}

// NewUnionPV create UnionPV
func NewUnionPV() *UnionPV {
	return &UnionPV{
		pv:  NewPV(),
		pvc: NewPVC(),
	}
}

// SetName set pvc and pv name
func (un *UnionPV) SetName(name string) *UnionPV {
	pvname, pvcname := name+"-pv", name+"-pvc"
	un.pv.SetName(pvname)
	un.pvc.SetName(pvcname)
	return un
}

// SetNamespace set pvc namespace,pv not have namespece
func (un *UnionPV) SetNamespace(namespace string) *UnionPV {
	un.pvc.SetNameSpace(namespace)
	return un
}

// SetAccessMode set pvc and pv accessmode
func (un *UnionPV) SetAccessMode(mode core.PersistentVolumeAccessMode) *UnionPV {
	un.pvc.SetAccessMode(mode)
	un.pv.SetAccessMode(mode)
	return un
}

// SetAccessModes set pv and pv  access modes
func (un *UnionPV) SetAccessModes(modes []core.PersistentVolumeAccessMode) *UnionPV {
	un.pvc.SetAccessModes(modes)
	un.pv.SetAccessModes(modes)
	return un
}

// SetCapacity set pv capacity and set pvc resource request
func (un *UnionPV) SetCapacity(capMaps map[core.ResourceName]string) *UnionPV {
	un.pv.SetCapacity(capMaps)
	un.pvc.SetResourceRequests(capMaps)
	return un
}

// SetVolumeMode set pvc volume mode Filesystem or Block
func (un *UnionPV) SetVolumeMode(volumeMode core.PersistentVolumeMode) *UnionPV {
	un.pvc.SetVolumeMode(volumeMode)
	return un
}

// SetLabels set pv labels ,set pvc labels and set pvc selector
func (un *UnionPV) SetLabels(labels map[string]string) *UnionPV {
	pvlabels := make(map[string]string, 0)
	pvclabels := make(map[string]string, 0)
	for k, v := range labels {
		pvclabels[k] = v + "-pvc"
		pvlabels[k] = v + "-pv"
	}
	un.pv.SetLabels(pvlabels)
	un.pvc.SetLabels(pvclabels)
	un.pvc.SetSelector(pvlabels)
	return un
}

// SetNFS set pv volume source is nfs
func (un *UnionPV) SetNFS(nfs *core.NFSVolumeSource) *UnionPV {
	un.pv.SetNFS(nfs)
	return un
}

// SetRBD set pv volume source is rbd
func (un *UnionPV) SetRBD(rbd *core.RBDPersistentVolumeSource) *UnionPV {
	un.pv.SetRBD(rbd)
	return un
}

func (un *UnionPV) verify() {
	if un.err != nil {
		return
	}
	pvname, pvcname := un.pv.GetName(), un.pvc.GetName()
	if !verifyString(pvname) || !verifyString(pvcname) {
		un.err = errors.New("pvc or pv name is empty not allow")
		return
	}
	//check labels and selector
	pvlabels, pvclabels, pvcselector := un.pv.GetLabels(), un.pvc.GetLabels(), un.pvc.GetSelector()
	if !reflect.DeepEqual(pvcselector, pvlabels) {
		un.err = errors.New("UnionPV, it is not allow to pvc selector and pv labels not equal")
		return
	}
	if !verifyString(un.pvc.GetNamespace()) {
		un.SetNamespace("default")
	}

	if pvlabels == nil {
		pvlabels = map[string]string{"name": pvname}
		un.pv.SetLabels(pvlabels)
		un.pvc.SetSelector(pvlabels)
	}
	if pvclabels == nil {
		pvclabels = map[string]string{"name": pvcname}
		un.pvc.SetLabels(pvclabels)
	}

}

// Finish the finalstep, will return kubernetes resource object pvc,pv and error
func (un *UnionPV) Finish() (pv *v1.PersistentVolume, pvc *v1.PersistentVolumeClaim, err error) {
	un.verify()
	if un.err != nil {
		err = un.err
		return
	}
	pv, err = un.pv.Finish()
	if err != nil {
		return
	}
	pvc, err = un.pvc.Finish()
	return
}
