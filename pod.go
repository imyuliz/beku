package beku

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/ghodss/yaml"
	"k8s.io/api/core/v1"
)

// Pod include Kubernetes resource bject Pod and error
type Pod struct {
	pod *v1.Pod
	err error
}

// NewPod create Pod and hain function call begin with this function.
func NewPod() *Pod { return &Pod{pod: &v1.Pod{}} }

// JSONNew use json data create Pod
func (obj *Pod) JSONNew(jsonbyts []byte) *Pod {
	obj.error(json.Unmarshal(jsonbyts, obj.pod))
	return obj
}

// YAMLNew use yaml data create Pod
func (obj *Pod) YAMLNew(yamlbyts []byte) *Pod {
	obj.error(yaml.Unmarshal(yamlbyts, obj.pod))
	return obj
}

// Finish Chain function call end with this function
// return Kubernetes resource object Pod and error.
// In the function, it will check necessary parametersainput the default field
func (obj *Pod) Finish() (pod *v1.Pod, err error) {
	obj.verify()
	pod, err = obj.pod, obj.err
	return
}

// SetName set Pod name
func (obj *Pod) SetName(name string) *Pod {
	obj.pod.SetName(name)
	return obj
}

// SetNamespace set Pod namespace and set Pod namespace.
func (obj *Pod) SetNamespace(namespace string) *Pod {
	obj.pod.SetNamespace(namespace)
	return obj
}

// SetNamespaceAndName set Pod Namespace and Pod name
func (obj *Pod) SetNamespaceAndName(namespace, name string) *Pod {
	obj.SetName(name)
	obj.SetNamespace(namespace)
	return obj
}

// SetLabels set pod labels
func (obj *Pod) SetLabels(labels map[string]string) *Pod {
	obj.pod.SetLabels(labels)
	return obj
}

// SetContainer set pod container
// name:name is container name ,default ""
// image:image is image name ,must input image
// containerPort: image expose containerPort,must input containerPort
func (obj *Pod) SetContainer(name, image string, containerPort int32) *Pod {
	// This must be a valid port number, 0 < x < 65536.
	if containerPort <= 0 || containerPort >= 65536 {
		obj.error(errors.New("SetContainer err, container Port range: 0 < containerPort < 65536"))
		return obj
	}
	if !verifyString(image) {
		obj.error(errors.New("SetContainer err, image is not allowed to be empty"))
		return obj

	}
	port := v1.ContainerPort{ContainerPort: containerPort}
	container := v1.Container{
		Name:  name,
		Image: image,
		Ports: []v1.ContainerPort{port},
	}
	containersLen := len(obj.pod.Spec.Containers)
	if containersLen < 1 {
		obj.pod.Spec.Containers = []v1.Container{container}
		return obj
	}
	for index := 0; index < containersLen; index++ {
		img := strings.TrimSpace(obj.pod.Spec.Containers[index].Image)
		if img == "" || len(img) <= 0 {
			obj.pod.Spec.Containers[index].Name = name
			obj.pod.Spec.Containers[index].Image = image
			obj.pod.Spec.Containers[index].Ports = []v1.ContainerPort{port}
			return obj
		}
	}
	obj.pod.Spec.Containers = append(obj.pod.Spec.Containers, container)
	return obj
}

func (obj *Pod) error(err error) {
	if obj.err != nil {
		return
	}
	obj.err = err
}

func (obj *Pod) verify() {
	if !verifyString(obj.pod.GetName()) {
		obj.err = errors.New("Pod name is not allowed to be empty")
		return
	}
	obj.error(containerRepeated(obj.pod.Spec.Containers))
	obj.pod.Kind = "Pod"
	obj.pod.APIVersion = "v1"
}
