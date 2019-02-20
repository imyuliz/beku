package beku

/*
import (
	"errors"

	"k8s.io/kubernetes/pkg/apis/scheduling"
)

// PriorityClass defines the mapping from a priority class name to the priority
// integer value. The value can be any valid integer and err
type PriorityClass struct {
	pc  *scheduling.PriorityClass
	err error
}

// NewPriorityClass create PriorityClass and Chain function call begin with this function.
func NewPriorityClass() *PriorityClass { return &PriorityClass{pc: &scheduling.PriorityClass{}} }

// Finish Chain function call end with this function
// return real PriorityClass(really service is kubernetes resource object PriorityClass and error
// In the function, it will check necessary parametersainput the default field
func (obj *PriorityClass) Finish() (pc *scheduling.PriorityClass, err error) {
	obj.verify()
	pc, err = obj.pc, obj.err
	return
}

// SetName set priorityClass name
func (obj *PriorityClass) SetName(name string) *PriorityClass {
	obj.pc.SetName(name)
	return obj
}

// SetValue set priorityClass priority value,The higher the value, the higher the priority.
// The value range of 0<=prioriry<1000000000
func (obj *PriorityClass) SetValue(prioriry int32) *PriorityClass {
	if prioriry < 0 && prioriry > 1000000000 {
		obj.error(errors.New("PriorityClass.Value must be in the range of 0 to 1000000000"))
		return obj
	}
	obj.pc.Value = prioriry
	return obj
}

// SetGlobalDefault set priorityClass global default value,
// the default priority for pods that do not have any priority class.
// Only one PriorityClass can be marked as `globalDefault`. However, if more than
// one PriorityClasses exists with their `globalDefault` field set to true,
// the smallest value of such global default PriorityClasses will be used as the default priority.
// +optional
func (obj *PriorityClass) SetGlobalDefault(global bool) *PriorityClass {
	obj.pc.GlobalDefault = global
	return obj
}

// SetNameAnddValue set PriorityClass name priority value
func (obj *PriorityClass) SetNameAnddValue(name string, prioriry int32) *PriorityClass {
	obj.SetName(name)
	obj.SetValue(prioriry)
	return obj
}

// SetDescription set priorityClass Description is an arbitrary string that usually provides guidelines on
// when this priority class should be used.
func (obj *PriorityClass) SetDescription(desc string) *PriorityClass {
	obj.pc.Description = desc
	return obj
}

func (obj *PriorityClass) error(err error) {
	if obj.err != nil {
		return
	}
	obj.err = err
}

func (obj *PriorityClass) verify() {
	if obj.err != nil {
		return
	}

	if !verifyString(obj.pc.GetName()) {
		obj.err = errors.New("pc.Name is not allowed to be empty")
		return
	}
	obj.pc.APIVersion = "scheduling.k8s.io/v1beta1"
	obj.pc.Kind = "PriorityClass"
}
*/
