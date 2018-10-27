package beku

import (
	"errors"
	"strings"

	"k8s.io/api/core/v1"
)

// ServicePort service ports
type ServicePort struct {
	Name       string   `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	Protocol   Protocol `json:"protocol,omitempty" protobuf:"bytes,2,opt,name=protocol,casttype=Protocol"`
	Port       int32    `json:"port" protobuf:"varint,3,opt,name=port"`
	TargetPort int      `json:"targetPort,omitempty" protobuf:"bytes,4,opt,name=targetPort"`
	NodePort   int32    `json:"nodePort,omitempty" protobuf:"varint,5,opt,name=nodePort"`
}

var (
	resourceLimit   = make(map[ResourceName]string, 0)
	resourceRequest = make(map[ResourceName]string, 0)
)

func defaultLimit() map[ResourceName]string {
	return resourceLimit
}

func defaultRequest() map[ResourceName]string {
	return resourceRequest

}

// RegisterResourceLimit register you need default resource limit, resource only include CPU and MEMORY
func RegisterResourceLimit(limits map[ResourceName]string) error {
	if len(limits) == 2 && verifyString(limits[ResourceCPU]) && verifyString(limits[ResourceMemory]) {
		resourceLimit[ResourceCPU] = limits[ResourceCPU]
		resourceLimit[ResourceMemory] = limits[ResourceMemory]
		return nil
	}
	return errors.New("resource limit must include cpu and memory and only include cpu and memory")

}

// RegisterResourceRequest register you need default resource limit, resource only include CPU and MEMORY
func RegisterResourceRequest(request map[ResourceName]string) error {
	if len(request) == 2 && verifyString(request[ResourceCPU]) && verifyString(request[ResourceMemory]) {
		resourceRequest[ResourceCPU] = request[ResourceCPU]
		resourceRequest[ResourceMemory] = request[ResourceMemory]
		return nil
	}
	return errors.New("resource request must include cpu and memory and only include cpu and memory")
}

// ServiceType service type
type ServiceType string

// ServiceType
const (
	// ServiceTypeClusterIP means a service will only be accessible inside the
	// cluster, via the cluster IP.
	ServiceTypeClusterIP ServiceType = "ClusterIP"
	// ServiceTypeNodePort means a service will be exposed on one port of
	// every node, in addition to 'ClusterIP' type.
	ServiceTypeNodePort ServiceType = "NodePort"
	// ServiceTypeLoadBalancer means a service will be exposed via an
	// external load balancer (if the cloud provider supports it), in addition
	// to 'NodePort' type.
	ServiceTypeLoadBalancer ServiceType = "LoadBalancer"
	// ServiceTypeExternalName means a service consists of only a reference to
	// an external name that kubedns or equivalent will return as a CNAME
	// record, with no exposing or proxying of any pods involved.
	ServiceTypeExternalName ServiceType = "ExternalName"
)

// ToK8s  translate into Kubernetes ServiceType
func (sty ServiceType) ToK8s() v1.ServiceType {
	switch t := string(sty); t {
	case "ClusterIP":
		return v1.ServiceTypeClusterIP
	case "NodePort":
		return v1.ServiceTypeNodePort
	// case "ExternalName":
	// 	return v1.ServiceTypeExternalName
	case "LoadBalancer":
		return v1.ServiceTypeLoadBalancer
	default:
		return v1.ServiceTypeClusterIP
	}
}

// Protocol defines network protocols supported for things like container ports.
type Protocol string

const (
	// ProtocolTCP is the TCP protocol.
	ProtocolTCP Protocol = "TCP"
	// ProtocolUDP is the UDP protocol.
	ProtocolUDP Protocol = "UDP"
)

// ToK8s translate into Kubernetes Protocol
func (pro Protocol) ToK8s() v1.Protocol {
	switch p := string(pro); strings.ToUpper(p) {
	case "TCP":
		return v1.ProtocolTCP
	case "UDP":
		return v1.ProtocolUDP
	}
	return v1.ProtocolTCP
}

// PersistentVolumeMode describes how a volume is intended to be consumed, either Block or Filesystem.
type PersistentVolumeMode string

const (
	// PersistentVolumeBlock means the volume will not be formatted with a filesystem and will remain a raw block device.
	PersistentVolumeBlock PersistentVolumeMode = "Block"
	// PersistentVolumeFilesystem means the volume will be or is formatted with a filesystem.
	PersistentVolumeFilesystem PersistentVolumeMode = "Filesystem"
)

// ToK8s   PersistentVolumeMode translate into   k8s PersistentVolumeMode
func (vMode PersistentVolumeMode) ToK8s() *v1.PersistentVolumeMode {
	switch v := string(vMode); v {
	case "Block":
		b := v1.PersistentVolumeBlock
		return &b
	case "Filesystem":
		f := v1.PersistentVolumeFilesystem
		return &f
	}
	return nil
}

// LabelSelector : A label selector is a label query over a set of resources. The result of matchLabels and
// matchExpressions are ANDed. An empty label selector matches all objects. A null
// label selector matches no objects.
type LabelSelector struct {
	// matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels
	// map is equivalent to an element of matchExpressions, whose key field is "key", the
	// operator is "In", and the values array contains only "value". The requirements are ANDed.
	// +optional
	MatchLabels map[string]string `json:"matchLabels,omitempty" protobuf:"bytes,1,rep,name=matchLabels"`
	// matchExpressions is a list of label selector requirements. The requirements are ANDed.
	// +optional
	MatchExpressions []LabelSelectorRequirement `json:"matchExpressions,omitempty" protobuf:"bytes,2,rep,name=matchExpressions"`
}

//LabelSelectorRequirement : A label selector requirement is a selector that contains values, a key, and an operator that
// relates the key and values.
type LabelSelectorRequirement struct {
	// key is the label key that the selector applies to.
	// +patchMergeKey=key
	// +patchStrategy=merge
	Key string `json:"key" patchStrategy:"merge" patchMergeKey:"key" protobuf:"bytes,1,opt,name=key"`
	// operator represents a key's relationship to a set of values.
	// Valid operators are In, NotIn, Exists and DoesNotExist.
	Operator LabelSelectorOperator `json:"operator" protobuf:"bytes,2,opt,name=operator,casttype=LabelSelectorOperator"`
	// values is an array of string values. If the operator is In or NotIn,
	// the values array must be non-empty. If the operator is Exists or DoesNotExist,
	// the values array must be empty. This array is replaced during a strategic
	// merge patch.
	// +optional
	Values []string `json:"values,omitempty" protobuf:"bytes,3,rep,name=values"`
}

// LabelSelectorOperator :A label selector operator is the set of operators that can be used in a selector requirement.
type LabelSelectorOperator string

// LabelSelectorOperator params
const (
	LabelSelectorOpIn           LabelSelectorOperator = "In"
	LabelSelectorOpNotIn        LabelSelectorOperator = "NotIn"
	LabelSelectorOpExists       LabelSelectorOperator = "Exists"
	LabelSelectorOpDoesNotExist LabelSelectorOperator = "DoesNotExist"
)

// CephFSPersistentVolumeSource  ceph volume setting
type CephFSPersistentVolumeSource struct {
	// Required: Monitors is a collection of Ceph monitors
	// More info: https://releases.k8s.io/HEAD/examples/volumes/cephfs/README.md#how-to-use-it
	Monitors []string `json:"monitors" protobuf:"bytes,1,rep,name=monitors"`
	// Optional: Used as the mounted root, rather than the full Ceph tree, default is /
	// +optional
	Path string `json:"path,omitempty" protobuf:"bytes,2,opt,name=path"`
	// Optional: User is the rados user name, default is admin
	// More info: https://releases.k8s.io/HEAD/examples/volumes/cephfs/README.md#how-to-use-it
	// +optional
	User string `json:"user,omitempty" protobuf:"bytes,3,opt,name=user"`
	// Optional: SecretFile is the path to key ring for User, default is /etc/ceph/user.secret
	// More info: https://releases.k8s.io/HEAD/examples/volumes/cephfs/README.md#how-to-use-it
	// +optional
	SecretFile string `json:"secretFile,omitempty" protobuf:"bytes,4,opt,name=secretFile"`
	// Optional: SecretRef is reference to the authentication secret for User, default is empty.
	// More info: https://releases.k8s.io/HEAD/examples/volumes/cephfs/README.md#how-to-use-it
	// +optional
	SecretRef *SecretReference `json:"secretRef,omitempty" protobuf:"bytes,5,opt,name=secretRef"`
	// Optional: Defaults to false (read/write). ReadOnly here will force
	// the ReadOnly setting in VolumeMounts.
	// More info: https://releases.k8s.io/HEAD/examples/volumes/cephfs/README.md#how-to-use-it
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,6,opt,name=readOnly"`
}

