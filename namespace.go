package beku

import (
	"errors"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// Release release Namespace on Kubernetes
func (obj *Namespace) Release() (*v1.Namespace, error) {
	ns, err := obj.Finish()
	if err != nil {
		return nil, err
	}
	client, err := GetKubeClient()
	if err != nil {
		return nil, err
	}
	return client.CoreV1().Namespaces().Create(ns)
}

// Apply  it will be updated when this resource object exists in K8s,
// it will be created when it does not exist.
func (obj *Namespace) Apply() (*v1.Namespace, error) {
	ns, err := obj.Finish()
	if err != nil {
		return nil, err
	}
	client, err := GetKubeClient()
	if err != nil {
		return nil, err
	}
	_, err = client.CoreV1().Namespaces().Get(ns.GetName(), metav1.GetOptions{})
	if err != nil {
		return client.CoreV1().Namespaces().Create(ns)
	}
	return client.CoreV1().Namespaces().Update(ns)
}

func (obj *Namespace) verify() {
	if obj.ns.GetName() == "" {
		obj.err = errors.New("Namespace.Name is not allowed to be empty")
		return
	}
	obj.ns.APIVersion = "v1"
	obj.ns.Kind = "Namespace"
}
