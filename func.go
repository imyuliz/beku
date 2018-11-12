package beku

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	apiresource "k8s.io/apimachinery/pkg/api/resource"

	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// ToYAML struct translate into yaml
func ToYAML(o interface{}) (byts []byte, err error) {
	byts, err = yaml.Marshal(o)
	return
}

// ToJSON struct translate into json
func ToJSON(v interface{}) (byts []byte, err error) {
	byts, err = json.Marshal(v)
	return
}

// JSONToYAML json data translate into yaml
func JSONToYAML(jbyts []byte) (ybyts []byte, err error) {
	ybyts, err = yaml.JSONToYAML(jbyts)
	return
}

// YAMLToJSON yaml data translate into json
func YAMLToJSON(ybyts []byte) (jbyts []byte, err error) {
	jbyts, err = yaml.YAMLToJSON(ybyts)
	return
}

// DeploymentToSvc  Use the Deployment to generate the associated SVC
// autoRelease[0] if true,beku will auto Release Service On Kubernetes,default can't Release Service On Kubernetes
func DeploymentToSvc(dp *appsv1.Deployment, sty ServiceType, autoRelease ...bool) (*v1.Service, error) {
	var ports []ServicePort
	for _, data := range dp.Spec.Template.Spec.Containers {
		ports = append(ports, ServicePort{
			Name:     data.Name,
			Protocol: Protocol(data.Ports[0].Protocol),
			Port:     data.Ports[0].ContainerPort,
		})
	}
	svc := NewSvc().SetNamespaceAndName(dp.GetNamespace(), dp.GetName()).SetSelector(dp.Spec.Template.GetLabels()).SetPorts(ports).SetServiceType(sty)
	if len(autoRelease) > 0 && autoRelease[0] == true {
		return svc.Release()
	}
	return svc.Finish()
}

// StatefulSetToSvc  Use the StatefulSet to generate the associated SVC
// autoRelease[0] if true,beku will auto Release Service On Kubernetes,default can't Release Service On Kubernetes
func StatefulSetToSvc(sts *appsv1.StatefulSet, sty ServiceType, isHeadless bool, autoRelease ...bool) (*v1.Service, error) {
	var ports []ServicePort
	for _, data := range sts.Spec.Template.Spec.Containers {
		ports = append(ports, ServicePort{
			Name:     data.Name,
			Protocol: Protocol(data.Ports[0].Protocol),
			Port:     data.Ports[0].ContainerPort,
		})
	}
	svc := NewSvc().SetNamespaceAndName(sts.GetNamespace(), sts.GetName()).SetSelector(sts.Spec.Template.GetLabels()).SetPorts(ports)
	if isHeadless {
		svc = svc.Headless()
	} else {
		svc = svc.SetServiceType(sty)
	}
	if len(autoRelease) > 0 && autoRelease[0] == true {
		return svc.Release()
	}
	return svc.Finish()
}

// DaemonSetToSvc  Use the Set to generate the associated SVC
// autoRelease[0] if true,beku will auto Release Service On Kubernetes,default can't Release Service On Kubernetes
func DaemonSetToSvc(ds *appsv1.DaemonSet, sty ServiceType, autoRelease ...bool) (*v1.Service, error) {
	var ports []ServicePort
	for _, data := range ds.Spec.Template.Spec.Containers {
		ports = append(ports, ServicePort{
			Name:     data.Name,
			Protocol: Protocol(data.Ports[0].Protocol),
			Port:     data.Ports[0].ContainerPort,
		})
	}
	svc := NewSvc().SetNamespaceAndName(ds.GetNamespace(), ds.GetName()).SetSelector(ds.Spec.Template.GetLabels()).SetPorts(ports).SetServiceType(sty)
	if len(autoRelease) > 0 && autoRelease[0] == true {
		return svc.Release()
	}
	return svc.Finish()
}

