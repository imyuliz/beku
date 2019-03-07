package beku

import (
	"encoding/json"

	"github.com/ghodss/yaml"
	"k8s.io/api/core/v1"
)

// Node include Kubernetes resource object node and error
type Node struct {
	node *v1.Node
	err  error
}

/*
notes:
It is not allowed to create a new Node. because
Serious errors may occur
*/

// JSONNewNode use json data create Node
func JSONNewNode(jsonbyts []byte) *Node {
	obj := &Node{node: &v1.Node{}}
	obj.error(json.Unmarshal(jsonbyts, obj.node))
	return obj
}

// YAMLNewNode use yaml data create Node
func YAMLNewNode(yamlbyts []byte) *Node {
	obj := &Node{node: &v1.Node{}}
	obj.error(yaml.Unmarshal(yamlbyts, obj.node))
	return obj
}

// ReadNewNode read new node
func ReadNewNode(coreNode *v1.Node) *Node { return &Node{node: coreNode} }

// Finish Chain function call end with this function
// return Kubernetes resource object Node and error.
// In the function, it will check necessary parametersainput the default field
func (obj *Node) Finish() (node *v1.Node, err error) {
	obj.verify()
	node, err = obj.node, obj.err
	return
}

// SetLabels set node Label
// If the key already exists on node.Labels, the value will be replaced with the new one.
func (obj *Node) SetLabels(labels map[string]string) *Node {
	if len(obj.node.Labels) <= 0 {
		obj.node.Labels = labels
		return obj
	}
	for key, value := range labels {
		obj.node.Labels[key] = value
	}
	return obj
}

// DelNodeLabels del node labels
// Skip if there is no such key
func (obj *Node) DelNodeLabels(labels map[string]string) *Node {
	if len(obj.node.GetLabels()) <= 0 {
		return obj
	}
	for k := range labels {
		delete(obj.node.Labels, k)
	}
	return obj
}

// SetAnnotations set Node annotations
func (obj *Node) SetAnnotations(annotations map[string]string) *Node {
	if len(obj.node.Annotations) <= 0 {
		obj.node.Annotations = annotations
		return obj
	}
	for key, value := range annotations {
		obj.node.Annotations[key] = value
	}
	return obj
}

// SetTaints set Taint
func (obj *Node) SetTaints(key, value string, effect TaintEffect) *Node {
	taint := v1.Taint{
		Key:    key,
		Value:  value,
		Effect: effect.ToK8s(),
	}
	if len(obj.node.Spec.Taints) <= 0 {
		obj.node.Spec.Taints = []v1.Taint{taint}
		return obj
	}
	obj.node.Spec.Taints = append(obj.node.Spec.Taints, taint)
	return obj
}

func (obj *Node) error(err error) {
	if obj.err != nil {
		return
	}
	obj.err = err
}

func (obj *Node) verify() {
	if obj.err != nil {
		return
	}
}
