package beku

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ghodss/yaml"
	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Deployment include Kubernetes resource object Deployment and error
type Deployment struct {
	dp  *v1.Deployment
	err error
}

// NewDeployment create Deployment and Chain function call begin with this function.
func NewDeployment() *Deployment { return &Deployment{dp: &v1.Deployment{}} }

// Finish Chain function call end with this function
// return Kubernetes resource object Deployment and error.
// In the function, it will check necessary parametersainput the default field
func (obj *Deployment) Finish() (dp *v1.Deployment, err error) {
	obj.verify()
	dp, err = obj.dp, obj.err
	return
}

// JSONNew use json data create Deployment
func (obj *Deployment) JSONNew(jsonbyts []byte) *Deployment {
	obj.error(json.Unmarshal(jsonbyts, obj.dp))
	return obj
}

// YAMLNew use yaml data create Deployment
func (obj *Deployment) YAMLNew(yamlbyts []byte) *Deployment {
	obj.error(yaml.Unmarshal(yamlbyts, obj.dp))
	return obj
}

// Replace replace Deployment by Kubernetes resource object
func (obj *Deployment) Replace(dp *v1.Deployment) *Deployment {
	if dp != nil {
		obj.dp = dp
	}
	return obj
}

// SetName set Deployment name
func (obj *Deployment) SetName(name string) *Deployment {
	obj.dp.SetName(name)
	return obj
}

// SetNamespace set Deployment namespace and set Pod namespace.
func (obj *Deployment) SetNamespace(namespace string) *Deployment {
	obj.dp.SetNamespace(namespace)
	obj.dp.Spec.Template.SetNamespace(namespace)
	return obj
}

// SetNamespaceAndName set Deployment namespace,set Pod namespace,set Deployment name.
func (obj *Deployment) SetNamespaceAndName(namespace, name string) *Deployment {
	obj.SetNamespace(namespace)
	obj.SetName(name)
	return obj
}

// SetLabels set Deployment labels
func (obj *Deployment) SetLabels(labels map[string]string) *Deployment {
	obj.dp.SetLabels(labels)
	return obj
}

// SetSelector set Deployment selector
// set:
// 1. Deployment.Spec.Selector
// 2. Deployment.Spec.Template.Label(the Field is Pod Labels.)
// and you can not be SetLabels
func (obj *Deployment) SetSelector(labels map[string]string) *Deployment {
	if len(labels) <= 0 {
		obj.error(errors.New("SetSelector err,label is not allowed to be empty"))
		return obj
	}
	if obj.dp.Spec.Selector == nil {
		obj.dp.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: labels,
		}
		obj.dp.Spec.Template.SetLabels(labels)
		return obj
	}
	obj.dp.Spec.Template.SetLabels(labels)
	obj.dp.Spec.Selector.MatchLabels = labels
	return obj
}

// SetAnnotations set Deployment annotations
func (obj *Deployment) SetAnnotations(annotations map[string]string) *Deployment {
	if len(obj.dp.Annotations) <= 0 {
		obj.dp.Annotations = annotations
		return obj
	}
	for key, value := range annotations {
		obj.dp.Annotations[key] = value
	}
	return obj
}

// SetReplicas set Deployment replicas default 1
func (obj *Deployment) SetReplicas(replicas int32) *Deployment {
	obj.dp.Spec.Replicas = &replicas
	return obj
}

// SetMinReadySeconds set Deployment minreadyseconds default 600
func (obj *Deployment) SetMinReadySeconds(sec int32) *Deployment {
	if sec < 0 {
		sec = 0
	}
	obj.dp.Spec.MinReadySeconds = sec
	return obj
}

// SetHistoryLimit set Deployment history version numbers, limit default 10
// the field is used to Rollback
func (obj *Deployment) SetHistoryLimit(limit int32) *Deployment {
	if limit <= 0 {
		limit = 10
	}
	obj.dp.Spec.RevisionHistoryLimit = &limit
	return obj
}

// SetHTTPLiveness set container liveness of http style
// port: required
// path: http request URL,eg: /api/v1/posts/1
// initDelaySec: how long time after the first start of the program the probe is executed for the first time.(sec)
// timeoutSec: http request timeout seconds,defaults to 1 second. Minimum value is 1.
// periodSec: how often does the probe? defaults to 1 second. Minimum value is 1,Except for the first time?
// headers: headers[0] is HTTP Header, do not fill if you do not need to set
// on the other hand, only **first container** will be set livenessProbe
func (obj *Deployment) SetHTTPLiveness(port int, path string, initDelaySec, timeoutSec, periodSec int32, headers ...map[string]string) *Deployment {
	setLiveness(&obj.dp.Spec.Template, httpProbe(port, path, initDelaySec, timeoutSec, periodSec, headers...))
	return obj
}

