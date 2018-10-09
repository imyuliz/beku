package beku

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"

	"github.com/ghodss/yaml"
	"github.com/yulibaozi/beku/core"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Service Kubernetes svc resource
type Service struct {
	v1  *v1.Service
	err error
}

// NewSVC   create service
func NewSVC() *Service {
	return &Service{
		v1: &v1.Service{},
	}
}

// JSONNew json style create service
func (svc *Service) JSONNew(jsonbyts []byte) *Service {
	err := json.Unmarshal(jsonbyts, svc.v1)
	if err != nil {
		svc.err = err
	}
	return svc
}

// YAMLNew yaml style create service
func (svc *Service) YAMLNew(yamlbyts []byte) *Service {
	err := yaml.Unmarshal(yamlbyts, svc.v1)
	if err != nil {
		svc.err = err
	}
	return svc
}

// Finish check err  and return service
func (svc *Service) Finish() (v1 *v1.Service, err error) {
	svc.verify()
	v1, err = svc.v1, svc.err
	return
}

// SetName setname
func (svc *Service) SetName(name string) *Service {
	svc.v1.SetName(name)
	return svc
}

// SetNameSpace set namespace
func (svc *Service) SetNameSpace(namespace string) *Service {
	svc.v1.SetNamespace(namespace)
	return svc
}

// SetNameAndNameSpace set name and namespace
func (svc *Service) SetNameAndNameSpace(name, namespace string) *Service {
	svc.v1.SetName(name)
	svc.v1.SetNamespace(namespace)
	return svc
}

// SetLabels set labels
func (svc *Service) SetLabels(labels map[string]string) *Service {
	svc.v1.SetLabels(labels)
	return svc
}

// SetSelector set seletor
func (svc *Service) SetSelector(selector map[string]string) *Service {
	svc.v1.Spec.Selector = selector
	return svc
}

// SetServiceType set service type
func (svc *Service) SetServiceType(st core.ServiceType) *Service {
	svc.v1.Spec.Type = st.ToK8s()
	return svc
}

// SetAnnotations set annotations
func (svc *Service) SetAnnotations(annotations map[string]string) *Service {
	svc.v1.SetAnnotations(annotations)
	return svc
}

// SetPorts set service ports
func (svc *Service) SetPorts(ports []core.ServicePort) *Service {
	svcPorts := make([]v1.ServicePort, 0)
	for _, data := range ports {
		svcPorts = append(svcPorts, v1.ServicePort{
			Name:       data.Name,
			Protocol:   data.Protocol.ToK8s(),
			Port:       data.Port,
			TargetPort: FromInt(data.TargetPort),
			NodePort:   data.NodePort,
		})
	}
	svc.v1.Spec.Ports = svcPorts
	return svc
}

// SetSessionAffinity set session affinity
func (svc *Service) SetSessionAffinity(affinity core.ServiceAffinity) *Service {
	svc.v1.Spec.SessionAffinity = affinity.ToK8s()
	return svc
}

// Verify 验证数据的可用性
func (svc *Service) verify() {
	if svc.err != nil {
		return
	}
	if !verifyString(svc.v1.Kind) {
		svc.v1.Kind = "Service"
	}
	if !verifyString(svc.v1.GetName()) {
		svc.err = errors.New("service name is allow empty")
		return
	}
	if !verifyString(svc.v1.APIVersion) {
		svc.v1.APIVersion = "v1"
	}
	portLen := len(svc.v1.Spec.Ports)
	if verifyMap(svc.v1.Spec.Selector) {
		if portLen < 1 {
			svc.err = errors.New("service ports not allow empty when spec.selector exists")
			return
		}
	}
	if portLen > 1 {
		nameMaps := make(map[string]bool, 0)
		for _, data := range svc.v1.Spec.Ports {
			if !verifyString(data.Name) {
				svc.err = errors.New("spec.port[x].name not allow empty when len(ports)>1")
				return
			}
			nameMaps[data.Name] = true
		}
		if len(nameMaps) != portLen {
			svc.err = errors.New("spec.ports name not allow repetition when len(ports)>1 ")
			return
		}
	}
}

func verifyString(str string) bool          { return !(str == "" || len(str) <= 0) }
func verifyMap(maps map[string]string) bool { return len(maps) > 0 }

// FromInt creates an IntOrString object with an int32 value. It is
// your responsibility not to call this method with a value greater
// than int32.
// TODO: convert to (val int32)
func FromInt(val int) intstr.IntOrString {
	if val > math.MaxInt32 || val < math.MinInt32 {
	}
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
