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
