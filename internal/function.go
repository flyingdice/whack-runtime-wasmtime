package internal

import (
	"github.com/bytecodealliance/wasmtime-go"
	"github.com/flyingdice/whack-sdk/pkg/sdk"
)

// GuestFunc represents a function that can be used with wasmtime.
type GuestFunc func(*wasmtime.Caller, []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap)

// HostFunc translates between GuestFunc and sdk.Function.
type HostFunc func(...wasmtime.Val) (interface{}, error)

var int32Type = wasmtime.NewValType(wasmtime.KindI32)

// function represents an SDK function represented as a wasmtime guest function.
type function struct {
	name   string
	args   []*wasmtime.ValType
	retval []*wasmtime.ValType
	fn     GuestFunc
}

// Name returns the name of the function.
func (f *function) Name() string { return f.name }

// Bind creates a new function in the given linker.
func (f *function) Bind(namespace string, linker *wasmtime.Linker) error {
	fnType := wasmtime.NewFuncType(f.args, f.retval)
	return linker.FuncNew(namespace, f.name, fnType, f.fn)
}

// TODO need to return a function that can be consumed by the linker.

// NewFunction creates a function for the given sdk.Function.
func NewFunction(runtime *Runtime, instanceWrn sdk.WRN, fn sdk.Function) *function {
	// Translate golang args into int32.
	args := make([]*wasmtime.ValType, fn.NumIn())
	for i := 0; i < fn.NumIn(); i++ {
		args[i] = int32Type
	}

	// Translate golang return value into variable number (0 or 1) of int32 return values.
	retval := make([]*wasmtime.ValType, fn.NumOut())
	for i := 0; i < fn.NumOut(); i++ {
		retval[i] = int32Type
	}

	// create fn type
	// create wrapper

	return &function{
		name:   fn.Name(),
		args:   args,
		retval: retval,
		fn:     guestFunc(hostFunc(runtime, instanceWrn, fn)),
	}
}

// hostFunc decorates the given sdk.Function so it can be called by wasmtime.
//
// This is responsible for calling the actual golang host function.
func hostFunc(runtime *Runtime, instanceWrn sdk.WRN, fn sdk.Function) HostFunc {
	return func(args ...wasmtime.Val) (interface{}, error) {
		hostArgs := hostFuncArgs(args...)
		return fn.Func()(runtime, instanceWrn, hostArgs...)
	}
}

// hostFuncArgs translates array of wasmtime arg types to golang types.
func hostFuncArgs(args ...wasmtime.Val) []interface{} {
	retval := make([]interface{}, len(args))
	for i, arg := range args {
		retval[i] = arg.I32()
	}
	return retval
}

// guestFunc decorates the given HostFunc, so it can be called by sdk.Function.
// This is responsible for translating a WASM function invocation into golang and back.
func guestFunc(fn HostFunc) GuestFunc {
	return func(c *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
		result, err := fn(args...)
		if err != nil {
			return nil, wasmtime.NewTrap(err.Error())
		}

		retval := make([]wasmtime.Val, 0)
		if result != nil {
			retval = append(retval, wasmtime.ValI32(result.(int32)))
		}

		return retval, nil
	}
}