// SecretReference represents a Secret Reference. It has enough information to retrieve secret
// in any namespace
type SecretReference struct {
	// Name is unique within a namespace to reference a secret resource.
	// +optional
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	// Namespace defines the space within which the secret name must be unique.
	// +optional
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,2,opt,name=namespace"`
}

// NFSVolumeSource : Represents  an NFS mount that lasts the lifetime of a pod.
// NFS volumes do not support ownership management or SELinux relabeling.
type NFSVolumeSource struct {
	// Server is the hostname or IP address of the NFS server.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs
	Server string `json:"server" protobuf:"bytes,1,opt,name=server"`

	// Path that is exported by the NFS server.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs
	Path string `json:"path" protobuf:"bytes,2,opt,name=path"`

	// ReadOnly here will force
	// the NFS export to be mounted with read-only permissions.
	// Defaults to false.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,3,opt,name=readOnly"`
}

// ResourceName is the name identifying various resources in a ResourceList.
type ResourceName string

// Resource names must be not more than 63 characters, consisting of upper- or lower-case alphanumeric characters,
// with the -, _, and . characters allowed anywhere, except the first or last character.
// The default convention, matching that for annotations, is to use lower-case names, with dashes, rather than
// camel case, separating compound words.
// Fully-qualified resource typenames are constructed from a DNS-style subdomain, followed by a slash `/` and a name.
const (
	// CPU, in cores. (500m = .5 cores)
	ResourceCPU ResourceName = "cpu"
	// Memory, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	ResourceMemory ResourceName = "memory"
	// Volume size, in bytes (e,g. 5Gi = 5GiB = 5 * 1024 * 1024 * 1024)
	ResourceStorage ResourceName = "storage"
	// Local ephemeral storage, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	// The resource name for ResourceEphemeralStorage is alpha and it can change across releases.
	ResourceEphemeralStorage ResourceName = "ephemeral-storage"
	// NVIDIA GPU, in devices. Alpha, might change: although fractional and allowing values >1, only one whole device per node is assigned.
	ResourceNvidiaGPU ResourceName = "alpha.kubernetes.io/nvidia-gpu"
)

