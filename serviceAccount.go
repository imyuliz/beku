package beku

import (
	"errors"

	corev1 "k8s.io/api/core/v1"
)

// ServiceAccount include kubernetes resource object ServiceAccount(sa) and error
type ServiceAccount struct {
	sa  *corev1.ServiceAccount
	err error
}

// NewSa  create  ServiceAccount(sa) and chain function call begin with this function.
func NewSa() *ServiceAccount { return &ServiceAccount{sa: &corev1.ServiceAccount{}} }

// Finish Chain function call end with this function
// return Kubernetes resource object ServiceAccount and error.
// In the function, it will check necessary parameters、input the default field。
func (obj *ServiceAccount) Finish() (*corev1.ServiceAccount, error) {
	obj.verify()
	return obj.sa, obj.err
}

func (obj *ServiceAccount) verify() {
	if obj.sa.GetName() == "" {
		obj.error(errors.New("Set Name err,name is not allowed to be empty"))
		return
	}
	obj.sa.APIVersion = "v1"
	obj.sa.Kind = "ServiceAccount"
}
func (obj *ServiceAccount) error(err error) {
	if obj.err != nil {
		return
	}
	obj.err = err
}

// SetNamespceAndName set namespace and name, name is not allowed to be empty, name default is ""
func (obj *ServiceAccount) SetNamespceAndName(namespace, name string) *ServiceAccount {
	obj.sa.SetNamespace(namespace)
	obj.sa.SetName(name)
	return obj
}
