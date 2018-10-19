package beku

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/ghodss/yaml"
	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// StatefulSet include kubernetes resource object StatefulSet(sts) and error
type StatefulSet struct {
	sts *v1.StatefulSet
	err error
}

// NewSts  create StatefulSet(sts) and chain funtion call begin with this funtion.
func NewSts() *StatefulSet { return &StatefulSet{sts: &v1.StatefulSet{}} }

// Finish Chain function call end with this function
// return Kubernetes resource object StatefulSet and error.
// In the function, it will check necessary parameters、input the default field。
func (obj *StatefulSet) Finish() (*v1.StatefulSet, error) {
	if obj.err != nil {
		return nil, obj.err
	}
	return obj.sts, nil
}

// JSONNew use json data create StatelfulSet
func (obj *StatefulSet) JSONNew(jsonbyts []byte) *StatefulSet {
	obj.err = json.Unmarshal(jsonbyts, obj.sts)
	return obj
}

// YAMLNew use yaml data create StatefulSet
func (obj *StatefulSet) YAMLNew(yamlbyts []byte) *StatefulSet {
	obj.err = yaml.Unmarshal(yamlbyts, obj.sts)
	return obj
}

// SetName set StatefulSet(sts) name
func (obj *StatefulSet) SetName(name string) *StatefulSet {
	obj.sts.SetName(name)
	return obj
}

// SetNameSpace set StatefulSet(sts) namespace ,default namespace is 'default'
func (obj *StatefulSet) SetNameSpace(namespace string) *StatefulSet {
	obj.sts.SetNamespace(namespace)
	return obj
}

// SetLabels set StatefulSet(sts) Labels
func (obj *StatefulSet) SetLabels(labels map[string]string) *StatefulSet {
	obj.sts.SetLabels(labels)
	return obj
}

// SetReplicas set StatefulSet(sts) replicas
func (obj *StatefulSet) SetReplicas(replicas int32) *StatefulSet {
	obj.sts.Spec.Replicas = &replicas
	return obj
}

// SetSelector set StatefulSet(sts) labels selector and set Pod Labels
func (obj *StatefulSet) SetSelector(labels map[string]string) *StatefulSet {
	if len(labels) <= 0 {
		obj.err = errors.New("SetSelector err,label is not allowed to be empty")
		return obj
	}
	if obj.sts.Spec.Selector == nil {
		obj.sts.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: labels,
		}
		obj.sts.Spec.Selector.MatchLabels = labels
		obj.sts.Spec.Template.SetLabels(labels)
		return obj
	}
	obj.sts.Spec.Selector.MatchLabels = labels
	obj.sts.Spec.Template.SetLabels(labels)
	return obj
}

// GetPodLabel get Pod labels
func (obj *StatefulSet) GetPodLabel() map[string]string { return obj.sts.Spec.Template.GetLabels() }

// SetPodLabels set Pod labels and set StatefulSet(sts) selector
func (obj *StatefulSet) SetPodLabels(labels map[string]string) *StatefulSet {
	obj.sts.Spec.Template.SetLabels(labels)
	obj.SetSelector(labels)
	return obj
}

// SetContainer set StatefulSet(sts) container
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
	containersLen := len(obj.sts.Spec.Template.Spec.Containers)
	if containersLen < 1 {
		obj.sts.Spec.Template.Spec.Containers = []corev1.Container{container}
		return obj
	}

	for index := 0; index < containersLen; index++ {
		img := strings.TrimSpace(obj.sts.Spec.Template.Spec.Containers[index].Image)
		if img == "" || len(img) <= 0 {
			obj.sts.Spec.Template.Spec.Containers[index].Name = name
			obj.sts.Spec.Template.Spec.Containers[index].Image = image
			obj.sts.Spec.Template.Spec.Containers[index].Ports = []corev1.ContainerPort{port}
			return obj
		}
	}
	obj.sts.Spec.Template.Spec.Containers = append(obj.sts.Spec.Template.Spec.Containers, container)
	return obj
}

// verify check service necessary value, input the default field and input related data.
func (obj *StatefulSet) verify() {
	if obj.err != nil {
		return
	}
	if !verifyString(obj.sts.GetName()) {
		obj.err = errors.New("StatefulSet.Name is not allowed to be empty")
		return
	}
	if obj.sts.Spec.Selector == nil {
		obj.err = errors.New("StatefulSet.Spec.Selector is not allowed to be empty")
		return
	}
	if len(obj.sts.Spec.Template.Spec.Containers) < 1 {
		obj.err = errors.New("StatefulSet.Container is not allowed to be empty")
		return
	}
}