// ToK8s translate into k8s ResourceName
func (r ResourceName) ToK8s() v1.ResourceName {
	switch re := string(r); re {
	case "cpu":
		return v1.ResourceCPU
	case "storage":
		return v1.ResourceStorage
	case "memory":
		return v1.ResourceMemory
	case "ephemeral-storage":
		return v1.ResourceEphemeralStorage
	case "alpha.kubernetes.io/nvidia-gpu":
		return v1.ResourceNvidiaGPU
	}
	return ""
}

// PersistentVolumeAccessMode volume access mode read,write
type PersistentVolumeAccessMode string

// VolumeAccessMode params
const (
	// can be mounted read/write mode to exactly 1 host
	ReadWriteOnce PersistentVolumeAccessMode = "ReadWriteOnce"
	// can be mounted in read-only mode to many hosts
	ReadOnlyMany PersistentVolumeAccessMode = "ReadOnlyMany"
	// can be mounted in read/write mode to many hosts
	ReadWriteMany PersistentVolumeAccessMode = "ReadWriteMany"
	RWO           PersistentVolumeAccessMode = "RWO"
	ROX           PersistentVolumeAccessMode = "ROX"
	RWX           PersistentVolumeAccessMode = "RWX"
)

// ToK8s translate into k8s accessMode
func (pvm PersistentVolumeAccessMode) ToK8s() v1.PersistentVolumeAccessMode {
	switch m := string(pvm); m {
	case "ReadWriteOnce", "RWO":
		return v1.ReadWriteOnce
	case "ReadOnlyMany", "ROX":
		return v1.ReadOnlyMany
	case "ReadWriteMany", "RWX":
		return v1.ReadWriteMany
	}
	return ""
}

