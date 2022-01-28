package internal

import (
	"github.com/flyingdice/whack-runtime-wasmtime/internal/consts"
	"github.com/flyingdice/whack-sdk/pkg/sdk"
	"github.com/flyingdice/whack-sdk/pkg/sdk/runtime"
	"github.com/pkg/errors"
)

// bindHostImports binds functions/globals into the linker to be accessible by WASM code.
func bindHostImports(
	runtime *Runtime,
	instanceWrn sdk.WRN,
	imports runtime.HostImports,
) error {
	for name, fn := range importFunctions(runtime, instanceWrn, imports) {
		if err := fn.Bind(consts.Namespace, runtime.linker); err != nil {
			return errors.Wrapf(err, "failed to bind function %s", name)
		}
	}

	for name, gbl := range importGlobals(runtime, instanceWrn, imports) {
		if err := gbl.Bind(consts.Namespace, runtime.linker); err != nil {
			return errors.Wrapf(err, "failed to bind global %s", name)
		}
	}

	return nil
}

// importFunctions builds a map of host functions that can be accessed by WASM code.
func importFunctions(runtime *Runtime, instanceWrn sdk.WRN, imports runtime.HostImports) map[string]*function {
	functions := make(map[string]*function)
	for _, fn := range imports.Functions() {
		functions[fn.Name()] = NewFunction(runtime, instanceWrn, fn)
	}
	return functions
}

// importGlobals builds a map of globals that can be accessed by WASM code.
func importGlobals(runtime *Runtime, instanceWrn sdk.WRN, imports runtime.HostImports) map[string]*global {
	globals := make(map[string]*global)
	return globals
}