// SetCMDLiveness set container liveness of cmd style
// cmd: execute liveness probe as commond line
// timeoutSec: http request timeout seconds,defaults to 1 second. Minimum value is 1.
// periodSec: how often does the probe??defaults to 1 second. Minimum value is 1,Except for the first time?
// headers: headers[0] is HTTP Header, do not fill if you do not need to set
// on the other hand, only **first container** will be set livenessProbe
func (obj *Deployment) SetCMDLiveness(cmd []string, initDelaySec, timeoutSec, periodSec int32) *Deployment {
	setLiveness(&obj.dp.Spec.Template, cmdProbe(cmd, initDelaySec, timeoutSec, periodSec))
	return obj
}

// SetTCPLiveness set container liveness of tcp style
// host: default is ""
// port: required
// timeoutSec: http request timeout seconds,defaults to 1 second. Minimum value is 1.
// periodSec: how often does the probe??defaults to 1 second. Minimum value is 1,Except for the first time?
// headers: headers[0] is HTTP Header, do not fill if you do not need to set
// on the other hand, only **first container** will be set livenessProbe
func (obj *Deployment) SetTCPLiveness(host string, port int, initDelaySec, timeoutSec, periodSec int32) *Deployment {
	setLiveness(&obj.dp.Spec.Template, tcpProbe(host, port, initDelaySec, timeoutSec, periodSec))
	return obj
}

// SetHTTPReadness set container readness
// initDelaySec: how long time after the first start of the program the probe is executed for the first time.(sec)
// timeoutSec: http request timeout seconds,defaults to 1 second. Minimum value is 1.
// periodSec: how often does the probe??defaults to 1 second. Minimum value is 1,Except for the first time?
// on the other hand, only **first container** will be set livenessProbe
func (obj *Deployment) SetHTTPReadness(port int, path string, initDelaySec, timeoutSec, periodSec int32, headers ...map[string]string) *Deployment {
	setReadness(&obj.dp.Spec.Template, httpProbe(port, path, initDelaySec, timeoutSec, periodSec, headers...))
	return obj
}

// SetCMDReadness set container readness of cmd style
// cmd: execute readness probe as commond line
// timeoutSec: http request timeout seconds,defaults to 1 second. Minimum value is 1.
// periodSec: how often does the probe? defaults to 1 second. Minimum value is 1,Except for the first time?
// headers: headers[0] is HTTP Header, do not fill if you do not need to set
// on the other hand, only **first container** will be set livenessProbe
func (obj *Deployment) SetCMDReadness(cmd []string, initDelaySec, timeoutSec, periodSec int32) *Deployment {
	setReadness(&obj.dp.Spec.Template, cmdProbe(cmd, initDelaySec, timeoutSec, periodSec))
	return obj
}

// SetTCPReadness set container readness of tcp style
// host: default is ""
// port: required
// timeoutSec: http request timeout seconds,defaults to 1 second. Minimum value is 1.
// periodSec: how often does the probe? defaults to 1 second. Minimum value is 1,Except for the first time?
// headers: headers[0] is HTTP Header, do not fill if you do not need to set
// on the other hand, only **first container** will be set livenessProbe
func (obj *Deployment) SetTCPReadness(host string, port int, initDelaySec, timeoutSec, periodSec int32) *Deployment {
	setReadness(&obj.dp.Spec.Template, tcpProbe(host, port, initDelaySec, timeoutSec, periodSec))
	return obj
}

// SetPreStopExec set StatefulSet PreStop command
// PreStop is called immediately before a container is terminated.
// The container is terminated after the handler completes.
// The reason for termination is passed to the handler.
// Regardless of the outcome of the handler, the container is eventually terminated.
// Other management of the container blocks until the hook completes.
// More info: https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/#container-hooks
func (obj *Deployment) SetPreStopExec(command []string) *Deployment {
	setPreStopExec(&obj.dp.Spec.Template, command)
	return obj
}