// RBDPersistentVolumeSource Represents a Rados Block Device mount that lasts the lifetime of a pod.
// RBD volumes support ownership management and SELinux relabeling.
type RBDPersistentVolumeSource struct {
	// A collection of Ceph monitors.
	// More info: https://releases.k8s.io/HEAD/examples/volumes/rbd/README.md#how-to-use-it
	CephMonitors []string `json:"monitors" protobuf:"bytes,1,rep,name=monitors"`
	// The rados image name.
	// More info: https://releases.k8s.io/HEAD/examples/volumes/rbd/README.md#how-to-use-it
	RBDImage string `json:"image" protobuf:"bytes,2,opt,name=image"`
	// Filesystem type of the volume that you want to mount.
	// Tip: Ensure that the filesystem type is supported by the host operating system.
	// Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#rbd
	// TODO: how do we prevent errors in the filesystem from compromising the machine
	// +optional
	FSType string `json:"fsType,omitempty" protobuf:"bytes,3,opt,name=fsType"`
	// The rados pool name.
	// Default is rbd.
	// More info: https://releases.k8s.io/HEAD/examples/volumes/rbd/README.md#how-to-use-it
	// +optional
	RBDPool string `json:"pool,omitempty" protobuf:"bytes,4,opt,name=pool"`
	// The rados user name.
	// Default is admin.
	// More info: https://releases.k8s.io/HEAD/examples/volumes/rbd/README.md#how-to-use-it
	// +optional
	RadosUser string `json:"user,omitempty" protobuf:"bytes,5,opt,name=user"`
	// Keyring is the path to key ring for RBDUser.
	// Default is /etc/ceph/keyring.
	// More info: https://releases.k8s.io/HEAD/examples/volumes/rbd/README.md#how-to-use-it
	// +optional
	Keyring string `json:"keyring,omitempty" protobuf:"bytes,6,opt,name=keyring"`
	// SecretRef is name of the authentication secret for RBDUser. If provided
	// overrides keyring.
	// Default is nil.
	// More info: https://releases.k8s.io/HEAD/examples/volumes/rbd/README.md#how-to-use-it
	// +optional
	SecretRef *SecretReference `json:"secretRef,omitempty" protobuf:"bytes,7,opt,name=secretRef"`
	// ReadOnly here will force the ReadOnly setting in VolumeMounts.
	// Defaults to false.
	// More info: https://releases.k8s.io/HEAD/examples/volumes/rbd/README.md#how-to-use-it
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,8,opt,name=readOnly"`
}

// ServiceAffinity Type set affinity
type ServiceAffinity string

// affinity params
const (
	// ServiceAffinityClientIP is the Client IP based.
	ServiceAffinityClientIP ServiceAffinity = "ClientIP"

	// ServiceAffinityNone - no session affinity.
	ServiceAffinityNone ServiceAffinity = "None"
)

