package beku

import (
	"errors"
	"fmt"
	"strings"

	"github.com/yulibaozi/beku/core"
	"github.com/yulibaozi/mapper"
	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Deployment api/apps/deploy deployment
type Deployment struct {
	deploy *v1.Deployment
	err    error
}

// NewDeployment create deployment
func NewDeployment() *Deployment {
	return &Deployment{
		deploy: &v1.Deployment{},
	}
}

// SetName set deployment name
func (dep *Deployment) SetName(name string) *Deployment {
	dep.deploy.SetName(name)
	return dep
}

// SetNameSpace set deployment namespace
func (dep *Deployment) SetNameSpace(namespace string) *Deployment {
	dep.deploy.SetNamespace(namespace)
	dep.deploy.Spec.Template.SetNamespace(namespace)
	return dep
}

// SetAnnotations set deployment annotations
func (dep *Deployment) SetAnnotations(annotations map[string]string) *Deployment {
	dep.deploy.SetAnnotations(annotations)
	return dep
}

// SetLabels set deployment labels
// both deployment and pod labels will set
func (dep *Deployment) SetLabels(labels map[string]string) *Deployment {
	dep.deploy.SetLabels(labels)
	dep.deploy.Spec.Template.SetLabels(labels)
	return dep
}

// SetReplicas set replicas default 1
func (dep *Deployment) SetReplicas(replicas int32) *Deployment {
	dep.deploy.Spec.Replicas = &replicas
	return dep
}

// SetSelector set deployment selector
func (dep *Deployment) SetSelector(labels map[string]string) *Deployment {
	if len(labels) <= 0 {
		dep.err = errors.New("LabelSelector set  error,label not allow empty")
		return dep
	}
	if dep.deploy.Spec.Selector == nil {
		dep.deploy.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: labels,
		}
		return dep
	}
	dep.deploy.Spec.Selector.MatchLabels = labels
	return dep
}

// SetMinReadySeconds set  minreadyseconds  default 600
func (dep *Deployment) SetMinReadySeconds(sec int32) *Deployment {
	if sec < 0 {
		sec = 0
	}
	dep.deploy.Spec.MinReadySeconds = sec
	return dep
}

// SetHistoryLimit set history limit default 10
func (dep *Deployment) SetHistoryLimit(limit int32) *Deployment {
	if limit < 0 {
		limit = 10
	}
	dep.deploy.Spec.RevisionHistoryLimit = &limit
	return dep
}

// SetMatchExpressions set match expressions
func (dep *Deployment) SetMatchExpressions(ents []core.LabelSelectorRequirement) *Deployment {
	requirements := make([]metav1.LabelSelectorRequirement, 0)
	err := mapper.AutoMapper(ents, requirements)
	if err != nil {
		dep.err = fmt.Errorf("SetMatchExpressions error:%v", err)
		return dep
	}
	if dep.deploy.Spec.Selector == nil {
		dep.deploy.Spec.Selector = &metav1.LabelSelector{
			MatchExpressions: requirements,
		}
		return dep
	}
	dep.deploy.Spec.Selector.MatchExpressions = requirements
	return dep
}

// SetMaxDeployTime set deploy max time,default 600
func (dep *Deployment) SetMaxDeployTime(sec int32) *Deployment {
	if sec < 0 {
		sec = 600
	}
	dep.deploy.Spec.ProgressDeadlineSeconds = &sec
	return dep
}

// SetPodLabels set pod labels
func (dep *Deployment) SetPodLabels(labels map[string]string) *Deployment {
	dep.deploy.Spec.Template.SetLabels(labels)
	return dep
}

// GetPodLabel get pod labels
func (dep *Deployment) GetPodLabel() map[string]string {
	return dep.deploy.Spec.Template.GetLabels()
}

// SetContainer set deployment container
// name:name is container name ,default ""
// image:image is image name ,must input image
// containerPort: image expose containerPort,must input containerPort
func (dep *Deployment) SetContainer(name, image string, containerPort int32) *Deployment {
	// This must be a valid port number, 0 < x < 65536.
	if containerPort <= 0 || containerPort >= 65536 {
		dep.err = errors.New("containerPort error when SetContainer: 0 < containerPort < 65536")
		return dep
	}
	if !verifyString(image) {
		dep.err = errors.New("image not allow empty,must input image")
		return dep
	}
	port := corev1.ContainerPort{
		ContainerPort: containerPort,
	}
	container := corev1.Container{
		Name:  name,
		Image: image,
		Ports: []corev1.ContainerPort{port},
	}
	if dep.deploy.Spec.Template.Spec.Containers == nil {
		dep.deploy.Spec.Template.Spec.Containers = []corev1.Container{container}
		return dep
	}
	containersLen := len(dep.deploy.Spec.Template.Spec.Containers)
	for index := 0; index < containersLen; index++ {
		img := strings.TrimSpace(dep.deploy.Spec.Template.Spec.Containers[index].Image)
		if img == "" || len(img) <= 0 {
			dep.deploy.Spec.Template.Spec.Containers[index] = container
			return dep
		}
	}
	dep.deploy.Spec.Template.Spec.Containers = append(dep.deploy.Spec.Template.Spec.Containers, container)
	return dep
}

// SetEnvs set envs
func (dep *Deployment) SetEnvs(envMap map[string]string) *Deployment {
	if len(envMap) <= 0 {
		dep.err = errors.New("set env error: envMap is empty")
		return dep
	}
	var envs []corev1.EnvVar
	for k, v := range envMap {
		k, v = strings.TrimSpace(k), strings.TrimSpace(v)
		if k == "" || v == "" {
			dep.err = errors.New("set env error: name or value not allow")
			return dep
		}
		envs = append(envs, corev1.EnvVar{Name: k, Value: v})
	}
	if len(envs) <= 0 {
		dep.err = errors.New("set env error, envs is empty")
		return dep
	}
	containerLen := len(dep.deploy.Spec.Template.Spec.Containers)
	for index := 0; index < containerLen; index++ {
		if dep.deploy.Spec.Template.Spec.Containers[index].Env == nil {
			dep.deploy.Spec.Template.Spec.Containers[index].Env = envs
		}
	}
	return dep
}

func (dep *Deployment) verify() {
	if !verifyString(dep.deploy.GetName()) {
		dep.err = errors.New("deployment name not allow empty")
		return
	}
	if len(dep.deploy.GetLabels()) < 1 {
		dep.err = errors.New("deployment labels not allow empty")
		return
	}
	// if dep.deploy.Spec.Selector == nil {
	// 	dep.err = errors.New("deployment.Spec.Selector not allow empty")
	// 	return dep
	// }
	if len(dep.deploy.Spec.Template.GetLabels()) < 1 {
		dep.err = errors.New("deployment.Spec.Templata.label not allow empty")
		return
	}
	if dep.deploy.Spec.Template.Spec.Containers == nil || len(dep.deploy.Spec.Template.Spec.Containers) < 1 {
		dep.err = errors.New("Deployment.Spec.Template.Spec.Containers not allow nil")
		return
	}
	if dep.deploy.Spec.Selector == nil {
		dep.SetSelector(dep.GetPodLabel())
	}
	dep.deploy.Kind = "Deployment"
	dep.deploy.APIVersion = "apps/v1"
	return
}

// Finish finnaly used
func (dep *Deployment) Finish() (*v1.Deployment, error) {
	dep.verify()
	if dep.err != nil {
		return nil, dep.err
	}
	return dep.deploy, nil
}
