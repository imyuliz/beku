package beku

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/yulibaozi/mapper"
	"k8s.io/api/core/v1"
)

// PersistentVolume include Kubernetes resource object PersistentVolume(pv) and error.
type PersistentVolume struct {
	pv  *v1.PersistentVolume
	err error
}

// NewPV create PersistentVolume and chain function call begin with this function.
func NewPV() *PersistentVolume {
	return &PersistentVolume{
		pv: &v1.PersistentVolume{},
	}
}

// Finish chain function call end with this function
// return Kubernetes resource object PersistentVolume(pv) and error.
// In the function, it will check necessary parameters、input the default field。
func (obj *PersistentVolume) Finish() (*v1.PersistentVolume, error) {
	obj.verify()
	return obj.pv, obj.err
}

// JSONNew use json data create PersistentVolume(pv)
func (obj *PersistentVolume) JSONNew(jsonbyte []byte) *PersistentVolume {
	obj.error(json.Unmarshal(jsonbyte, obj.pv))
	return obj
}

// YAMLNew use yaml data create PersistentVolume(pv)
func (obj *PersistentVolume) YAMLNew(yamlbyts []byte) *PersistentVolume {
	obj.error(yaml.Unmarshal(yamlbyts, obj.pv))
	return obj
}

// SetLabels set PersistentVolume(pv) label
func (obj *PersistentVolume) SetLabels(labels map[string]string) *PersistentVolume {
	obj.pv.SetLabels(labels)
	return obj
}

// GetLabels get PersistentVolume(pv) labels
func (obj *PersistentVolume) GetLabels() map[string]string {
	return obj.pv.GetLabels()
}

// SetName set PersistentVolume(pv) name
func (obj *PersistentVolume) SetName(name string) *PersistentVolume {
	obj.pv.SetName(name)
	return obj
}

// GetName get PersistentVolume(pv) name
func (obj *PersistentVolume) GetName() string {
	return obj.pv.GetName()
}

// SetAnnotations set  PersistentVolume(pv) annotations
func (obj *PersistentVolume) SetAnnotations(annotations map[string]string) *PersistentVolume {
	obj.pv.SetAnnotations(annotations)
	return obj
}

// SetAccessMode set PersistentVolume(pv) access mode, only one
func (obj *PersistentVolume) SetAccessMode(mode PersistentVolumeAccessMode) *PersistentVolume {
	obj.pv.Spec.AccessModes = []v1.PersistentVolumeAccessMode{mode.ToK8s()}
	return obj
}

// SetAccessModes set PersistentVolume(pv) access mode, many modes
func (obj *PersistentVolume) SetAccessModes(modes []PersistentVolumeAccessMode) *PersistentVolume {
	var objModes []v1.PersistentVolumeAccessMode
	for _, m := range modes {
		objModes = append(objModes, m.ToK8s())
	}
	obj.pv.Spec.AccessModes = objModes
	return obj
}

// SetNFS set PersistentVolume(pv) volume source is nfs
func (obj *PersistentVolume) SetNFS(nfs *NFSVolumeSource) *PersistentVolume {
	if !verifyString(nfs.Server) {
		obj.error(errors.New("SetNFS err, nfs server is not allowed to be empty"))
		return obj
	}
	if !verifyString(nfs.Path) {
		obj.error(errors.New("SetNFS err, nfs path is not allowed to be empty"))
		return obj
	}
	nfsv := new(v1.NFSVolumeSource)
	err := mapper.Mapper(nfs, nfsv)
	if err != nil {
		obj.error(fmt.Errorf("SetNFS err:%v", err))
		return obj
	}
	obj.pv.Spec.PersistentVolumeSource.NFS = nfsv
	return obj
}

// SetCapacity set PersistentVolume(pv) capacity
func (obj *PersistentVolume) SetCapacity(capMaps map[ResourceName]string) *PersistentVolume {
	data, err := ResourceMapsToK8s(capMaps)
	if err != nil {
		obj.error(fmt.Errorf("SetCapacity err:%v", err))
		return obj
	}
	obj.pv.Spec.Capacity = data
	return obj
}

// SetCephFS set PersistentVolume(pv) volume source is ceph
func (obj *PersistentVolume) SetCephFS(cephFs *CephFSPersistentVolumeSource) *PersistentVolume {
	if len(cephFs.Monitors) < 1 {
		obj.error(errors.New("SetCephFS err,cephFS monitor is not allowed to be empty"))
		return obj
	}
	ceph := &v1.CephFSPersistentVolumeSource{
		Monitors:   cephFs.Monitors,
		Path:       cephFs.Path,
		User:       cephFs.User,
		SecretFile: cephFs.SecretFile,
		ReadOnly:   cephFs.ReadOnly,
	}
	if cephFs.SecretRef != nil {
		ceph.SecretRef = &v1.SecretReference{
			Name:      ceph.SecretRef.Name,
			Namespace: ceph.SecretRef.Namespace,
		}
	}
	obj.pv.Spec.PersistentVolumeSource.CephFS = ceph
	return obj
}

// SetRBD  set PersistentVolume(pv) volume source is RBD
func (obj *PersistentVolume) SetRBD(rbd *RBDPersistentVolumeSource) *PersistentVolume {
	if len(rbd.CephMonitors) < 1 {
		obj.error(errors.New("SetRBD err, CephMonitor is not allowed to be empty"))
		return obj
	}
	if !verifyString(rbd.RBDImage) {
		obj.error(errors.New("SetRBD err, RBDImage is not allowed to be empty"))
		return obj
	}
	if !verifyString(rbd.FSType) {
		obj.error(errors.New("SetRBD err, RBD.FSType is not allowed to be empty,maybe you can input one of  'ext4', 'xfs', 'ntfs'"))
		return obj
	}
	rbds := &v1.RBDPersistentVolumeSource{
		CephMonitors: rbd.CephMonitors,
		RBDImage:     rbd.RBDImage,
		FSType:       rbd.FSType,
		RBDPool:      rbd.RBDPool,
		RadosUser:    rbd.RadosUser,
		Keyring:      rbd.Keyring,
		ReadOnly:     rbd.ReadOnly,
	}
	if rbd.SecretRef != nil {
		rbds.SecretRef = &v1.SecretReference{
			Name:      rbd.SecretRef.Name,
			Namespace: rbd.SecretRef.Namespace,
		}
	}
	obj.pv.Spec.PersistentVolumeSource.RBD = rbds
	return obj
}

func (obj *PersistentVolume) error(err error) {
	if obj.err != nil {
		return
	}
	obj.err = err
}

// verify check service necessary value, input the default field and input related data.
func (obj *PersistentVolume) verify() {
	if obj.err != nil {
		return
	}
	if !verifyString(obj.pv.GetName()) {
		obj.err = errors.New("PersistentVolume.Name is not allowed to be empty")
		return
	}
	if obj.pv.Spec.AccessModes == nil || len(obj.pv.Spec.AccessModes) < 1 {
		obj.err = errors.New("PersistentVolume.Spec.AccessModes is not allowed to be empty")
		return
	}
	if obj.pv.Spec.Capacity == nil || len(obj.pv.Spec.Capacity) < 1 {
		obj.err = errors.New("PersistentVolume.Spec.Capacity is not allowed to be empty")
		return
	}
	var objs v1.PersistentVolumeSource
	if obj.pv.Spec.PersistentVolumeSource == objs {
		obj.err = errors.New("PersistentVolume.Spec.PersistentVolumeSource is not allowed to be empty")
		return
	}
	obj.pv.Kind = "PersistentVolume"
	obj.pv.APIVersion = "v1"
}