// SetPostStartExec set PostStart shell command style
// PostStart is called immediately after a container is created. If the handler fails,
// the container is terminated and restarted according to its restart policy.
// Other management of the container blocks until the hook completes.
// More info: https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/#container-hooks
func (obj *Deployment) SetPostStartExec(command []string) *Deployment {
	setPostStartExec(&obj.dp.Spec.Template, command)
	return obj
}

// SetPreStopHTTP set preStop  http style
// PreStop is called immediately before a container is terminated.
// The container is terminated after the handler completes.
// The reason for termination is passed to the handler.
// Regardless of the outcome of the handler, the container is eventually terminated.
// Other management of the container blocks until the hook completes.
// More info: https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/#container-hooks
func (obj *Deployment) SetPreStopHTTP(scheme URIScheme, host string, port int, path string, headers ...map[string]string) *Deployment {
	setPreStopHTTP(&obj.dp.Spec.Template, scheme, host, port, path, headers...)
	return obj
}

// SetPostStartHTTP set  PostStart http style
// PostStart is called immediately after a container is created. If the handler fails,
// the container is terminated and restarted according to its restart policy.
// Other management of the container blocks until the hook completes.
// More info: https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/#container-hooks
func (obj *Deployment) SetPostStartHTTP(scheme URIScheme, host string, port int, path string, headers ...map[string]string) *Deployment {
	setPostStartHTTP(&obj.dp.Spec.Template, scheme, host, port, path, headers...)
	return obj
}

// SetMatchExpressions set Deployment match expressions
// the field is used to set complicated Label.
func (obj *Deployment) SetMatchExpressions(ents []LabelSelectorRequirement) *Deployment {
	if len(ents) <= 0 {
		return obj
	}
	requirements := make([]metav1.LabelSelectorRequirement, 0)
	for index := range ents {
		requirements = append(requirements, metav1.LabelSelectorRequirement{
			Key:      ents[index].Key,
			Operator: metav1.LabelSelectorOperator(ents[index].Operator),
			Values:   ents[index].Values,
		})
	}
	if obj.dp.Spec.Selector == nil {
		obj.dp.Spec.Selector = &metav1.LabelSelector{
			MatchExpressions: requirements,
		}
		return obj
	}
	obj.dp.Spec.Selector.MatchExpressions = requirements
	return obj
}

// SetDeployMaxTime set Deployment deploy max time,default 600s.
// If real deploy time more than this value,Deployment controller return err:ProgressDeadlineExceeded
// and Pod will Redeploy.
func (obj *Deployment) SetDeployMaxTime(sec int32) *Deployment {
	if sec < 0 {
		sec = 600
	}
	obj.dp.Spec.ProgressDeadlineSeconds = &sec
	return obj
}

// SetPodQos set pod  quality of service
// qosClass: is quality of service,the value only 'Guaranteed','Burstable' and 'BestEffort'
// autoSet: If your previous settings do not meet the requirements of PodQoS, we will automatically set
func (obj *Deployment) SetPodQos(qosClass string, autoSet ...bool) *Deployment {
	obj.SetAnnotations(setQosMap(obj.dp.Annotations, qosClass, autoSet...))
	return obj
}

// SetPodLabels set Pod labels
// when call SetLabels(),you can not use this function.
func (obj *Deployment) SetPodLabels(labels map[string]string) *Deployment {
	obj.SetSelector(labels)
	return obj
}

// SetImagePullSecrets set pod pull secret
func (obj *Deployment) SetImagePullSecrets(secretName string) *Deployment {
	setImagePullSecrets(&obj.dp.Spec.Template, secretName)
	return obj
}

// GetPodLabel get Pod labels
func (obj *Deployment) GetPodLabel() map[string]string {
	return obj.dp.Spec.Template.GetLabels()
}

// SetPodPriorityClass set Deployment Pod Priority
// priorityClassName is Kubernetes resource object PriorityClass name
// priorityClassName must already exists in kubernetes cluster
func (obj *Deployment) SetPodPriorityClass(priorityClassName string) *Deployment {
	obj.error(setPodPriorityClass(&obj.dp.Spec.Template, priorityClassName))
	return obj
}

// SetPVClaim set Deployment PersistentVolumeClaimVolumeSource
// params:
// volumeName: this is Custom field,you can define VolumeSource name,will be used of the container MountPath,
// claimName: this is PersistentVolumeClaim(PVC) name,the PVC and Deployment must on same namespace and exist.
func (obj *Deployment) SetPVClaim(volumeName, claimName string) *Deployment {
	obj.error(setPVClaim(&obj.dp.Spec.Template, volumeName, claimName))
	return obj
}

