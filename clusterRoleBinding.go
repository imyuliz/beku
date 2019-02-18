package beku

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/ghodss/yaml"
	"k8s.io/api/rbac/v1beta1"
)

// ClusterRoleBinding include kubernetes resource object ClusterRoleBinding and error
type ClusterRoleBinding struct {
	crb *v1beta1.ClusterRoleBinding
	err error
}

// NewClusterRoleBinding create  NewClusterRoleBinding and chain function call begin with this function.
func NewClusterRoleBinding() *ClusterRoleBinding {
	return &ClusterRoleBinding{crb: &v1beta1.ClusterRoleBinding{}}
}

// Finish Chain function call end with this function
// return Kubernetes resource object ClusterRoleBinding and error.
// In the function, it will check necessary parameters、input the default field。
func (obj *ClusterRoleBinding) Finish() (*v1beta1.ClusterRoleBinding, error) {
	obj.verify()
	return obj.crb, obj.err
}

// JSONNew use json data create ClusterRoleBinding
func (obj *ClusterRoleBinding) JSONNew(jsonbyts []byte) *ClusterRoleBinding {
	obj.error(json.Unmarshal(jsonbyts, obj.crb))
	return obj
}

// YAMLNew use yaml data create ClusterRoleBinding
func (obj *ClusterRoleBinding) YAMLNew(yamlbyts []byte) *ClusterRoleBinding {
	obj.error(yaml.Unmarshal(yamlbyts, obj.crb))
	return obj
}

// SetName set ClusterRoleBinding name
func (obj *ClusterRoleBinding) SetName(name string) *ClusterRoleBinding {
	obj.crb.SetName(name)
	return obj
}

// SubKind subject kind
type SubKind string

// subject kinds
const (
	User  SubKind = "User"
	Group SubKind = "Group"
	SA    SubKind = "ServiceAccount"
)

var (
	kindMaps = map[SubKind]string{
		User:  "rbac.authorization.k8s.io",
		Group: "rbac.authorization.k8s.io",
		SA:    "namespace",
	}
)

// Subject set  ClusterRoleBinding subject
// kind only support "User", "Group", "ServiceAccount"
// namespace  it is Required when kind is "ServiceAccount" default is "". it is Optional when kind is "User" or "Group"
func (obj *ClusterRoleBinding) Subject(name string, kind SubKind, namespace string) *ClusterRoleBinding {
	if kindMaps[kind] == "" {
		obj.error(fmt.Errorf("Set subject err. kind:%v is not supported, only support User/Group/ServiceAccount ", kind))
		return obj
	}
	subject := v1beta1.Subject{
		Name: name,
		Kind: string(kind),
	}
	if kindMaps[kind] == "namespace" {
		subject.Namespace = namespace
	} else {
		subject.APIGroup = kindMaps[kind]
	}
	if len(obj.crb.Subjects) <= 0 {
		obj.crb.Subjects = []v1beta1.Subject{subject}
		return obj
	}
	obj.crb.Subjects = append(obj.crb.Subjects, subject)
	return obj
}

// SetRoleRef set ClusterRoleBinding RoleRef
func (obj *ClusterRoleBinding) SetRoleRef(name string) *ClusterRoleBinding {
	if emptyString(name) {
		obj.error(errors.New("set SetRoleRef err. name is not allow to be empty"))
		return obj
	}
	obj.crb.RoleRef = v1beta1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "ClusterRole",
		Name:     name,
	}
	return obj
}

var tmpRoleRef = v1beta1.RoleRef{}

func (obj *ClusterRoleBinding) verify() {

	if obj.crb.GetName() == "" {
		obj.error(errors.New("Set Name err,name is not allowed to be empty"))
		return
	}

	if len(obj.crb.Subjects) <= 0 {
		obj.error(errors.New("Set ClusterRoleBinding err,subejects is not allowed to be empty"))
		return
	}

	if reflect.DeepEqual(tmpRoleRef, obj.crb.RoleRef) {
		obj.error(errors.New("Set ClusterRoleBinding err,RoleRef is not allowed to be empty"))
		return
	}

	obj.crb.APIVersion = "rbac.authorization.k8s.io/v1beta1"
	obj.crb.Kind = "ClusterRoleBinding"
}

func (obj *ClusterRoleBinding) error(err error) {
	if obj.err != nil {
		return
	}
	obj.err = err
}
