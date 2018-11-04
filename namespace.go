package beku

import (
	"errors"

	"k8s.io/api/core/v1"
)

// Namespace include Kubernets resource object Namespace and err
type Namespace struct {
	ns  *v1.Namespace
	err error
}

// NewNs create Namespace and Chain function call begin with this function.
func NewNs() *Namespace { return &Namespace{ns: &v1.Namespace{}} }

// Finish Chain function call end with this function
// return Kubernetes resource object Namespace and error.
// In the function, it will check necessary parametersăinput the default field
func (obj *Namespace) Finish() (*v1.Namespace, error) {
	obj.verify()
	return obj.ns, obj.err
}

// SetName set namespace name
func (obj *Namespace) SetName(name string) *Namespace {
	obj.ns.SetName(name)
	return obj
}

func (obj *Namespace) verify() {
	if obj.ns.GetName() == "" {
		obj.err = errors.New("Namespace.Name is not allowed to be empty")
		return
	}
	obj.ns.APIVersion = "v1"
	obj.ns.Kind = "Namespace"
}
