package beku

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
	"k8s.io/api/core/v1"
	apiresource "k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
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