//SetPVCMounts mount PersistentVolumeClaim on container
// params:
// volumeName:the param is SetPVClaim() function volumeName,and when you call SetPVCMounts function you must call SetPVClaim function,and no order.
// on the other hand SetPVCMounts() function only mount first Container,and On the Container you can volumeMount many PersistentVolumeClaim.
// mountPath: runtime container dir eg:/var/lib/mysql
func (obj *Deployment) SetPVCMounts(volumeName, mountPath string) *Deployment {
	obj.error(setPVCMounts(&obj.dp.Spec.Template, volumeName, mountPath))
	return obj
}

func (obj *Deployment) error(err error) {
	if obj.err != nil {
		return
	}
	obj.err = err
}

// ImagePullPolicy  Deployment  pull image policy:Always,Never,IfNotPresent
func (obj *Deployment) ImagePullPolicy(pullPolicy PullPolicy) *Deployment {
	if len(obj.dp.Annotations) <= 0 {
		obj.dp.Annotations = make(map[string]string, 0)
	}
	obj.dp.Annotations[ImagePullPolicyKey] = string(pullPolicy)
	return obj
}

// SetContainer set Deployment container
// name:name is container name ,default ""
// image:image is image name ,must input image
// containerPort: image expose containerPort,must input containerPort
func (obj *Deployment) SetContainer(name, image string, containerPort int32) *Deployment {
	obj.error(setContainer(&obj.dp.Spec.Template, name, image, containerPort))
	return obj
}

// SetContainerOne set one container
func (obj *Deployment) SetContainerOne(container corev1.Container) *Deployment {
	if obj.dp.Spec.Template.Spec.Containers == nil {
		obj.dp.Spec.Template.Spec.Containers = []corev1.Container{container}
		return obj
	}
	obj.dp.Spec.Template.Spec.Containers = append(obj.dp.Spec.Template.Spec.Containers, container)
	return obj
}

// SetResourceLimit set container of deployment resource limit,eg:CPU and MEMORY
func (obj *Deployment) SetResourceLimit(limits map[ResourceName]string) *Deployment {
	obj.error(setResourceLimit(&obj.dp.Spec.Template, limits))
	return obj
}

// SetResourceRequst set container of deployment resource request,only CPU and MEMORY
func (obj *Deployment) SetResourceRequst(requests map[ResourceName]string) *Deployment {
	obj.error(setResourceRequests(&obj.dp.Spec.Template, requests))
	return obj
}

// SetEnvs set Pod Environmental variable
func (obj *Deployment) SetEnvs(envMap map[string]string) *Deployment {
	obj.error(setEnvs(&obj.dp.Spec.Template, envMap))
	return obj
}

// Release release Deployment on Kubernetes
func (obj *Deployment) Release() (*v1.Deployment, error) {
	dp, err := obj.Finish()
	if err != nil {
		return nil, err
	}
	client, err := GetKubeClient()
	if err != nil {
		return nil, err
	}
	return client.AppsV1().Deployments(dp.GetNamespace()).Create(dp)
}

// Apply  it will be updated when this resource object exists in K8s,
// it will be created when it does not exist.
func (obj *Deployment) Apply() (*v1.Deployment, error) {
	dp, err := obj.Finish()
	if err != nil {
		return nil, err
	}
	client, err := GetKubeClient()
	if err != nil {
		return nil, err
	}
	_, err = client.AppsV1().Deployments(dp.GetNamespace()).Get(dp.GetName(), metav1.GetOptions{})
	if err != nil {
		return client.AppsV1().Deployments(dp.GetNamespace()).Create(dp)
	}
	return client.AppsV1().Deployments(dp.GetNamespace()).Update(dp)
}

// DelNodeAffinity delete node affinitys
// keys is delete key list
func (obj *Deployment) DelNodeAffinity(keys []string) *Deployment {
	delNodeAffinity(&obj.dp.Spec.Template, keys)
	return obj
}

// SetRequiredORNodeAffinity set node affinity  for RequiredDuringSchedulingIgnoredDuringExecution style
// A list of keys, many key do OR operation.
func (obj *Deployment) SetRequiredORNodeAffinity(key string, value []string, operator NodeSelectorOperator) *Deployment {
	nsRequirement := corev1.NodeSelectorRequirement{
		Key:      key,
		Operator: operator.ToK8s(),
		Values:   value,
	}
	obj.error(setRequiredORNodeAffinity(&obj.dp.Spec.Template, nsRequirement))
	return obj
}

