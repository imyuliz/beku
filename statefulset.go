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
	obj.verify()
	return obj.sts, obj.err
}

// JSONNew use json data create StatelfulSet
func (obj *StatefulSet) JSONNew(jsonbyts []byte) *StatefulSet {
	obj.error(json.Unmarshal(jsonbyts, obj.sts))
	return obj
}

// YAMLNew use yaml data create StatefulSet
func (obj *StatefulSet) YAMLNew(yamlbyts []byte) *StatefulSet {
	obj.error(yaml.Unmarshal(yamlbyts, obj.sts))
	return obj
}

// Replace replace StatefulSet by Kubernetes resource object
func (obj *StatefulSet) Replace(sts *v1.StatefulSet) *StatefulSet {
	if sts != nil {
		obj.sts = sts
	}
	return obj
}

// SetName set StatefulSet(sts) name
func (obj *StatefulSet) SetName(name string) *StatefulSet {
	obj.sts.SetName(name)
	return obj
}

// SetNamespace set StatefulSet(sts) namespace ,default namespace is 'default'
func (obj *StatefulSet) SetNamespace(namespace string) *StatefulSet {
	obj.sts.SetNamespace(namespace)
	return obj
}

// SetNamespaceAndName set StatefulSet namespace,set Pod namespace,set Deployment name.
func (obj *StatefulSet) SetNamespaceAndName(namespace, name string) *StatefulSet {
	obj.SetNamespace(namespace)
	obj.SetName(name)
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

// SetPodPriorityClass set StatefulSet Pod Priority
// priorityClassName is Kubernetes resource object PriorityClass name
// priorityClassName must already exists in kubernetes cluster
func (obj *StatefulSet) SetPodPriorityClass(priorityClassName string) *StatefulSet {
	obj.error(setPodPriorityClass(&obj.sts.Spec.Template, priorityClassName))
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
// mountPath: runtime container dir eg:/var/lib/mysql
func (obj *StatefulSet) SetPVCMounts(volumeName, mountPath string) *StatefulSet {
	obj.error(setPVCMounts(&obj.sts.Spec.Template, volumeName, mountPath))
	return obj
}

// SetPVCTemp set StatefulSet PersistentVolumeClaimTemplate
// can't call SetPVCMounts() function when you call the function,
// because SetPVCMounts() function has been called automatically,
// Don't worry
func (obj *StatefulSet) SetPVCTemp(pvcName, mountPath string, mode PersistentVolumeAccessMode, requests map[ResourceName]string) *StatefulSet {
	if !verifyString(pvcName) {
		obj.error(errors.New("SetPVCTemp failed,pvcName is not allowed to be empty"))
		return obj
	}
	if !verifyString(mountPath) {
		obj.error(errors.New("SetPVCTemp failed,mountPath is not allowed to be empty"))
		return obj
	}
	reqs, err := ResourceMapsToK8s(requests)
	if err != nil {
		obj.error(err)
		return obj
	}
	obj.SetPVCMounts(pvcName, mountPath)
	temp := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{Name: pvcName},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{mode.ToK8s()},
			Resources:   corev1.ResourceRequirements{Requests: reqs},
		},
	}
	if len(obj.sts.Spec.VolumeClaimTemplates) <= 0 {
		obj.sts.Spec.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{temp}
		return obj
	}
	obj.sts.Spec.VolumeClaimTemplates = append(obj.sts.Spec.VolumeClaimTemplates, temp)
	return obj

}

// SetPreStopExec set StatefulSet PreStop command
// PreStop is called immediately before a container is terminated.
// The container is terminated after the handler completes.
// The reason for termination is passed to the handler.
// Regardless of the outcome of the handler, the container is eventually terminated.
// Other management of the container blocks until the hook completes.
// More info: https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/#container-hooks
func (obj *StatefulSet) SetPreStopExec(command []string) *StatefulSet {
	setPreStopExec(&obj.sts.Spec.Template, command)
	return obj
}

// SetPostStartExec set PostStart shell command style
// PostStart is called immediately after a container is created. If the handler fails,
// the container is terminated and restarted according to its restart policy.
// Other management of the container blocks until the hook completes.
// More info: https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/#container-hooks
func (obj *StatefulSet) SetPostStartExec(command []string) *StatefulSet {
	setPostStartExec(&obj.sts.Spec.Template, command)
	return obj
}

