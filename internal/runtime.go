package internal

import (
	"github.com/bytecodealliance/wasmtime-go"
	"github.com/flyingdice/whack-sdk/pkg/sdk"
	"github.com/flyingdice/whack-sdk/pkg/sdk/runtime"
	"github.com/pkg/errors"
)

var _ sdk.Runtime = (*Runtime)(nil)

// TODO (ahawker) Cache compiled modules?
// TODO (ahawker) Reactor support? Do we call start/init during creation or lazy?

// Runtime encapsulates all Wasmtime state necessary to run WASM code.
type Runtime struct {
	cfg    runtime.Config
	engine *wasmtime.Engine
	wasi   *wasmtime.WasiConfig
	linker *wasmtime.Linker
	module *wasmtime.Module
	pool   sdk.InstancePool
	store  *wasmtime.Store
}

// New creates a new Instance for the runtime.
func (r *Runtime) New() (sdk.Instance, error) {
	instance, err := NewInstance(r, r.cfg.Host().Imports())
	if err != nil {
		return nil, err
	}
	return r.pool.Set(instance)
}

// Get returns the pool instance for the given wrn.
func (r *Runtime) Get(wrn sdk.WRN) (sdk.Instance, error) {
	return r.pool.Get(wrn)
}

// NewRuntime creates a new runtime.
func NewRuntime(mod sdk.Module, cfg runtime.Config) (*Runtime, error) {
	// Create global state.
	engine := wasmtime.NewEngine()

	// Compile raw module bytes into code (WAT).
	compiled, err := wasmtime.NewModule(engine, mod.Bytes())
	if err != nil {
		return nil, errors.Wrap(err, "failed to compile module")
	}

	// Create linker with WASI support.
	linker, err := newLinker(engine)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create linker")
	}

	// Create WASI config.
	wasi, err := wasiConfig(mod.Wrn(), cfg.Wasi())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create wasi config")
	}

	// Create store with WASI support.
	store, err := newStore(engine, wasi)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create store")
	}

	// Create runtime instance pool.
	pool, err := runtime.NewPool()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create instance pool")
	}

	return &Runtime{
		engine: engine,
		wasi:   wasi,
		linker: linker,
		module: compiled,
		cfg:    cfg,
		pool:   pool,
		store:  store,
	}, nil
}

// newLinker creates a wasmtime linker with WASI support.
func newLinker(engine *wasmtime.Engine) (*wasmtime.Linker, error) {
	linker := wasmtime.NewLinker(engine)
	if err := linker.DefineWasi(); err != nil {
		return nil, errors.Wrap(err, "failed to create wasmtime wasi")
	}
	return linker, nil
}

func newStore(engine *wasmtime.Engine, config *wasmtime.WasiConfig) (*wasmtime.Store, error) {
	store := wasmtime.NewStore(engine)
	store.SetWasi(config)
	return store, nil
}

func wasiConfig(wrn sdk.WRN, config runtime.WasiConfig) (*wasmtime.WasiConfig, error) {
	cfg := wasmtime.NewWasiConfig()

	// Arguments
	cfg.SetArgv(config.Arguments())

	// Environment Variables
	keys, vals := envToSlices(config.EnvVars())
	cfg.SetEnv(keys, vals)

	//// Preopen directories
	//for name, path := range config.Directories() {
	//	// TODO (ahawker) correct order?
	//	if err := cfg.PreopenDir(path, name); err != nil {
	//		return nil, errors.Wrapf(err, "failed to add preopen directory %s: %s", name, path)
	//	}
	//}

	// Stdio
	if config.CaptureStdout() {
		cfg.InheritStdout()
	}
	if config.CaptureStderr() {
		cfg.InheritStderr()
	}

	// Workdir
	if err := cfg.PreopenDir(config.Workdir(), "."); err != nil {
		return nil, errors.Wrapf(err, "failed to preopen workdir %s", config.Workdir())
	}

	return cfg, nil
}

func envToSlices(m map[string]string) ([]string, []string) {
	keys := make([]string, 0, len(m))
	vals := make([]string, 0, len(m))

	for k, v := range m {
		keys = append(keys, k)
		vals = append(vals, v)
	}

	return keys, vals
}