// SetRequiredAndNodeAffinity set node affinity  for RequiredDuringSchedulingIgnoredDuringExecution style
// A list of keys, many key do AND operation.
func (obj *Deployment) SetRequiredAndNodeAffinity(key string, value []string, operator NodeSelectorOperator) *Deployment {
	nsRequirement := corev1.NodeSelectorRequirement{
		Key:      key,
		Operator: operator.ToK8s(),
		Values:   value,
	}
	obj.error(setRequiredAndNodeAffinity(&obj.dp.Spec.Template, nsRequirement))
	return obj
}

// SetToleration set Taints Tolerations
// delayTimeSec  TolerationSeconds represents the period of time the toleration (which must be
// of effect NoExecute, otherwise this field is ignored) tolerates the taint. By default,
// it is not set, which means tolerate the taint forever (do not evict). Zero and
// negative values will be treated as 0 (evict immediately) by the system.
// +optional
// operator default is Equal
func (obj *Deployment) SetToleration(key, value string, operator TolerationOperator, effect TaintEffect, delayTimeSec ...int64) *Deployment {
	var tolerationSeconds int64
	toleration := corev1.Toleration{
		Key:      key,
		Value:    value,
		Operator: operator.ToK8s(),
		Effect:   effect.ToK8s(),
	}
	if len(delayTimeSec) > 0 && delayTimeSec[0] > 0 {
		tolerationSeconds = delayTimeSec[0]
	}
	if effect == TaintEffectNoExecute {
		toleration.TolerationSeconds = &tolerationSeconds
	}

	if operator == TolerationOpExists {
		toleration.Value = ""
	}

	setTolerations(&obj.dp.Spec.Template, toleration)
	return obj
}

// SetPreferredNodeAffinity set node affinity for PreferredDuringSchedulingIgnoredDuringExecution style
func (obj *Deployment) SetPreferredNodeAffinity(weight int32, key string, value []string, operator NodeSelectorOperator) *Deployment {
	nsRequirement := corev1.NodeSelectorRequirement{
		Key:      key,
		Operator: operator.ToK8s(),
		Values:   value,
	}
	obj.error(setPreferredNodeAffinity(&obj.dp.Spec.Template, nsRequirement, weight))
	return obj
}

// verify check service necessary value, input the default field and input related data.
func (obj *Deployment) verify() {
	if obj.err != nil {
		return
	}
	if !verifyString(obj.dp.GetName()) {
		obj.err = errors.New("Deployment name is not allowed to be empty")
		return
	}
	if len(obj.dp.Spec.Template.GetLabels()) < 1 {
		obj.err = errors.New("Deployment.Spec.Templata.Labels is not allowed to be empty")
		return
	}

	if err := containerRepeated(obj.dp.Spec.Template.Spec.Containers); err != nil {
		obj.err = fmt.Errorf("Deployment.Spec.Template.Spec.Containers err:%s", err.Error())
		return
	}
	if obj.dp.Spec.Selector == nil {
		obj.SetSelector(obj.GetPodLabel())
	}

	//check qos set,if err!=nil, check need auto set qos
	presentQos, err := qosCheck(obj.dp.Annotations[qosKey], obj.dp.Spec.Template.Spec)
	if err != nil {
		if obj.dp.Annotations[autoQosKey] == "true" {
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
	obj.dp.Kind = "Deployment"
	obj.dp.APIVersion = "apps/v1"
	if obj.dp.Annotations[ImagePullPolicyKey] == "" {
		for index := range obj.dp.Spec.Template.Spec.Containers {
			obj.dp.Spec.Template.Spec.Containers[index].ImagePullPolicy = corev1.PullIfNotPresent
		}
		return
	}
	policy := PullPolicy(obj.dp.Annotations[ImagePullPolicyKey]).ToK8s()
	for index := range obj.dp.Spec.Template.Spec.Containers {
		obj.dp.Spec.Template.Spec.Containers[index].ImagePullPolicy = policy
	}
	delete(obj.dp.Annotations, ImagePullPolicyKey)

}

// autoSetQos auto set Pod of Deployment QOS
func (obj *Deployment) autoSetQos(presentQos string) error {
	return autoSetQos(obj.dp.Annotations[qosKey], presentQos, &obj.dp.Spec.Template.Spec)
}