// SetPreStopHTTP set preStop  http style
// PreStop is called immediately before a container is terminated.
// The container is terminated after the handler completes.
// The reason for termination is passed to the handler.
// Regardless of the outcome of the handler, the container is eventually terminated.
// Other management of the container blocks until the hook completes.
// More info: https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/#container-hooks
func (obj *StatefulSet) SetPreStopHTTP(scheme URIScheme, host string, port int, path string, headers ...map[string]string) *StatefulSet {
	setPreStopHTTP(&obj.sts.Spec.Template, scheme, host, port, path, headers...)
	return obj
}

// SetPostStartHTTP set  PostStart http style
// PostStart is called immediately after a container is created. If the handler fails,
// the container is terminated and restarted according to its restart policy.
// Other management of the container blocks until the hook completes.
// More info: https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/#container-hooks
func (obj *StatefulSet) SetPostStartHTTP(scheme URIScheme, host string, port int, path string, headers ...map[string]string) *StatefulSet {
	setPostStartHTTP(&obj.sts.Spec.Template, scheme, host, port, path, headers...)
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

// SetImagePullSecrets set pod pull secret
func (obj *StatefulSet) SetImagePullSecrets(secretName string) *StatefulSet {
	setImagePullSecrets(&obj.sts.Spec.Template, secretName)
	return obj
}

// SetNodeAffinity set node Affinity
//  corev1: "k8s.io/api/core/v1"
// func (obj *StatefulSet) SetNodeAffinity(nodeAffinity *corev1.NodeAffinity) *StatefulSet {
// 	obj.error(setNodeAffinity(&obj.sts.Spec.Template, nodeAffinity))
// 	return obj
// }

// ImagePullPolicy  StatefulSet  pull image policy:Always,Never,IfNotPresent
func (obj *StatefulSet) ImagePullPolicy(pullPolicy PullPolicy) *StatefulSet {
	if len(obj.sts.Annotations) <= 0 {
		obj.sts.Annotations = make(map[string]string, 0)
	}
	obj.sts.Annotations[ImagePullPolicyKey] = string(pullPolicy)
	return obj
}

// Release release StatefulSet on Kubernetes
func (obj *StatefulSet) Release() (*v1.StatefulSet, error) {
	sts, err := obj.Finish()
	if err != nil {
		return nil, err
	}
	client, err := GetKubeClient()
	if err != nil {
		return nil, err
	}
	return client.AppsV1().StatefulSets(sts.GetNamespace()).Create(sts)
}

// Apply  it will be updated when this resource object exists in K8s,
// it will be created when it does not exist.
func (obj *StatefulSet) Apply() (*v1.StatefulSet, error) {
	sts, err := obj.Finish()
	if err != nil {
		return nil, err
	}
	client, err := GetKubeClient()
	if err != nil {
		return nil, err
	}
	_, err = client.AppsV1().StatefulSets(sts.GetNamespace()).Get(sts.GetName(), metav1.GetOptions{})
	if err != nil {
		return client.AppsV1().StatefulSets(sts.GetNamespace()).Create(sts)
	}
	return client.AppsV1().StatefulSets(sts.GetNamespace()).Update(sts)
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
		obj.err = errors.New("StatefulSet.Spec.Selector.MatchLabels is not allowed to be empty")
		return
	}
	if err := containerRepeated(obj.sts.Spec.Template.Spec.Containers); err != nil {
		obj.err = fmt.Errorf("StatefulSet.Spec.Template.Spec.Containers err:%s", err.Error())
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
	obj.sts.Kind = "StatefulSet"
	obj.sts.APIVersion = "apps/v1"
	if obj.sts.Annotations[ImagePullPolicyKey] == "" {
		for index := range obj.sts.Spec.Template.Spec.Containers {
			obj.sts.Spec.Template.Spec.Containers[index].ImagePullPolicy = corev1.PullIfNotPresent
		}
		return
	}
	policy := PullPolicy(obj.sts.Annotations[ImagePullPolicyKey]).ToK8s()
	for index := range obj.sts.Spec.Template.Spec.Containers {
		obj.sts.Spec.Template.Spec.Containers[index].ImagePullPolicy = policy
	}
	delete(obj.sts.Annotations, ImagePullPolicyKey)
}

// autoSetQos auto set Pod of StatefulSet QOS
func (obj *StatefulSet) autoSetQos(presentQos string) error {
	return autoSetQos(obj.sts.Annotations[qosKey], presentQos, &obj.sts.Spec.Template.Spec)
}
