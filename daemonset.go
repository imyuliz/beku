package beku

import (
	"errors"
	"strings"

	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DaemonSet include Kubernets resource object DaemonSet and error
type DaemonSet struct {
	ds  *v1.DaemonSet
	err error
}

// NewDS create DaemonSet(ds) and chain function call begin with this function.
func (obj *DaemonSet) NewDS() *DaemonSet { return &DaemonSet{ds: &v1.DaemonSet{}} }

// Finish Chain function call end with this function
// return real DaemonSet(really DaemonSet is kubernetes resource object DaemonSet and error
// In the function, it will check necessary parameters、input the default field。
func (obj *DaemonSet) Finish() (*v1.DaemonSet, error) {
	if obj.err != nil {
		return nil, obj.err
	}
	return obj.ds, nil
}

// SetName set DaemonSet(ds) name
func (obj *DaemonSet) SetName(name string) *DaemonSet {
	obj.ds.SetName(name)
	return obj
}

// SetNamespace set DaemonSet(ds) namespace, default namespace value is 'default'
func (obj *DaemonSet) SetNamespace(namespace string) *DaemonSet {
	obj.ds.SetNamespace(namespace)
	return obj
}

// SetNamespaceAndName set DaemonSet namespace and DaemonSet name and  Pod Namespace
func (obj *DaemonSet) SetNamespaceAndName(namespace, name string) *DaemonSet {
	obj.ds.SetName(name)
	obj.ds.SetNamespace(namespace)
	obj.ds.Spec.Template.SetNamespace(namespace)
	return obj
}

// SetLabels set DaemonSet(ds) Labels,set Pod Labels and set DaemonSet selector.
func (obj *DaemonSet) SetLabels(labels map[string]string) *DaemonSet {
	obj.ds.SetLabels(labels)
	obj.SetPodLabels(labels)
	return obj
}

// SetSelector set DaemonSet(ds) Selector and Set Pod Label
// The Pod that matches the seletor will be selected, DaemonSet will controller the Pod.
func (obj *DaemonSet) SetSelector(selector map[string]string) *DaemonSet {
	obj.ds.Spec.Template.SetLabels(selector)
	if obj.ds.Spec.Selector == nil {
		obj.ds.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: selector,
		}
		return obj
	}
	obj.ds.Spec.Selector.MatchLabels = selector
	return obj
}

// SetPodLabels set Pod Label and set DaemonSet Selector
func (obj *DaemonSet) SetPodLabels(labels map[string]string) *DaemonSet {
	obj.ds.Spec.Template.SetLabels(labels)
	obj.SetPodLabels(labels)
	return obj
}

// SetContainer set DaemonSet container
// name Not required when only one Conatiner,you can input "".
// when many container this Field is necessary and cann't repeat
// image is necessary, image very important
// containerPort container port,this is necessary
func (obj *DaemonSet) SetContainer(name, image string, containerPort int32) *DaemonSet {
	if containerPort <= 0 || containerPort >= 65536 {
		obj.err = errors.New("SetContainer err, container Port range: 0 < containerPort < 65536")
		return obj
	}
	if !verifyString(image) {
		obj.err = errors.New("SetContainer err, image is not allowed to be empty")
		return obj
	}
	port := corev1.ContainerPort{ContainerPort: containerPort}
	container := corev1.Container{
		Name:  name,
		Image: image,
		Ports: []corev1.ContainerPort{port},
	}
	containersLen := len(obj.ds.Spec.Template.Spec.Containers)
	if containersLen < 1 {
		obj.ds.Spec.Template.Spec.Containers = []corev1.Container{container}
		return obj
	}
	for index := 0; index < containersLen; index++ {
		// if iamge not exist
		img := strings.TrimSpace(obj.ds.Spec.Template.Spec.Containers[index].Image)
		if img == "" || len(img) <= 0 {
			obj.ds.Spec.Template.Spec.Containers[index].Name = name
			obj.ds.Spec.Template.Spec.Containers[index].Image = image
			obj.ds.Spec.Template.Spec.Containers[index].Ports = []corev1.ContainerPort{port}
			return obj
		}
	}
	obj.ds.Spec.Template.Spec.Containers = append(obj.ds.Spec.Template.Spec.Containers, container)
	return obj
}

// SetEnvs set Pod Environmental variable
func (obj *DaemonSet) SetEnvs(envMap map[string]string) *DaemonSet {
	envs, err := mapToEnvs(envMap)
	if err != nil {
		obj.err = err
		return obj
	}
	containerLen := len(obj.ds.Spec.Template.Spec.Containers)
	if containerLen < 1 {
		obj.ds.Spec.Template.Spec.Containers = []corev1.Container{corev1.Container{Env: envs}}
		return obj
	}
	for index := 0; index < containerLen; index++ {
		if obj.ds.Spec.Template.Spec.Containers[index].Env == nil {
			obj.ds.Spec.Template.Spec.Containers[index].Env = envs
		}
	}
	return obj
}

// GetPodLabel get pod labels
func (obj *DaemonSet) GetPodLabel() map[string]string {
	return obj.ds.Spec.Template.GetLabels()
}

// verify check service necessary value, input the default field and input related data.
func (obj *DaemonSet) verify() {
	if obj.err != nil {
		return
	}
	if !verifyString(obj.ds.Name) {
		obj.err = errors.New("DaemonSet.Name is not allowed to be empty")
		return
	}
	if obj.ds.Spec.Template.Spec.Containers == nil || len(obj.ds.Spec.Template.Spec.Containers) < 1 {
		obj.err = errors.New("DaemonSet.Spec.Template.Spec.Containers is not allowed to be empty")
		return
	}
	if len(obj.GetPodLabel()) < 1 {
		obj.err = errors.New("Pod Labels is not allowed to be empty,you can call SetPodLabels input")
	}
	obj.ds.Kind = "DaemonSet"
	obj.ds.APIVersion = "app/v1"
}
