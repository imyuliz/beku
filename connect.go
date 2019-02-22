package beku

import (
	"errors"
	"fmt"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Client k8s client
type client struct {
	Host     string
	CAData   []byte
	CertData []byte
	KeyData  []byte
}

var defaultClient = new(client)

func getClientConfig() *client {
	return defaultClient
}

func setClientConfig(host string, ca, cert, key []byte) error {
	defaultClient.Host = host
	if len(ca) <= 1 && len(cert) <= 1 && len(key) <= 1 {
		return nil
	}
	defaultClient.CAData = ca
	defaultClient.CertData = cert
	defaultClient.KeyData = key
	return nil
}

// GetKubeClient get Kubernetes apiServer
func GetKubeClient(isInCluster ...bool) (*kubernetes.Clientset, error) {
	// Incluster  call apiserver
	if len(isInCluster) > 0 && isInCluster[0] {
		restConf, err := rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("get InClusterConfig err:%s", err.Error())
		}
		return kubernetes.NewForConfig(restConf)
	}
	config := getClientConfig()
	if config.Host == "" {
		return nil, errors.New("get kubernetes apiserver error,Because Host is empty,you can call function RegisterK8sClient() register")
	}
	if ViaTLS(config.CAData, config.CertData, config.KeyData) {
		return getTLSKubeClient(config.Host, config.CAData, config.CertData, config.KeyData)
	}
	return getKubeClient(config.Host)
}

// ViaTLS  verify Kubernetes apiServer cert
func ViaTLS(ca, cert, key []byte) bool {
	return len(ca) > 1 && len(cert) > 1 && len(key) > 1
}

func getTLSKubeClient(host string, ca, cert, key []byte) (*kubernetes.Clientset, error) {
	return kubernetes.NewForConfig(&rest.Config{
		Host: host,
		TLSClientConfig: rest.TLSClientConfig{
			CAData:   ca,
			CertData: cert,
			KeyData:  key,
		},
	})

}

func getKubeClient(host string) (*kubernetes.Clientset, error) {
	return kubernetes.NewForConfig(&rest.Config{
		Host: host,
	})
}

// RegisterK8sClient register k8s apiServer Client on Beku
// If the certificate is not required, ca,cert,key field is ""
func RegisterK8sClient(host, ca, cert, key string) error {
	if strings.TrimSpace(host) == "" {
		return errors.New("RegisterK8sClient failed,host is not allowed to be empty")
	}
	return setClientConfig(host, []byte(ca), []byte(cert), []byte(key))
}

// RegisterK8sClientBase64 register k8s apiServer Client on Beku
// use the function when ca,cert,key were base64 encode.
// the function will base64 decode ca,cert,key
// ca is certificate-authority-data
// cert is client-certificate-data
// key is  client-key-data
func RegisterK8sClientBase64(host, ca, cert, key string) error {
	if strings.TrimSpace(host) == "" {
		return errors.New("RegisterK8sClient failed,host is not allowed to be empty")
	}
	var (
		caByts, certByts, keyByts []byte
		err                       error
	)
	if ca != "" && cert != "" && key != "" {
		caByts, err = Base64Decode(ca)
		if err != nil {
			return err
		}
		certByts, err = Base64Decode(cert)
		if err != nil {
			return err
		}
		keyByts, err = Base64Decode(key)
		if err != nil {
			return err
		}
	}
	setClientConfig(host, caByts, certByts, keyByts)
	return nil
}
