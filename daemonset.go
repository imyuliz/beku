package beku

import (
	"errors"
	"strings"

	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DaemonSet include k8s resource object DaemonSet and error
type DaemonSet struct {
	v1  *v1.DaemonSet
	err error
}

// NewDS  new  DaemonSet
func (ds *DaemonSet) NewDS() *DaemonSet {
	return &DaemonSet{
		v1: &v1.DaemonSet{},
	}
}

// SetName set daemonSet name
func (ds *DaemonSet) SetName(name string) *DaemonSet {
	ds.v1.SetName(name)
	return ds
}

// SetNamespace set daemonSet namespace, default namespace is default
func (ds *DaemonSet) SetNamespace(namespace string) *DaemonSet {
	ds.v1.SetNamespace(namespace)
	return ds
}

// SetNamespaceAndName set namespace and name
func (ds *DaemonSet) SetNamespaceAndName(namespace, name string) *DaemonSet {
	ds.v1.SetName(name)
	ds.v1.SetNamespace(namespace)
	return ds
}

// SetLabels set daemonSet labels
func (ds *DaemonSet) SetLabels(labels map[string]string) *DaemonSet {
	ds.v1.SetLabels(labels)
	return ds
}

// SetSelector set Selector ,will check match label pod
func (ds *DaemonSet) SetSelector(selector map[string]string) *DaemonSet {
	if ds.v1.Spec.Selector == nil {
		ds.v1.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: selector,
		}
	}
	return ds
}

// SetPodLabels set  pod Label and will auto set DaemonSet SetSelector
func (ds *DaemonSet) SetPodLabels(labels map[string]string) *DaemonSet {
	ds.v1.Spec.Template.SetLabels(labels)
	ds.SetSelector(labels)
	return ds
}

// SetContainer set daemonSet container
// name Not required when only one Conatiner,you can input "",when many container this Field is necessary and cann't repeat
// image is necessary, image very important
// containerPort container port,this is necessary
func (ds *DaemonSet) SetContainer(name, image string, containerPort int32) *DaemonSet {
	if containerPort <= 0 || containerPort >= 65536 {
		ds.err = errors.New("containerPort error when SetContainer: 0 < containerPort < 65536")
		return ds
	}
	if !verifyString(image) {
		ds.err = errors.New("image not allow empty,must input image")
		return ds
	}
	port := corev1.ContainerPort{
		ContainerPort: containerPort,
	}
	container := corev1.Container{
		Name:  name,
		Image: image,
		Ports: []corev1.ContainerPort{port},
	}
	if ds.v1.Spec.Template.Spec.Containers == nil {
		ds.v1.Spec.Template.Spec.Containers = []corev1.Container{container}
		return ds
	}
	containersLen := len(ds.v1.Spec.Template.Spec.Containers)
	for index := 0; index < containersLen; index++ {
		img := strings.TrimSpace(ds.v1.Spec.Template.Spec.Containers[index].Image)
		if img == "" || len(img) <= 0 {
			ds.v1.Spec.Template.Spec.Containers[index] = container
			return ds
		}
	}
	ds.v1.Spec.Template.Spec.Containers = append(ds.v1.Spec.Template.Spec.Containers, container)
	return ds
}

// SetEnvs set pod Environmental variable
func (ds *DaemonSet) SetEnvs(envMap map[string]string) *DaemonSet {

	if len(envMap) <= 0 {
		ds.err = errors.New("set env error: envMap is empty")
		return ds
	}
	var envs []corev1.EnvVar
	for k, v := range envMap {
		k, v = strings.TrimSpace(k), strings.TrimSpace(v)
		if k == "" || v == "" {
			ds.err = errors.New("set env error: name or value not allow")
			return ds
		}
		envs = append(envs, corev1.EnvVar{Name: k, Value: v})
	}
	if len(envs) <= 0 {
		ds.err = errors.New("set env error, envs is empty")
		return ds
	}
	containerLen := len(ds.v1.Spec.Template.Spec.Containers)
	for index := 0; index < containerLen; index++ {
		if ds.v1.Spec.Template.Spec.Containers[index].Env == nil {
			ds.v1.Spec.Template.Spec.Containers[index].Env = envs
		}
	}
	return ds
}

// GetPodLabel get pod labels
func (ds *DaemonSet) GetPodLabel() map[string]string {
	return ds.v1.Spec.Template.GetLabels()
}

func (ds *DaemonSet) verify() {
	if !verifyString(ds.v1.Name) {
		ds.err = errors.New("daemonSet name not allow empty")
		return
	}
	if len(ds.v1.GetLabels()) < 1 {
		ds.err = errors.New("daemonSet labels not allow empty")
		return
	}
	if ds.v1.Spec.Template.Spec.Containers == nil || len(ds.v1.Spec.Template.Spec.Containers) < 1 {
		ds.err = errors.New("DaemonSet.Spec.Template.Spec.Containers not allow nil")
		return
	}
	if ds.v1.Spec.Selector == nil {
		ds.SetSelector(ds.GetPodLabel())
	}
	ds.v1.Kind = "DaemonSet"
	ds.v1.APIVersion = "app/v1"
}

// Finish the final step, will return kubernetes resource object  DaemonSet and error
func (ds *DaemonSet) Finish() (*v1.DaemonSet, error) {
	if ds.err != nil {
		return nil, ds.err
	}
	return ds.v1, nil
}
