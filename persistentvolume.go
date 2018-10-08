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
	v1  *v1.PersistentVolume
	err error
}

// NewPV create pv
func NewPV() *PersistentVolume {
	return &PersistentVolume{
		v1: &v1.PersistentVolume{},
	}
}

// JSONNew json create  pv
func (pv *PersistentVolume) JSONNew(jsonbyte []byte) *PersistentVolume {
	pv.err = json.Unmarshal(jsonbyte, pv.v1)
	return pv
}

// SetLabels set persistentVolume label
func (pv *PersistentVolume) SetLabels(labels map[string]string) *PersistentVolume {
	pv.v1.SetLabels(labels)
	return pv
}

// SetName set pv name
func (pv *PersistentVolume) SetName(name string) *PersistentVolume {
	pv.v1.SetName(name)
	return pv
}

// SetNamespace set pv namespace ,default namesapce is default
func (pv *PersistentVolume) SetNamespace(namespace string) *PersistentVolume {
	pv.v1.SetNamespace(namespace)
	return pv
}

// SetAnnotations set pv annotations
func (pv *PersistentVolume) SetAnnotations(annotations map[string]string) *PersistentVolume {
	pv.v1.SetAnnotations(annotations)
	return pv
}

// SetAccessMode set pv access mode,only one
func (pv *PersistentVolume) SetAccessMode(mode core.PersistentVolumeAccessMode) *PersistentVolume {
	pv.v1.Spec.AccessModes = []v1.PersistentVolumeAccessMode{mode.ToK8s()}
	return pv
}

// SetAccessModes set pv access mode ,many modes
func (pv *PersistentVolume) SetAccessModes(modes []core.PersistentVolumeAccessMode) *PersistentVolume {
	var pvModes []v1.PersistentVolumeAccessMode
	for _, m := range modes {
		pvModes = append(pvModes, m.ToK8s())
	}
	pv.v1.Spec.AccessModes = pvModes
	return pv
}

// SetNFS set pv volume source nfs
func (pv *PersistentVolume) SetNFS(nfs *core.NFSVolumeSource) *PersistentVolume {
	if !verifyString(nfs.Server) {
		pv.err = errors.New("nfs server not allow empty")
		return pv
	}
	if !verifyString(nfs.Path) {
		pv.err = errors.New("nfs path not allow empty")
		return pv
	}
	nfsv := new(v1.NFSVolumeSource)
	err := mapper.Mapper(nfs, nfsv)
	if err != nil {
		pv.err = fmt.Errorf("set nfs error:%v", err)
		return pv
	}
	pv.v1.Spec.PersistentVolumeSource.NFS = nfsv
	return pv
}

// SetCapacity set pv capacity
func (pv *PersistentVolume) SetCapacity(capMaps map[core.ResourceName]string) *PersistentVolume {
	data, err := core.ResourceMapsToK8s(capMaps)
	if err != nil {
		pv.err = fmt.Errorf("set capacity err:%v", err)
		return pv
	}
	pv.v1.Spec.Capacity = data
	return pv
}

// SetCephFS set pv  volume source ceph
func (pv *PersistentVolume) SetCephFS(cephFs *core.CephFSPersistentVolumeSource) *PersistentVolume {
	if len(cephFs.Monitors) < 1 {
		pv.err = errors.New("cephFs monitor not allow empty")
		return pv
	}
	ceph := &v1.CephFSPersistentVolumeSource{
		SecretRef: new(v1.SecretReference),
	}
	err := mapper.Mapper(cephFs, ceph)
	if err != nil {
		pv.err = fmt.Errorf("set CephFs error:%v", err)
	}
	pv.v1.Spec.PersistentVolumeSource.CephFS = ceph
	return pv
}

// SetRBD  set pv volume source is rbd
func (pv *PersistentVolume) SetRBD(rbd *core.RBDPersistentVolumeSource) *PersistentVolume {
	if len(rbd.CephMonitors) < 1 {
		pv.err = errors.New("rbd CephMonitor not allow empty")
		return pv
	}
	if !verifyString(rbd.RBDImage) {
		pv.err = errors.New("rbd RBDImage not allow empty")
		return pv
	}
	if !verifyString(rbd.FSType) {
		pv.err = errors.New("rbd FSType not allow empty,maybe you can input one of  'ext4', 'xfs', 'ntfs'")

	}
	rbds := new(v1.RBDPersistentVolumeSource)
	err := mapper.Mapper(rbd, rbds)
	if err != nil {
		pv.err = fmt.Errorf("SetRBD error:%v", err)
		return pv
	}
	pv.v1.Spec.PersistentVolumeSource.RBD = rbds
	return pv
}

// Verify Verify pv
func (pv *PersistentVolume) verify() {
	if !verifyString(pv.v1.GetName()) {
		pv.err = errors.New("pv name not allow empty")
		return
	}
	if pv.v1.Spec.AccessModes == nil || len(pv.v1.Spec.AccessModes) < 1 {
		pv.err = errors.New("pv accessModes not allow empty")
		return
	}
	if pv.v1.Spec.Capacity == nil || len(pv.v1.Spec.Capacity) < 1 {
		pv.err = errors.New("pv capacity not allow  empty")
		return
	}
	var pvs v1.PersistentVolumeSource
	if pv.v1.Spec.PersistentVolumeSource == pvs {
		pv.err = errors.New("pv persistentVolumessource not allow empty")
		return
	}
	pv.v1.Kind = "PersistentVolume"
	pv.v1.APIVersion = "v1"
}

// Finish  set pv finnaly step, will return pv and err
func (pv *PersistentVolume) Finish() (*v1.PersistentVolume, error) {
	pv.verify()
	if pv.err != nil {
		return nil, pv.err
	}
	return pv.v1, nil
}
