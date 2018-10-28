package beku

import (
	"encoding/json"
	"errors"

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

// NewSts  create StatefulSet(sts) and chain function call begin with this function.
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

// SetReplicas set StatefulSet(sts) replicas default 1
func (obj *StatefulSet) SetReplicas(replicas int32) *StatefulSet {
	obj.sts.Spec.Replicas = &replicas
	return obj
}

// SetSelector set StatefulSet(sts) labels selector and set Pod Labels
func (obj *StatefulSet) SetSelector(labels map[string]string) *StatefulSet {
	if len(labels) <= 0 {
		obj.error(errors.New("SetSelector err,label is not allowed to be empty"))
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

// SetAnnotations set StatefulSet annotations
func (obj *StatefulSet) SetAnnotations(annotations map[string]string) *StatefulSet {
	if len(obj.sts.Annotations) <= 0 {
		obj.sts.Annotations = annotations
		return obj
	}
	for key, value := range annotations {
		obj.sts.Annotations[key] = value
	}
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
	obj.error(setContainer(&obj.sts.Spec.Template, name, image, containerPort))
	return obj
}

// SetResourceLimit set container of StatefulSet resource limit,eg:CPU and MEMORY
func (obj *StatefulSet) SetResourceLimit(limits map[ResourceName]string) *StatefulSet {
	obj.error(setResourceLimit(&obj.sts.Spec.Template, limits))
	return obj
}

func (obj *StatefulSet) error(err error) {
	if obj.err != nil {
		return
	}
	obj.err = err
}

// SetResourceRequst set container of StatefulSet resource request,only CPU and MEMORY
func (obj *StatefulSet) SetResourceRequst(requests map[ResourceName]string) *StatefulSet {
	obj.error(setResourceRequests(&obj.sts.Spec.Template, requests))
	return obj
}

// SetEnvs set Pod Environmental variable
func (obj *StatefulSet) SetEnvs(envMap map[string]string) *StatefulSet {
	obj.error(setEnvs(&obj.sts.Spec.Template, envMap))
	return obj
}

// SetPVClaim set StatefulSet PersistentVolumeClaimVolumeSource
// params:
// volumeName: this is Custom field,you can define VolumeSource name,will be used of the container MountPath,
// claimName: this is PersistentVolumeClaim(PVC) name,the PVC and StatefulSet must on same namespace and exist.
func (obj *StatefulSet) SetPVClaim(volumeName, claimName string) *StatefulSet {
	obj.error(setPVClaim(&obj.sts.Spec.Template, volumeName, claimName))
	return obj
}

//SetPVCMounts mount PersistentVolumeClaim on container
// params:
// volumeName:the param is SetPVClaim() function volumeName,and when you call SetPVCMounts function you must call SetPVClaim function,and no order.
// on the other hand SetPVCMounts() function only mount first Container,and On the Container you can volumeMount many PersistentVolumeClaim.
// mounthPath: runtime container dir eg:/var/lib/mysql
func (obj *StatefulSet) SetPVCMounts(volumeName, mounthPath string) *StatefulSet {
	obj.error(setPVCMounts(&obj.sts.Spec.Template, volumeName, mounthPath))
	return nil
}

// SetHTTPLiveness set container liveness of http style
// port: required
// path: http request URL,eg: /api/v1/posts/1
// initDelaySec: how long time after the first start of the program the probe is executed for the first time.(sec)
// timeoutSec: http request timeout seconds,defaults to 1 second. Minimum value is 1.
// periodSec: how often does the probe??defaults to 1 second. Minimum value is 1,Except for the first time?
// headers: headers[0] is HTTP Header, do not fill if you do not need to set
// on the other hand, only **first container** will be set livenessProbe
func (obj *StatefulSet) SetHTTPLiveness(port int, path string, initDelaySec, timeoutSec, periodSec int32, headers ...map[string]string) *StatefulSet {
	setLiveness(&obj.sts.Spec.Template, httpProbe(port, path, initDelaySec, timeoutSec, periodSec, headers...))
	return obj
}

// SetCMDLiveness set container liveness of cmd style
// cmd: execute liveness probe as commond line
// timeoutSec: http request timeout seconds,defaults to 1 second. Minimum value is 1.
// periodSec: how often does the probe??defaults to 1 second. Minimum value is 1,Except for the first time?
// headers: headers[0] is HTTP Header, do not fill if you do not need to set
// on the other hand, only **first container** will be set livenessProbe
func (obj *StatefulSet) SetCMDLiveness(cmd []string, initDelaySec, timeoutSec, periodSec int32) *StatefulSet {
	setLiveness(&obj.sts.Spec.Template, cmdProbe(cmd, initDelaySec, timeoutSec, periodSec))
	return obj
}

// SetTCPLiveness set container liveness of tcp style
// host: default is ""
// port: required
// timeoutSec: http request timeout seconds,defaults to 1 second. Minimum value is 1.
// periodSec: how often does the probe??defaults to 1 second. Minimum value is 1,Except for the first time?
// headers: headers[0] is HTTP Header, do not fill if you do not need to set
// on the other hand, only **first container** will be set livenessProbe
func (obj *StatefulSet) SetTCPLiveness(host string, port int, initDelaySec, timeoutSec, periodSec int32) *StatefulSet {
	setLiveness(&obj.sts.Spec.Template, tcpProbe(host, port, initDelaySec, timeoutSec, periodSec))
	return obj
}

// SetHTTPReadness set container readness
// initDelaySec: how long time after the first start of the program the probe is executed for the first time.(sec)
// timeoutSec: http request timeout seconds,defaults to 1 second. Minimum value is 1.
// periodSec: how often does the probe??defaults to 1 second. Minimum value is 1,Except for the first time?
// on the other hand, only **first container** will be set livenessProbe
func (obj *StatefulSet) SetHTTPReadness(port int, path string, initDelaySec, timeoutSec, periodSec int32, headers ...map[string]string) *StatefulSet {
	setReadness(&obj.sts.Spec.Template, httpProbe(port, path, initDelaySec, timeoutSec, periodSec, headers...))
	return obj
}

// SetCMDReadness set container readness of cmd style
// cmd: execute readness probe as commond line
// timeoutSec: http request timeout seconds,defaults to 1 second. Minimum value is 1.
// periodSec: how often does the probe? defaults to 1 second. Minimum value is 1,Except for the first time?
// headers: headers[0] is HTTP Header, do not fill if you do not need to set
// on the other hand, only **first container** will be set livenessProbe
func (obj *StatefulSet) SetCMDReadness(cmd []string, initDelaySec, timeoutSec, periodSec int32) *StatefulSet {
	setReadness(&obj.sts.Spec.Template, cmdProbe(cmd, initDelaySec, timeoutSec, periodSec))
	return obj
}

// SetTCPReadness set container readness of tcp style
// host: default is ""
// port: required
// timeoutSec: http request timeout seconds,defaults to 1 second. Minimum value is 1.
// periodSec: how often does the probe? defaults to 1 second. Minimum value is 1,Except for the first time?
// headers: headers[0] is HTTP Header, do not fill if you do not need to set
// on the other hand, only **first container** will be set livenessProbe
func (obj *StatefulSet) SetTCPReadness(host string, port int, initDelaySec, timeoutSec, periodSec int32) *StatefulSet {
	setReadness(&obj.sts.Spec.Template, tcpProbe(host, port, initDelaySec, timeoutSec, periodSec))
	return obj
}

// SetPodQos set pod  quality of service
// qosClass: is quality of service,the value only 'Guaranteed','Burstable' and 'BestEffort'
// autoSet: If your previous settings do not meet the requirements of PodQoS, we will automatically set
func (obj *StatefulSet) SetPodQos(qosClass string, autoSet ...bool) *StatefulSet {
	obj.SetAnnotations(setQosMap(obj.sts.Annotations, qosClass, autoSet...))
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
	//check qos set,if err!=nil, check need auto set qos
	presentQos, err := qosCheck(obj.sts.Annotations[qosKey], obj.sts.Spec.Template.Spec)
	if err != nil {
		if obj.sts.Annotations[autoQosKey] == "true" {
			err := obj.autoSetQos(presentQos)
			if err != nil {
				obj.err = err
				return
			}
		} else {
			obj.err = err
			return
		}
	}
	for index := range obj.sts.Spec.Template.Spec.Containers {
		obj.sts.Spec.Template.Spec.Containers[index].ImagePullPolicy = corev1.PullIfNotPresent
	}
	obj.sts.Kind = "StatefulSet"
	obj.sts.APIVersion = "apps/v1"
}

// autoSetQos auto set Pod of StatefulSet QOS
func (obj *StatefulSet) autoSetQos(presentQos string) error {
	return autoSetQos(obj.sts.Annotations[qosKey], presentQos, &obj.sts.Spec.Template.Spec)
}
