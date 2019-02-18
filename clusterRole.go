package beku

import (
	"encoding/json"
	"errors"

	"github.com/ghodss/yaml"
	"k8s.io/api/rbac/v1beta1"
)

// ClusterRole include kubernetes resource object ClusterRole and error
type ClusterRole struct {
	role *v1beta1.ClusterRole
	err  error
}

// NewClusterRole  create  ClusterRole and chain function call begin with this function.
func NewClusterRole() *ClusterRole { return &ClusterRole{role: &v1beta1.ClusterRole{}} }

// Finish Chain function call end with this function
// return Kubernetes resource object ClusterRole and error.
// In the function, it will check necessary parameters、input the default field。
func (obj *ClusterRole) Finish() (*v1beta1.ClusterRole, error) {
	obj.verify()
	return obj.role, obj.err
}

// JSONNew use json data create ClusterRole
func (obj *ClusterRole) JSONNew(jsonbyts []byte) *ClusterRole {
	obj.error(json.Unmarshal(jsonbyts, obj.role))
	return obj
}

// YAMLNew use yaml data create ClusterRole
func (obj *ClusterRole) YAMLNew(yamlbyts []byte) *ClusterRole {
	obj.error(yaml.Unmarshal(yamlbyts, obj.role))
	return obj
}

func (obj *ClusterRole) verify() {
	if obj.role.GetName() == "" {
		obj.error(errors.New("Set Name err,name is not allowed to be empty"))
		return
	}
	if len(obj.role.Rules) <= 0 {
		obj.error(errors.New("Set ClusterRole err,rules is not allowed to be empty"))
		return
	}
	obj.role.APIVersion = "rbac.authorization.k8s.io/v1beta1"
	obj.role.Kind = "ClusterRole"
}
func (obj *ClusterRole) error(err error) {
	if obj.err != nil {
		return
	}
	obj.err = err
}

// SetName set ClusterRole name
func (obj *ClusterRole) SetName(name string) *ClusterRole {
	obj.role.SetName(name)
	return obj
}

// SetRole set cluster role
// verbs is func method. such as "get", "watch", "list","create", "delete" ..., you can set "*" if you want to use all the func method
// apiGroups is resource apiGroup. such as "v1", "apps/v1", "rbac.authorization.k8s.io"... , you can set "*" if you want to use all the resource apiGroup
// resources is resources object. such as "daemonsets", "deployments","replicasets", you can set "*" if you want to use all the resource object.
func (obj *ClusterRole) SetRole(verbs, apiGroups, resources []string) *ClusterRole {
	rule := v1beta1.PolicyRule{
		Verbs:     verbs,
		APIGroups: apiGroups,
		Resources: resources,
	}
	if obj.role.Rules != nil {
		obj.role.Rules = append(obj.role.Rules, rule)
		return obj
	}
	obj.role.Rules = []v1beta1.PolicyRule{rule}
	return obj
}
