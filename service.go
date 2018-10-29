package beku

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ghodss/yaml"
	"k8s.io/api/core/v1"
)

// Service include Kubernetes resource object Service and error
type Service struct {
	svc *v1.Service
	err error
}

// NewSvc create service(svc) and chain function call begin with this function.
func NewSvc() *Service { return &Service{svc: &v1.Service{}} }

// Finish Chain function call end with this function
// return real service(really service is kubernetes resource object Service and error
// In the function, it will check necessary parametersainput the default field
func (obj *Service) Finish() (svc *v1.Service, err error) {
	obj.verify()
	svc, err = obj.svc, obj.err
	return
}

// JSONNew use json data create service(svc)
func (obj *Service) JSONNew(jsonbyts []byte) *Service {
	obj.error(json.Unmarshal(jsonbyts, obj.svc))
	return obj
}

// YAMLNew use yaml data create service(svc)
func (obj *Service) YAMLNew(yamlbyts []byte) *Service {
	obj.error(yaml.Unmarshal(yamlbyts, obj.svc))
	return obj
}

// SetName set service(svc) name
func (obj *Service) SetName(name string) *Service {
	obj.svc.SetName(name)
	return obj
}

// SetNamespace set service(svc) namespace
func (obj *Service) SetNamespace(namespace string) *Service {
	obj.svc.SetNamespace(namespace)
	return obj
}

// SetNamespaceAndName set service(svc) namespace and name
// namespace default value is 'default'
func (obj *Service) SetNamespaceAndName(namespace, name string) *Service {
	obj.svc.SetName(name)
	obj.svc.SetNamespace(namespace)
	return obj
}

// SetLabels set service(svc) labels
func (obj *Service) SetLabels(labels map[string]string) *Service {
	obj.svc.SetLabels(labels)
	return obj
}

// SetSelector set service(svc) seletor
// The Pod that matches the selector will be selected
// the function Required call when you create service(svc)
func (obj *Service) SetSelector(selector map[string]string) *Service {
	obj.svc.Spec.Selector = selector
	return obj
}

// SetServiceType set service(svc) type,you can choose NodePort,ClusterIP
// many info please redirect to ServiceType
func (obj *Service) SetServiceType(sty ServiceType) *Service {
	obj.svc.Spec.Type = sty.ToK8s()
	return obj
}

// SetAnnotations set service(svc) annotations
func (obj *Service) SetAnnotations(annotations map[string]string) *Service {
	obj.svc.SetAnnotations(annotations)
	return obj
}

// SetPorts set service(svc) ports
func (obj *Service) SetPorts(ports []ServicePort) *Service {
	objPorts := make([]v1.ServicePort, 0)
	for _, data := range ports {
		objPorts = append(objPorts, v1.ServicePort{
			Name:       data.Name,
			Protocol:   data.Protocol.ToK8s(),
			Port:       data.Port,
			TargetPort: FromInt(data.TargetPort),
			NodePort:   data.NodePort,
		})
	}
	obj.svc.Spec.Ports = objPorts
	return obj
}

// SetPort set service(svc) Port. port params required input on ServicePort
// default TargetPort same as Port
// NodePort is random number when not input or NodePort <= 0
// Protocol default value 'TCP'
func (obj *Service) SetPort(sp ServicePort) *Service {
	sPort := v1.ServicePort{
		Name:       sp.Name,
		Protocol:   sp.Protocol.ToK8s(),
		Port:       sp.Port,
		TargetPort: FromInt(sp.TargetPort),
		NodePort:   sp.NodePort,
	}
	if len(obj.svc.Spec.Ports) > 0 {
		obj.svc.Spec.Ports = append(obj.svc.Spec.Ports, sPort)
		return obj
	}
	obj.svc.Spec.Ports = []v1.ServicePort{sPort}
	return obj
}

// SetSessionAffinity set service(svc) session affinity
func (obj *Service) SetSessionAffinity(affinity ServiceAffinity) *Service {
	obj.svc.Spec.SessionAffinity = affinity.ToK8s()
	return obj
}

func (obj *Service) error(err error) {
	if obj.err != nil {
		return
	}
	obj.err = err
}

// verify check service necessary value, input the default field and input related data.
func (obj *Service) verify() {
	if obj.err != nil {
		return
	}
	if !verifyString(obj.svc.GetName()) {
		obj.err = errors.New("svc.Name is not allowed to be empty")
		return
	}
	portLen := len(obj.svc.Spec.Ports)
	if verifyMap(obj.svc.Spec.Selector) {
		if portLen < 1 {
			obj.err = errors.New("svc.Spec.Ports is not allowed be empty when svc.Spec.Selector exist")
			return
		}
	}
	if portLen > 1 {
		nameMaps := make(map[string]bool, 0)
		for index, data := range obj.svc.Spec.Ports {
			if !verifyString(data.Name) {
				obj.err = fmt.Errorf("svc.Spec.port[%d].name is not allowed be empty when len(svc.Spec.Ports) > 1", index)
				return
			}
			nameMaps[data.Name] = true
		}
		if len(nameMaps) != portLen {
			obj.err = errors.New("len(svc.Spec.Ports.Name) != len(svc.Spec.Ports) is not allowed when len(svc.Sepc.Ports) > 1")
			return
		}
	}
	obj.svc.Kind = "Service"
	obj.svc.APIVersion = "v1"
}