// Base64Encode base64 encode
func Base64Encode(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

// Base64Decode base64 decode
func Base64Decode(src string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(src)
}

// ResourceMapsToK8s to K8s resourceList
func ResourceMapsToK8s(maps map[ResourceName]string) (v1.ResourceList, error) {
	data := make(v1.ResourceList, 0)
	for k, v := range maps {
		q, err := apiresource.ParseQuantity(v)
		if err != nil {
			return nil, err
		}
		reName := k.ToK8s()
		if reName == "" {
			return nil, errors.New("resource name not allow")
		}
		data[reName] = q
	}
	if len(data) < 1 {
		return nil, errors.New("source cann't allow empty")
	}
	return data, nil
}

// FromInt creates an IntOrString object with an int32 value. It is
// your responsibility not to call this method with a value greater
// than int32.
// TODO: convert to (val int32)
func FromInt(val int) intstr.IntOrString {
	return intstr.IntOrString{Type: intstr.Int, IntVal: int32(val)}
}

// FromString creates an IntOrString object with a string value.
func FromString(val string) intstr.IntOrString {
	return intstr.IntOrString{Type: intstr.String, StrVal: val}
}

// Parse the given string and try to convert it to an integer before
// setting it as a string value.
func Parse(val string) intstr.IntOrString {
	i, err := strconv.Atoi(val)
	if err != nil {
		return FromString(val)
	}
	return FromInt(i)
}

// httpProbe  container health check  readness and liveness  http probe
func httpProbe(port int, path string, initDelaySec, timeoutSec, periodSec int32, headers ...map[string]string) *v1.Probe {
	if initDelaySec <= 0 {
		initDelaySec = 30
	}
	return &v1.Probe{
		Handler:             v1.Handler{HTTPGet: &v1.HTTPGetAction{Path: path, Port: FromInt(port), HTTPHeaders: mapsToHeaders(headers)}},
		InitialDelaySeconds: initDelaySec,
		TimeoutSeconds:      timeoutSec,
		PeriodSeconds:       periodSec,
	}
}

func mapsToHeaders(headers []map[string]string) []v1.HTTPHeader {
	if len(headers) <= 0 {
		return nil
	}
	return mapToHeaders(headers[0])
}

func mapToHeaders(header map[string]string) []v1.HTTPHeader {
	var headers []v1.HTTPHeader
	for key := range header {
		headers = append(headers, v1.HTTPHeader{Name: key, Value: header[key]})
	}
	return headers
}

// cmdProbe container health check readness and liveness cmd probe
func cmdProbe(cmd []string, initDelaySec, timeoutSec, periodSec int32) *v1.Probe {
	if initDelaySec <= 0 {
		initDelaySec = 30
	}
	return &v1.Probe{
		Handler:             v1.Handler{Exec: &v1.ExecAction{Command: cmd}},
		InitialDelaySeconds: initDelaySec,
		TimeoutSeconds:      timeoutSec,
		PeriodSeconds:       periodSec,
	}
}

// tcpProbe container health check readness and liveness tcp probe
func tcpProbe(host string, port int, initDelaySec, timeoutSec, periodSec int32) *v1.Probe {
	if initDelaySec <= 0 {
		initDelaySec = 30
	}
	return &v1.Probe{
		Handler:             v1.Handler{TCPSocket: &v1.TCPSocketAction{Port: FromInt(port), Host: host}},
		InitialDelaySeconds: initDelaySec,
		TimeoutSeconds:      timeoutSec,
		PeriodSeconds:       periodSec,
	}
}
func verifyString(str string) bool          { return !(str == "" || len(str) <= 0) }
func verifyMap(maps map[string]string) bool { return len(maps) > 0 }

func mapToEnvs(envMap map[string]string) ([]v1.EnvVar, error) {
	if len(envMap) <= 0 {
		return nil, errors.New("SetEnvs error, envMap is not allowed to be empty")
	}
	var envs []v1.EnvVar
	for k, v := range envMap {
		k, v = strings.TrimSpace(k), strings.TrimSpace(v)
		if k == "" || v == "" {
			return nil, fmt.Errorf("SetEnvs error, key or value is not allowed to be empty,data(%s:%s)", k, v)
		}
		envs = append(envs, v1.EnvVar{Name: k, Value: v})
	}
	if len(envs) <= 0 {
		return nil, fmt.Errorf("SetEnvs error, envs is not allowed to be empty")
	}
	return envs, nil
}

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

// GetKubeClient get Kubernetes apiServer
func GetKubeClient() (*kubernetes.Clientset, error) {
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
	if ca != "" && cert != "" && key != "" {
		defaultClient.Host = host
		defaultClient.CAData = []byte(ca)
		defaultClient.CertData = []byte(cert)
		defaultClient.KeyData = []byte(key)
		return nil
	}
	defaultClient.Host = host
	return nil
}
