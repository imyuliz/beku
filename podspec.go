package beku

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/sets"
)

func setImagePullSecrets(podTemp *v1.PodTemplateSpec, secretName string) {
	if len(podTemp.Spec.ImagePullSecrets) <= 0 {
		podTemp.Spec.ImagePullSecrets = []v1.LocalObjectReference{v1.LocalObjectReference{Name: secretName}}
		return
	}
	podTemp.Spec.ImagePullSecrets = append(podTemp.Spec.ImagePullSecrets, v1.LocalObjectReference{Name: secretName})
}

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

var supportedQoSComputeResources = sets.NewString(string(ResourceCPU), string(ResourceMemory))

// QOSList is a set of (resource name, QoS class) pairs.
type QOSList map[v1.ResourceName]v1.PodQOSClass

func isSupportedQoSComputeResource(name v1.ResourceName) bool {
	return supportedQoSComputeResources.Has(string(name))
}

// GetPodQOS returns the QoS class of a pod.
// A pod is besteffort if none of its containers have specified any requests or limits.
// A pod is guaranteed only when requests and limits are specified for all the containers and they are equal.
// A pod is burstable if limits and requests do not match across all containers.
func GetPodQOS(pod v1.PodSpec) v1.PodQOSClass {
	requests := v1.ResourceList{}
	limits := v1.ResourceList{}
	zeroQuantity := resource.MustParse("0")
	isGuaranteed := true
	for _, container := range pod.Containers {
		// process requests
		for name, quantity := range container.Resources.Requests {
			if !isSupportedQoSComputeResource(name) {
				continue
			}
			if quantity.Cmp(zeroQuantity) == 1 {
				delta := quantity.Copy()
				if _, exists := requests[name]; !exists {
					requests[name] = *delta
				} else {
					delta.Add(requests[name])
					requests[name] = *delta
				}
			}
		}
		// process limits
		qosLimitsFound := sets.NewString()
		for name, quantity := range container.Resources.Limits {
			if !isSupportedQoSComputeResource(name) {
				continue
			}
			if quantity.Cmp(zeroQuantity) == 1 {
				qosLimitsFound.Insert(string(name))
				delta := quantity.Copy()
				if _, exists := limits[name]; !exists {
					limits[name] = *delta
				} else {
					delta.Add(limits[name])
					limits[name] = *delta
				}
			}
		}

		if !qosLimitsFound.HasAll(string(v1.ResourceMemory), string(v1.ResourceCPU)) {
			isGuaranteed = false
		}
	}
	if len(requests) == 0 && len(limits) == 0 {
		return v1.PodQOSBestEffort
	}
	// Check is requests match limits for all resources.
	if isGuaranteed {
		for name, req := range requests {
			if lim, exists := limits[name]; !exists || lim.Cmp(req) != 0 {
				isGuaranteed = false
				break
			}
		}
	}
	if isGuaranteed &&
		len(requests) == len(limits) {
		return v1.PodQOSGuaranteed
	}
	return v1.PodQOSBurstable
}

func autoSetQos(targetQos, presentQos string, pod *v1.PodSpec) error {
	if qosRanks[presentQos] >= qosRanks[targetQos] {
		return nil
	}
	if qosRanks[targetQos] == GuaranteedRank &&
		qosRanks[presentQos] == BestEffortRank {
		if len(defaultLimit()) == 2 && len(defaultRequest()) == 2 {
			if reflect.DeepEqual(defaultLimit(), defaultRequest()) {
				containers := len(pod.Containers)
				requests, _ := ResourceMapsToK8s(defaultRequest())
				for index := 0; index < containers; index++ {
					pod.Containers[index].Resources.Limits = requests
					pod.Containers[index].Resources.Requests = requests
				}
				return nil
			}
			return fmt.Errorf("set QOS rank:%s failed,because,Because the default addition is not satisfied,notice:%s", targetQos, qosNotices[targetQos])
		}
		return fmt.Errorf("set QOS rank:%s failed,because,Because the default addition is not satisfied,you can call func RegisterResourceLimit() and RegisterResourceRequest() register default resource limits and requests", targetQos)
	}

	//If what you expect is Burstable
	if len(defaultRequest()) > 0 {
		//set container of Pod resoource requests value.
		containers := len(pod.Containers)
		requests, _ := ResourceMapsToK8s(defaultRequest())
		for index := 0; index < containers; index++ {
			pod.Containers[index].Resources.Requests = requests
		}
		return nil
	}
	return errors.New("set Qos Rank failed,you can call func RegisterResourceLimit() and RegisterResourceRequest() register default resource limits and requests")
}

func qosCheck(qosClass string, podTem v1.PodSpec) (string, error) {
	qosClass = strings.TrimSpace(qosClass)
	if qosClass == "" || qosClass == "BestEffort" {
		return "BestEffort", nil
	}
	//Kubernetes evaluate qos grade
	evaQos := string(GetPodQOS(podTem))
	if qosClass == evaQos {
		return evaQos, nil
	}
	return evaQos, fmt.Errorf("qos check failed, notice:%s", qosNotices[qosClass])
}

// setQosMap set Pod Qos
func setQosMap(dec map[string]string, qosClass string, autoSet ...bool) map[string]string {
	var (
		auto = "false"
	)
	if len(autoSet) > 0 && autoSet[0] {
		auto = "true"

	}
	if dec == nil {
		dec = map[string]string{qosKey: qosClass, autoQosKey: auto}
		return dec
	}
	dec[qosKey] = qosClass
	dec[autoQosKey] = auto
	return dec
}
