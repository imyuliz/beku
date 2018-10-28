package beku

import (
	"errors"
	"reflect"

	"k8s.io/api/core/v1"
)

// UnionPV output pvc and pv
type UnionPV struct {
	pv  *PersistentVolume
	pvc *PersistentVolumeClaim
	err error
}

// NewUnionPV create PersistentVolume,PersistentVolumeClaim and error
// and chain function call begin with this function.
func NewUnionPV() *UnionPV { return &UnionPV{pv: NewPV(), pvc: NewPVC()} }

// Finish Chain function call end with this function
// return Kubernetes resource object(PersistentVolume,PersistentVolumeClaim) and error
// In the function, it will check necessary parametersainput the default field
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

// SetName set PersistentVolume and PersistentVolumeClaim name
func (un *UnionPV) SetName(name string) *UnionPV {
	un.pv.SetName(name)
	un.pvc.SetName(name)
	return un
}

// SetNamespace set PersistentVolumeClaim namespace,
// and PersistentVolume can't set namespece because of no such attribute
func (un *UnionPV) SetNamespace(namespace string) *UnionPV {
	un.pvc.SetNamespace(namespace)
	return un
}

// SetAccessMode set PersistentVolume and PersistentVolumeClaim accessmode
func (un *UnionPV) SetAccessMode(mode PersistentVolumeAccessMode) *UnionPV {
	un.pvc.SetAccessMode(mode)
	un.pv.SetAccessMode(mode)
	return un
}

// SetAccessModes PersistentVolume and PersistentVolumeClaim access modes
func (un *UnionPV) SetAccessModes(modes []PersistentVolumeAccessMode) *UnionPV {
	un.pvc.SetAccessModes(modes)
	un.pv.SetAccessModes(modes)
	return un
}

// SetCapacity set PersistentVolume capacity and set PersistentVolumeClaim resource request
func (un *UnionPV) SetCapacity(capMaps map[ResourceName]string) *UnionPV {
	un.pv.SetCapacity(capMaps)
	un.pvc.SetResourceRequests(capMaps)
	return un
}

// SetVolumeMode set pvc volume mode Filesystem or Block
func (un *UnionPV) SetVolumeMode(volumeMode PersistentVolumeMode) *UnionPV {
	un.pvc.SetVolumeMode(volumeMode)
	return un
}

// SetLabels set PersistentVolume labels ,set PersistentVolumeClaim labels and set PersistentVolumeClaim selector
func (un *UnionPV) SetLabels(labels map[string]string) *UnionPV {
	un.pv.SetLabels(labels)
	un.pvc.SetLabels(labels)
	un.pvc.SetSelector(labels)
	return un
}

// SetNFS set PersistentVolume volume source is NFS
func (un *UnionPV) SetNFS(nfs *NFSVolumeSource) *UnionPV {
	un.pv.SetNFS(nfs)
	return un
}

// SetRBD set PersistentVolume volume source is RBD
func (un *UnionPV) SetRBD(rbd *RBDPersistentVolumeSource) *UnionPV {
	un.pv.SetRBD(rbd)
	return un
}

// verify check UnionPV necessary value, input the default field and input related data.
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