// ToK8s translate into k8s aserviceAffinity
func (sa ServiceAffinity) ToK8s() v1.ServiceAffinity {
	switch p := string(sa); strings.ToUpper(p) {
	case "CLIENTIP":
		return v1.ServiceAffinityClientIP
	case "NONE":
		return v1.ServiceAffinityNone
	}
	return v1.ServiceAffinityNone
}

// SecretType 'Opaque' or 'kubernetes.io/service-account-token'
type SecretType string

const (
	// SecretTypeOpaque is the default. Arbitrary user-defined data
	SecretTypeOpaque SecretType = "Opaque"

	// SecretTypeServiceAccountToken contains a token that identifies a service account to the API
	//
	// Required fields:
	// - Secret.Annotations["kubernetes.io/service-account.name"] - the name of the ServiceAccount the token identifies
	// - Secret.Annotations["kubernetes.io/service-account.uid"] - the UID of the ServiceAccount the token identifies
	// - Secret.Data["token"] - a token that identifies the service account to the API
	SecretTypeServiceAccountToken SecretType = "kubernetes.io/service-account-token"
)

// ToK8s translate into Kubernets SecretType
func (ty SecretType) ToK8s() v1.SecretType {
	switch p := string(ty); p {
	case "kubernetes.io/service-account-token":
		return v1.SecretTypeServiceAccountToken
	default:
		return v1.SecretTypeOpaque
	}
}

// PullPolicy describes a policy for if/when to pull a container image
type PullPolicy string

const (
	// PullAlways means that kubelet always attempts to pull the latest image. Container will fail If the pull fails.
	PullAlways PullPolicy = "Always"
	// PullNever means that kubelet never pulls an image, but only uses a local image. Container will fail if the image isn't present
	PullNever PullPolicy = "Never"
	// PullIfNotPresent means that kubelet pulls if the image isn't present on disk. Container will fail if the image isn't present and the pull fails.
	PullIfNotPresent PullPolicy = "IfNotPresent"
)

// ToK8s image pull policy
func (pp PullPolicy) ToK8s() v1.PullPolicy {
	switch p := string(pp); p {
	case "Always":
		return v1.PullAlways
	case "Never":
		return v1.PullNever
	default:
		return v1.PullIfNotPresent
	}
}

// PodQOSClass defines the supported qos classes of Pods.
type PodQOSClass string

const (
	// PodQOSGuaranteed is the Guaranteed qos class.
	PodQOSGuaranteed PodQOSClass = "Guaranteed"
	// PodQOSBurstable is the Burstable qos class.
	PodQOSBurstable PodQOSClass = "Burstable"
	// PodQOSBestEffort is the BestEffort qos class.
	PodQOSBestEffort PodQOSClass = "BestEffort"
)

// ToK8s set pod qos
func (qos PodQOSClass) ToK8s() v1.PodQOSClass {
	switch p := string(qos); p {
	case "Guaranteed":
		return v1.PodQOSGuaranteed
	case "Burstable":
		return v1.PodQOSBurstable
	case "BestEffort":
		return v1.PodQOSBestEffort
	}
	return ""
}

// qosNotices set Qos information
var (
	qosNotices = map[string]string{
		"Guaranteed": "Every Container in the Pod must have a memory limit and a memory request, and they must be the same, Every Container in the Pod must have a CPU limit and a CPU request, and they must be the same,more information: https://kubernetes.io/docs/tasks/configure-pod-container/quality-service-pod",
		"Burstable":  "The Pod does not meet the criteria for QoS class Guaranteed and at least one Container in the Pod has a memory or CPU request, more information: https://kubernetes.io/docs/tasks/configure-pod-container/quality-service-pod",
		"BestEffort": "The Containers in the Pod must not have any memory or CPU limits or requests, more information: https://kubernetes.io/docs/tasks/configure-pod-container/quality-service-pod",
	}
)
