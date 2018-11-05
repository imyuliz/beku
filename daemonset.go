package beku

import (
	"encoding/json"
	"errors"

	"github.com/ghodss/yaml"
	"k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DaemonSet include Kubernets resource object DaemonSet and error
type DaemonSet struct {
	ds  *v1.DaemonSet
	err error
}

// NewDS create DaemonSet(ds) and chain function call begin with this function.
func NewDS() *DaemonSet { return &DaemonSet{ds: &v1.DaemonSet{}} }

// Finish Chain function call end with this function
// return real DaemonSet(really DaemonSet is kubernetes resource object DaemonSet and error
// In the function, it will check necessary parameters、input the default field
func (obj *DaemonSet) Finish() (*v1.DaemonSet, error) {
	obj.verify()
	return obj.ds, obj.err
}

// JSONNew use json data create DaemonSet
func (obj *DaemonSet) JSONNew(jsonbyts []byte) *DaemonSet {
	obj.error(json.Unmarshal(jsonbyts, obj.ds))
	return obj
}

// YAMLNew use yaml data create DaemonSet
func (obj *DaemonSet) YAMLNew(yamlbyts []byte) *DaemonSet {
	obj.error(yaml.Unmarshal(yamlbyts, obj.ds))
	return obj
}

// Replace replace ds by Kubernetes resource object
func (obj *DaemonSet) Replace(ds *v1.DaemonSet) *DaemonSet {
	if ds != nil {
		obj.ds = ds
	}
	return obj
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

// SetLabels set DaemonSet(ds) Labels,set Pod Labels.
func (obj *DaemonSet) SetLabels(labels map[string]string) *DaemonSet {
	obj.ds.SetLabels(labels)
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
	obj.SetSelector(labels)
	return obj
}

// SetContainer set DaemonSet container
// name Not required when only one Container,you can input "".
// when many container this Field is necessary and cann't repeat
// image is necessary, image very important
// containerPort container port,this is necessary
func (obj *DaemonSet) SetContainer(name, image string, containerPort int32) *DaemonSet {
	obj.error(setContainer(&obj.ds.Spec.Template, name, image, containerPort))
	return obj
}

// SetAnnotations set DaemonSet annotations
func (obj *DaemonSet) SetAnnotations(annotations map[string]string) *DaemonSet {
	if len(obj.ds.Annotations) <= 0 {
		obj.ds.Annotations = annotations
		return obj
	}
	for key, value := range annotations {
		obj.ds.Annotations[key] = value
	}
	return obj
}

// SetPodQos set pod  quality of service
// qosClass: is quality of service,the value only 'Guaranteed','Burstable' and 'BestEffort'
// autoSet: If your previous settings do not meet the requirements of PodQoS, we will automatically set
func (obj *DaemonSet) SetPodQos(qosClass string, autoSet ...bool) *DaemonSet {
	obj.SetAnnotations(setQosMap(obj.ds.Annotations, qosClass, autoSet...))
	return obj
}

// SetHTTPLiveness set container liveness of http style
// port: required
// path: http request URL,eg: /api/v1/posts/1
// initDelaySec: how long time after the first start of the program the probe is executed for the first time.(sec)
// timeoutSec: http request timeout seconds,defaults to 1 second. Minimum value is 1.
// periodSec: how often does the probe??defaults to 1 second. Minimum value is 1,Except for the first time?
// headers: headers[0] is HTTP Header, do not fill if you do not need to set
// on the other hand, only **first container** will be set livenessProbe
func (obj *DaemonSet) SetHTTPLiveness(port int, path string, initDelaySec, timeoutSec, periodSec int32, headers ...map[string]string) *DaemonSet {
	setLiveness(&obj.ds.Spec.Template, httpProbe(port, path, initDelaySec, timeoutSec, periodSec, headers...))
	return obj
}

// SetCMDLiveness set container liveness of cmd style
// cmd: execute liveness probe as commond line
// timeoutSec: http request timeout seconds,defaults to 1 second. Minimum value is 1.
// periodSec: how often does the probe??defaults to 1 second. Minimum value is 1,Except for the first time?
// headers: headers[0] is HTTP Header, do not fill if you do not need to set
// on the other hand, only **first container** will be set livenessProbe
func (obj *DaemonSet) SetCMDLiveness(cmd []string, initDelaySec, timeoutSec, periodSec int32) *DaemonSet {
	setLiveness(&obj.ds.Spec.Template, cmdProbe(cmd, initDelaySec, timeoutSec, periodSec))
	return obj
}

// SetTCPLiveness set container liveness of tcp style
// host: default is ""
// port: required
// timeoutSec: http request timeout seconds,defaults to 1 second. Minimum value is 1.
// periodSec: how often does the probe??defaults to 1 second. Minimum value is 1,Except for the first time?
// headers: headers[0] is HTTP Header, do not fill if you do not need to set
// on the other hand, only **first container** will be set livenessProbe
func (obj *DaemonSet) SetTCPLiveness(host string, port int, initDelaySec, timeoutSec, periodSec int32) *DaemonSet {
	setLiveness(&obj.ds.Spec.Template, tcpProbe(host, port, initDelaySec, timeoutSec, periodSec))
	return obj
}

// SetHTTPReadness set container readness
// initDelaySec: how long time after the first start of the program the probe is executed for the first time.(sec)
// timeoutSec: http request timeout seconds,defaults to 1 second. Minimum value is 1.
// periodSec: how often does the probe??defaults to 1 second. Minimum value is 1,Except for the first time?
// on the other hand, only **first container** will be set livenessProbe
func (obj *DaemonSet) SetHTTPReadness(port int, path string, initDelaySec, timeoutSec, periodSec int32, headers ...map[string]string) *DaemonSet {
	setReadness(&obj.ds.Spec.Template, httpProbe(port, path, initDelaySec, timeoutSec, periodSec, headers...))
	return obj
}

// SetCMDReadness set container readness of cmd style
// cmd: execute readness probe as commond line
// timeoutSec: http request timeout seconds,defaults to 1 second. Minimum value is 1.
// periodSec: how often does the probe? defaults to 1 second. Minimum value is 1,Except for the first time?
// headers: headers[0] is HTTP Header, do not fill if you do not need to set
// on the other hand, only **first container** will be set livenessProbe
func (obj *DaemonSet) SetCMDReadness(cmd []string, initDelaySec, timeoutSec, periodSec int32) *DaemonSet {
	setReadness(&obj.ds.Spec.Template, cmdProbe(cmd, initDelaySec, timeoutSec, periodSec))
	return obj
}

// SetTCPReadness set container readness of tcp style
// host: default is ""
// port: required
// timeoutSec: http request timeout seconds,defaults to 1 second. Minimum value is 1.
// periodSec: how often does the probe? defaults to 1 second. Minimum value is 1,Except for the first time?
// headers: headers[0] is HTTP Header, do not fill if you do not need to set
// on the other hand, only **first container** will be set livenessProbe
func (obj *DaemonSet) SetTCPReadness(host string, port int, initDelaySec, timeoutSec, periodSec int32) *DaemonSet {
	setReadness(&obj.ds.Spec.Template, tcpProbe(host, port, initDelaySec, timeoutSec, periodSec))
	return obj
}

// SetPVClaim set DaemonSet PersistentVolumeClaimVolumeSource
// params:
// volumeName: this is Custom field,you can define VolumeSource name,will be used of the container MountPath,
// claimName: this is PersistentVolumeClaim(PVC) name,the PVC and DaemonSet must on same namespace and exist.
func (obj *DaemonSet) SetPVClaim(volumeName, claimName string) *DaemonSet {
	obj.error(setPVClaim(&obj.ds.Spec.Template, volumeName, claimName))
	return obj
}

//SetPVCMounts mount PersistentVolumeClaim on container
// params:
// volumeName:the param is SetPVClaim() function volumeName,and when you call SetPVCMounts function you must call SetPVClaim function,and no order.
// on the other hand SetPVCMounts() function only mount first Container,and On the Container you can volumeMount many PersistentVolumeClaim.
// mountPath: runtime container dir eg:/var/lib/mysql
func (obj *DaemonSet) SetPVCMounts(volumeName, mountPath string) *DaemonSet {
	obj.error(setPVCMounts(&obj.ds.Spec.Template, volumeName, mountPath))
	return nil
}

// SetEnvs set Pod Environmental variable
func (obj *DaemonSet) SetEnvs(envMap map[string]string) *DaemonSet {
	obj.error(setEnvs(&obj.ds.Spec.Template, envMap))
	return obj
}

// SetMinReadySeconds set DaemonSet minreadyseconds default 600
func (obj *DaemonSet) SetMinReadySeconds(sec int32) *DaemonSet {
	if sec < 0 {
		sec = 0
	}
	obj.ds.Spec.MinReadySeconds = sec
	return obj
}

// SetHistoryLimit set DaemonSet history version numbers, limit default 10
// the field is used to Rollback
func (obj *DaemonSet) SetHistoryLimit(limit int32) *DaemonSet {
	if limit <= 0 {
		limit = 10
	}
	obj.ds.Spec.RevisionHistoryLimit = &limit
	return obj
}

// SetImagePullSecrets set pod pull secret
func (obj *DaemonSet) SetImagePullSecrets(secretName string) *DaemonSet {
	setImagePullSecrets(&obj.ds.Spec.Template, secretName)
	return obj
}

// Release release DaemonSet on Kubernetes
func (obj *DaemonSet) Release() (*v1.DaemonSet, error) {
	ds, err := obj.Finish()
	if err != nil {
		return nil, err
	}
	client, err := GetKubeClient()
	if err != nil {
		return nil, err
	}
	return client.AppsV1().DaemonSets(ds.GetNamespace()).Create(ds)
}

// GetPodLabel get pod labels
func (obj *DaemonSet) GetPodLabel() map[string]string {
	return obj.ds.Spec.Template.GetLabels()
}

func (obj *DaemonSet) error(err error) {
	if obj.err != nil {
		return
	}
	obj.err = err
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
		return
	}
	//check qos set,if err!=nil, check need auto set qos
	presentQos, err := qosCheck(obj.ds.Annotations[qosKey], obj.ds.Spec.Template.Spec)
	if err != nil {
		if obj.ds.Annotations[autoQosKey] == "true" {
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
	obj.ds.Kind = "DaemonSet"
	obj.ds.APIVersion = "app/v1"
}

// autoSetQos auto set Pod of Deployment QOS
func (obj *DaemonSet) autoSetQos(presentQos string) error {
	return autoSetQos(obj.ds.Annotations[qosKey], presentQos, &obj.ds.Spec.Template.Spec)
}
