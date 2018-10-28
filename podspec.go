package beku

import (
	"errors"
	"fmt"
	"strings"

	"k8s.io/api/core/v1"
)

// setContainer set container
func setContainer(podTemp *v1.PodTemplateSpec, name, image string, containerPort int32) error {
	// This must be a valid port number, 0 < x < 65536.
	if containerPort <= 0 || containerPort >= 65536 {
		return errors.New("SetContainer err, container Port range: 0 < containerPort < 65536")
	}
	if !verifyString(image) {
		return errors.New("SetContainer err, image is not allowed to be empty")

	}
	port := v1.ContainerPort{ContainerPort: containerPort}
	container := v1.Container{
		Name:  name,
		Image: image,
		Ports: []v1.ContainerPort{port},
	}
	containersLen := len(podTemp.Spec.Containers)
	if containersLen < 1 {
		podTemp.Spec.Containers = []v1.Container{container}
		return nil
	}
	for index := 0; index < containersLen; index++ {
		img := strings.TrimSpace(podTemp.Spec.Containers[index].Image)
		if img == "" || len(img) <= 0 {
			podTemp.Spec.Containers[index].Name = name
			podTemp.Spec.Containers[index].Image = image
			podTemp.Spec.Containers[index].Ports = []v1.ContainerPort{port}
			return nil
		}
	}
	podTemp.Spec.Containers = append(podTemp.Spec.Containers, container)
	return nil
}

func setResourceLimit(podTemp *v1.PodTemplateSpec, limits map[ResourceName]string) error {
	data, err := ResourceMapsToK8s(limits)
	if err != nil {
		return fmt.Errorf("SetResourceLimit err:%v", err)
	}
	containerLen := len(podTemp.Spec.Containers)
	if containerLen < 1 {
		podTemp.Spec.Containers = []v1.Container{v1.Container{Resources: v1.ResourceRequirements{Limits: data}}}
		return nil
	}
	for index := 0; index < containerLen; index++ {
		if podTemp.Spec.Containers[index].Resources.Limits == nil {
			podTemp.Spec.Containers[index].Resources.Limits = data
		}
	}
	return nil
}

func setResourceRequests(podTemp *v1.PodTemplateSpec, requests map[ResourceName]string) error {
	data, err := ResourceMapsToK8s(requests)
	if err != nil {
		return fmt.Errorf("SetResourceLimit err:%v", err)
	}
	containerLen := len(podTemp.Spec.Containers)
	if containerLen < 1 {
		podTemp.Spec.Containers = []v1.Container{v1.Container{Resources: v1.ResourceRequirements{Requests: data}}}
		return nil
	}
	for index := 0; index < containerLen; index++ {
		if podTemp.Spec.Containers[index].Resources.Requests == nil {
			podTemp.Spec.Containers[index].Resources.Requests = data
		}
	}
	return nil
}

func setEnvs(podTemp *v1.PodTemplateSpec, envMap map[string]string) error {
	envs, err := mapToEnvs(envMap)
	if err != nil {
		return err
	}
	containerLen := len(podTemp.Spec.Containers)
	if containerLen < 1 {
		podTemp.Spec.Containers = []v1.Container{v1.Container{Env: envs}}
		return nil
	}
	for index := 0; index < containerLen; index++ {
		if podTemp.Spec.Containers[index].Env == nil {
			podTemp.Spec.Containers[index].Env = envs
		}
	}
	return nil
}

func setPVCMounts(podTemp *v1.PodTemplateSpec, volumeName, mountPath string) error {
	volumeMount := v1.VolumeMount{Name: volumeName, MountPath: mountPath}
	if len(podTemp.Spec.Containers) <= 0 {
		podTemp.Spec.Containers = append(podTemp.Spec.Containers, v1.Container{
			VolumeMounts: []v1.VolumeMount{volumeMount},
		})
		return nil
	}
	//only mount first container and first container can mount many data source.
	if len(podTemp.Spec.Containers[0].VolumeMounts) <= 0 {
		podTemp.Spec.Containers[0].VolumeMounts = []v1.VolumeMount{volumeMount}
		return nil
	}
	podTemp.Spec.Containers[0].VolumeMounts = append(podTemp.Spec.Containers[0].VolumeMounts, volumeMount)
	return nil
}

func setPVClaim(podTemp *v1.PodTemplateSpec, volumeName, claimName string) error {
	volume := v1.Volume{
		Name: volumeName,
		VolumeSource: v1.VolumeSource{
			PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
				ClaimName: claimName,
				ReadOnly:  false,
			},
		},
	}
	if len(podTemp.Spec.Volumes) <= 0 {
		podTemp.Spec.Volumes = []v1.Volume{volume}
		return nil
	}
	podTemp.Spec.Volumes = append(podTemp.Spec.Volumes, volume)
	return nil
}

func setLiveness(podTemp *v1.PodTemplateSpec, probe *v1.Probe) error {
	if len(podTemp.Spec.Containers) <= 0 {
		podTemp.Spec.Containers = []v1.Container{v1.Container{LivenessProbe: probe}}
		return nil
	}
	podTemp.Spec.Containers[0].LivenessProbe = probe
	return nil
}
func setReadness(podTemp *v1.PodTemplateSpec, probe *v1.Probe) error {
	if len(podTemp.Spec.Containers) <= 0 {
		podTemp.Spec.Containers = []v1.Container{v1.Container{ReadinessProbe: probe}}
		return nil
	}
	podTemp.Spec.Containers[0].ReadinessProbe = probe
	return nil
}
