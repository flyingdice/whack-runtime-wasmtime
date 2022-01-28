package sdk

import (
	"github.com/pkg/errors"
)

// Func represents a generic golang host func that can be called by WASM through runtime translation.
type Func func(Runtime, WRN, ...interface{}) (interface{}, error)

// TODO
type HostFunc func(Instance, ...int32) (interface{}, error)

// Function represents a generic function that can be imported into a module.
type Function interface {
	Name() string
	NumIn() int
	NumOut() int
	Func() Func
}

type function struct {
	name   string
	numIn  int
	numOut int
	fn     Func
}

func (f *function) Name() string { return f.name }
func (f *function) NumIn() int   { return f.numIn }
func (f *function) NumOut() int  { return f.numOut }
func (f *function) Func() Func   { return f.fn }

// NewFunction creates a new function from the given name/golang Func.
func NewFunction(name string, numIn, numOut int, fn HostFunc) *function {
	return &function{
		name:   name,
		numIn:  numIn,
		numOut: numOut,
		fn: func(runtime Runtime, instanceWrn WRN, args ...interface{}) (interface{}, error) {
			// Fetch runtime instance being used for this call.
			instance, err := runtime.Get(instanceWrn)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get instance for wrn %s", instanceWrn)
			}
			// Convert golang args into primitives supported by wasm.
			wargs, err := wasmArgs(args...)
			if err != nil {
				return nil, errors.Wrap(err, "failed to convert args to wasm args")
			}
			// TODO
			result, err := fn(instance, wargs...)
			instance.SetResult(result, err)
			return nil, nil
		},
	}
}

// wasmArgs takes a variable number of interface{} values and converts them to int32's.
// WASM only supports simple primitive types, so we convert virtually everything to int32
// between the boundary.
// TODO (ahawker) support floats
func wasmArgs(args ...interface{}) ([]int32, error) {
	retval := make([]int32, len(args))
	for i, arg := range args {
		asserted, ok := arg.(int32)
		if !ok {
			return nil, errors.Errorf("arg %d type assert to int32 failed", i)
		}
		retval[i] = asserted
	}
	return retval, nil
}
