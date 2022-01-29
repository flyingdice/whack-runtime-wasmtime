package internal

import (
	"github.com/bytecodealliance/wasmtime-go"
	"github.com/flyingdice/whack-runtime-wasmtime/internal/consts"
	"github.com/flyingdice/whack-sdk/pkg/sdk"
	"github.com/flyingdice/whack-sdk/pkg/sdk/runtime"
	"github.com/pkg/errors"
)

var _ sdk.RuntimeInstance = (*Instance)(nil)

// Instance wraps an active Wasmtime instance with a helpful
// API for interacting with it.
type Instance struct {
	env      *wasmtime.WasiConfig
	instance *wasmtime.Instance
	store    *wasmtime.Store
	wrn      sdk.WRN
}

// WRN returns the resource name for the instance.
func (i *Instance) WRN() sdk.WRN { return i.wrn }

// Stdout returns the stdout stream for the instance.
//
// This will be empty if InheritStdout was not set on WasiConfig.
func (i *Instance) Stdout() string { return "" }

// Stderr returns the stderr stream for the instance.z
//
// This will be empty if InheritStderr was not set on WasiConfigs.
func (i *Instance) Stderr() string { return "" }

// Call invokes an exported function by name with the given arguments.
func (i *Instance) Call(name string, args ...interface{}) (interface{}, error) {
	function := i.instance.GetFunc(i.store, name)
	if function == nil {
		return nil, errors.Errorf("exported function %s not found", name)
	}
	return function.Call(i.store, args...)
}

// Close will free resources used by the underlying wasmtime instance.
func (i *Instance) Close() error {
	i.instance = nil
	return nil
}

// Read memory of given length at a specific address location.
func (i *Instance) Read(mem sdk.Memory) ([]byte, error) {
	export := i.instance.GetExport(i.store, consts.ExportedMemoryName)
	if export == nil {
		return nil, errors.New("failed to get memory export")
	}
	memory := export.Memory()
	if memory == nil {
		return nil, errors.New("failed to get exported memory")
	}

	addr := mem.Address()
	length := mem.Length()

	data := memory.UnsafeData(i.store)
	if int(length) > len(data) {
		return nil, errors.Errorf("expected %d bytes; memory only %d bytes", length, len(data))
	}

	buffer := make([]byte, length)

	if read := copy(buffer, data[addr:addr+length]); int32(read) != length {
		return nil, errors.Errorf("expected to read %d; got %d", length, read)
	}

	return buffer, nil
}

// Write bytes to memory at a specific address location.
func (i *Instance) Write(addr int32, bytes []byte) (int32, error) {
	export := i.instance.GetExport(i.store, consts.ExportedMemoryName)
	if export == nil {
		return 0, errors.New("failed to get memory export")
	}
	memory := export.Memory()
	if memory == nil {
		return 0, errors.New("failed to get exported memory")
	}

	length := int32(len(bytes))
	data := memory.UnsafeData(i.store)

	if written := copy(data[addr:addr+length], bytes); int32(written) != length {
		return 0, errors.Errorf("expected to write %d; got %d", length, written)
	}

	return length, nil
}

func NewInstance(rt *Runtime, hostImports runtime.HostImports) (*Instance, error) {
	// Unique id for this instance so guest invoked host functions can
	// find the runtime instance they're executing in.
	id := sdk.RandomWRN()

	// Bind imports to runtime linker.
	err := bindHostImports(rt, id, hostImports)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create import object")
	}

	// Create Instance for module.
	inst, err := rt.linker.Instantiate(rt.store, rt.module)
	if err != nil {
		return nil, errors.Wrap(err, "fa	iled to create instance")
	}

	// Fetch and invoke wasi start ffi if one exists (commands).
	// Note: Start functions take no arguments and have no return value.
	start, err := wasiStartFunc(rt.store, inst)
	if err != nil {
		return nil, err
	}
	if start != nil {
		if _, err := start.Call(rt.store); err != nil {
			return nil, errors.Wrap(err, "failed to invoke wasi start function")
		}
	}

	// Fetch and invoke wasi initialize ffi if one exists (reactors).
	// Note: Init functions take no arguments and have no return value.
	init, err := wasiInitFunc(rt.store, inst)
	if err != nil {
		return nil, err
	}
	if init != nil {
		if _, err := init.Call(rt.store); err != nil {
			return nil, errors.Wrap(err, "failed to invoke wasi init function")
		}
	}

	return &Instance{
		wrn:      id,
		instance: inst,
		env:      rt.wasi,
	}, nil
}

// wasiStartFunc returns a NativeFunction start ffi for the Instance.
//
// If no NativeFunction is found, nil is returned without error.
//
// https://github.com/WebAssembly/WASI/blob/main/design/application-abi.md
func wasiStartFunc(store *wasmtime.Store, instance *wasmtime.Instance) (*wasmtime.Func, error) {
	return instance.GetFunc(store, consts.StartFunctionName), nil
}

// wasiInitFunc returns a NativeFunction initialize ffi for the Instance.
//
// If no NativeFunction is found, nil is returned without error.
//
// https://github.com/WebAssembly/WASI/blob/main/design/application-abi.md
func wasiInitFunc(store *wasmtime.Store, instance *wasmtime.Instance) (*wasmtime.Func, error) {
	return instance.GetFunc(store, consts.InitFunctionName), nil
}
