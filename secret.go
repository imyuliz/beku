package beku

import (
	"errors"

	"github.com/yulibaozi/beku/core"
	"k8s.io/api/core/v1"
)

// Secret 加密配置
type Secret struct {
	v1  *v1.Secret
	err error
}

// NewSecret create secret
func (secret *Secret) NewSecret() *Secret {
	return &Secret{
		v1: &v1.Secret{},
	}
}

// SetNamespace set secrete namespace ,default namespace is default
func (secret *Secret) SetNamespace(namespace string) *Secret {
	secret.v1.SetNamespace(namespace)
	return secret
}

// SetName set name
func (secret *Secret) SetName(name string) *Secret {
	secret.v1.SetName(name)
	return secret
}

// SetNameSpaceAndName set secret namespace and name
func (secret *Secret) SetNameSpaceAndName(namespace, name string) *Secret {
	secret.v1.SetNamespace(namespace)
	secret.v1.SetName(name)
	return secret
}

// SetLabels set secret labels
func (secret *Secret) SetLabels(labels map[string]string) *Secret {
	secret.v1.SetLabels(labels)
	return secret
}

// SetDataString set secret data
func (secret *Secret) SetDataString(datas map[string]string) *Secret {
	secret.v1.StringData = datas
	return secret
}

// SetDataBytes set secret data for byte
func (secret *Secret) SetDataBytes(bytes map[string][]byte) *Secret {
	secret.v1.Data = bytes
	return secret
}

// SetType set secret type,have Opaque and kubernetes.io/service-account-token two kind
// Opaque user-defined data
// kubernetes.io/service-account-token is used to kubernetes apiserver,because apiserver need to auth
func (secret *Secret) SetType(secType core.SecretType) *Secret {
	secret.v1.Type = secType.ToK8s()
	return secret
}

func (secret *Secret) verify() {
	if !verifyString(secret.v1.Name) {
		secret.err = errors.New("secret name not allow empty")
		return
	}
	secret.v1.APIVersion = "v1"
	secret.v1.Kind = "Secret"
}

// Finish the final step,will return kubernetes resource object secret and error
func (secret *Secret) Finish() (*v1.Secret, error) {
	secret.verify()
	if secret.err != nil {
		return nil, secret.err
	}
	return secret.v1, nil
}
