package beku

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/yulibaozi/beku/core"
	"github.com/yulibaozi/mapper"
	"k8s.io/api/core/v1"
)

// PersistentVolume pv
type PersistentVolume struct {
	pv  *v1.PersistentVolume
	err error
}

// Newobj create pv
func NewPV() *PersistentVolume {
	return &PersistentVolume{
		pv: &v1.PersistentVolume{},
	}
}

// JSONNew json create  pv
func (obj *PersistentVolume) JSONNew(jsonbyte []byte) *PersistentVolume {
	obj.err = json.Unmarshal(jsonbyte, obj.pv)
	return obj
}

// SetLabels set persistentVolume label
func (obj *PersistentVolume) SetLabels(labels map[string]string) *PersistentVolume {
	obj.pv.SetLabels(labels)
	return obj
}

// SetName set pv name
func (obj *PersistentVolume) SetName(name string) *PersistentVolume {
	obj.pv.SetName(name)
	return obj
}

// // SetNamespace set pv namespace ,default namesapce is default
// func (obj *PersistentVolume) SetNamespace(namespace string) *PersistentVolume {
// 	obj.pv.SetNamespace(namespace)
// 	return obj
// }

// SetAnnotations set pv annotations
func (obj *PersistentVolume) SetAnnotations(annotations map[string]string) *PersistentVolume {
	obj.pv.SetAnnotations(annotations)
	return obj
}

// SetAccessMode set pv access mode,only one
func (obj *PersistentVolume) SetAccessMode(mode core.PersistentVolumeAccessMode) *PersistentVolume {
	obj.pv.Spec.AccessModes = []v1.PersistentVolumeAccessMode{mode.ToK8s()}
	return obj
}

// SetAccessModes set pv access mode ,many modes
func (obj *PersistentVolume) SetAccessModes(modes []core.PersistentVolumeAccessMode) *PersistentVolume {
	var objModes []v1.PersistentVolumeAccessMode
	for _, m := range modes {
		objModes = append(objModes, m.ToK8s())
	}
	obj.pv.Spec.AccessModes = objModes
	return obj
}

// SetNFS set pv volume source nfs
func (obj *PersistentVolume) SetNFS(nfs *core.NFSVolumeSource) *PersistentVolume {
	if !verifyString(nfs.Server) {
		obj.err = errors.New("nfs server not allow empty")
		return obj
	}
	if !verifyString(nfs.Path) {
		obj.err = errors.New("nfs path not allow empty")
		return obj
	}
	nfsv := new(v1.NFSVolumeSource)
	err := mapper.Mapper(nfs, nfsv)
	if err != nil {
		obj.err = fmt.Errorf("set nfs error:%v", err)
		return obj
	}
	obj.pv.Spec.PersistentVolumeSource.NFS = nfsv
	return obj
}

// SetCapacity set pv capacity
func (obj *PersistentVolume) SetCapacity(capMaps map[core.ResourceName]string) *PersistentVolume {
	data, err := core.ResourceMapsToK8s(capMaps)
	if err != nil {
		obj.err = fmt.Errorf("set capacity err:%v", err)
		return obj
	}
	obj.pv.Spec.Capacity = data
	return obj
}

// SetCephFS set pv  volume source ceph
func (obj *PersistentVolume) SetCephFS(cephFs *core.CephFSPersistentVolumeSource) *PersistentVolume {
	if len(cephFs.Monitors) < 1 {
		obj.err = errors.New("cephFs monitor not allow empty")
		return obj
	}
	ceph := &v1.CephFSPersistentVolumeSource{
		SecretRef: new(v1.SecretReference),
	}
	err := mapper.Mapper(cephFs, ceph)
	if err != nil {
		obj.err = fmt.Errorf("set CephFs error:%v", err)
	}
	obj.pv.Spec.PersistentVolumeSource.CephFS = ceph
	return obj
}

// SetRBD  set pv volume source is rbd
func (obj *PersistentVolume) SetRBD(rbd *core.RBDPersistentVolumeSource) *PersistentVolume {
	if len(rbd.CephMonitors) < 1 {
		obj.err = errors.New("rbd CephMonitor not allow empty")
		return obj
	}
	if !verifyString(rbd.RBDImage) {
		obj.err = errors.New("rbd RBDImage not allow empty")
		return obj
	}
	if !verifyString(rbd.FSType) {
		obj.err = errors.New("rbd FSType not allow empty,maybe you can input one of  'ext4', 'xfs', 'ntfs'")

	}
	rbds := new(v1.RBDPersistentVolumeSource)
	err := mapper.Mapper(rbd, rbds)
	if err != nil {
		obj.err = fmt.Errorf("SetRBD error:%v", err)
		return obj
	}
	obj.pv.Spec.PersistentVolumeSource.RBD = rbds
	return obj
}

// Verify Verify pv
func (obj *PersistentVolume) verify() {
	if !verifyString(obj.pv.GetName()) {
		obj.err = errors.New("obj name not allow empty")
		return
	}
	if obj.pv.Spec.AccessModes == nil || len(obj.pv.Spec.AccessModes) < 1 {
		obj.err = errors.New("obj accessModes not allow empty")
		return
	}
	if obj.pv.Spec.Capacity == nil || len(obj.pv.Spec.Capacity) < 1 {
		obj.err = errors.New("obj capacity not allow  empty")
		return
	}
	var objs v1.PersistentVolumeSource
	if obj.pv.Spec.PersistentVolumeSource == objs {
		obj.err = errors.New("obj persistentVolumessource not allow empty")
		return
	}
	obj.pv.Kind = "PersistentVolume"
	obj.pv.APIVersion = "v1"
}

// Finish  set obj finnaly step, will return pv and err
func (obj *PersistentVolume) Finish() (*v1.PersistentVolume, error) {
	obj.verify()
	if obj.err != nil {
		return nil, obj.err
	}
	return obj.pv, nil
}
