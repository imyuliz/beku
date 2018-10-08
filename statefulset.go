package beku

import (
	"errors"
	"strings"

	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// StatefulSet include kubernetes resource object StatefulSet and error
type StatefulSet struct {
	sts *v1.StatefulSet
	err error
}

// NewSts  create sts
func NewSts() *StatefulSet {
	return &StatefulSet{
		sts: &v1.StatefulSet{},
	}

}

// SetName set name
func (obj *StatefulSet) SetName(name string) *StatefulSet {
	obj.sts.SetName(name)
	return obj
}

// SetNameSpace set sts namespace ,default namespace is default
func (obj *StatefulSet) SetNameSpace(namespace string) *StatefulSet {
	obj.sts.SetNamespace(namespace)
	return obj
}

// SetLabels set sts Labels
func (obj *StatefulSet) SetLabels(labels map[string]string) *StatefulSet {
	obj.sts.SetLabels(labels)
	return obj
}

// SetReplicas set sts replicas
func (obj *StatefulSet) SetReplicas(replicas int32) *StatefulSet {
	obj.sts.Spec.Replicas = &replicas
	return obj
}

// SetSelector set sts labels selector
func (obj *StatefulSet) SetSelector(labels map[string]string) *StatefulSet {
	if len(labels) <= 0 {
		obj.err = errors.New("LabelSelector set  error,label not allow empty")
		return obj
	}
	if obj.sts.Spec.Selector == nil {
		obj.sts.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: labels,
		}
	} else {
		obj.sts.Spec.Selector.MatchLabels = labels
	}
	obj.sts.Spec.Template.SetLabels(labels)
	return obj
}

// GetPodLabel get pod labels
func (obj *StatefulSet) GetPodLabel() map[string]string {
	return obj.sts.Spec.Template.GetLabels()
}

// SetPodLabels set pod labels and set sts selector
func (obj *StatefulSet) SetPodLabels(labels map[string]string) *StatefulSet {
	obj.sts.Spec.Template.SetLabels(labels)
	obj.SetSelector(labels)
	return obj
}

// SetContainer set deployment container
// name:name is container name ,default ""
// image:image is image name ,must input image
// containerPort: image expose containerPort,must input containerPort
func (obj *StatefulSet) SetContainer(name, image string, containerPort int32) *StatefulSet {
	// This must be a valid port number, 0 < x < 65536.
	if containerPort <= 0 || containerPort >= 65536 {
		obj.err = errors.New("containerPort error when SetContainer: 0 < containerPort < 65536")
		return obj
	}
	if !verifyString(image) {
		obj.err = errors.New("image not allow empty,must input image")
		return obj
	}
	port := corev1.ContainerPort{
		ContainerPort: containerPort,
	}
	container := corev1.Container{
		Name:  name,
		Image: image,
		Ports: []corev1.ContainerPort{port},
	}
	if obj.sts.Spec.Template.Spec.Containers == nil {
		obj.sts.Spec.Template.Spec.Containers = []corev1.Container{container}
		return obj
	}
	containersLen := len(obj.sts.Spec.Template.Spec.Containers)
	for index := 0; index < containersLen; index++ {
		img := strings.TrimSpace(obj.sts.Spec.Template.Spec.Containers[index].Image)
		if img == "" || len(img) <= 0 {
			obj.sts.Spec.Template.Spec.Containers[index] = container
			return obj
		}
	}
	obj.sts.Spec.Template.Spec.Containers = append(obj.sts.Spec.Template.Spec.Containers, container)
	return obj
}

func (obj *StatefulSet) verify() {
	if !verifyString(obj.sts.GetName()) {
		obj.err = errors.New("sts name  not allow empty")
		return
	}
	if obj.sts.Spec.Selector == nil {
		obj.err = errors.New("sts labels selector not allow empty")
		return
	}
	if len(obj.sts.Spec.Template.Spec.Containers) < 1 {
		obj.err = errors.New("sts container not allow empty")
		return
	}
}

// Finish the final step,will return kubernetes resource object secret and error
func (obj *StatefulSet) Finish() (*v1.StatefulSet, error) {
	if obj.err != nil {
		return nil, obj.err
	}
	return obj.sts, nil
}
