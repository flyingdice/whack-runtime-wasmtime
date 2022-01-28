package runtime

import (
	"github.com/flyingdice/whack-sdk/pkg/sdk"
)

type wasiConfig struct {
	args    []string
	stdout  bool
	stderr  bool
	dirs    map[string]string
	env     map[string]string
	workdir string
}

func (w *wasiConfig) Arguments() []string            { return w.args }
func (w *wasiConfig) CaptureStdout() bool            { return w.stdout }
func (w *wasiConfig) CaptureStderr() bool            { return w.stderr }
func (w *wasiConfig) Directories() map[string]string { return w.dirs }
func (w *wasiConfig) EnvVars() map[string]string     { return w.env }
func (w *wasiConfig) Workdir() string                { return w.workdir }

type hostConfig struct {
	imports HostImports
}

func (h *hostConfig) Imports() HostImports { return h.imports }

type hostImports struct {
	namespace string
	functions []sdk.Function
	globals   map[string]interface{}
	table     map[string]interface{}
	memory    map[string]interface{}
}

func (i *hostImports) Namespace() string               { return i.namespace }
func (i *hostImports) Functions() []sdk.Function       { return i.functions }
func (i *hostImports) Globals() map[string]interface{} { return i.globals }
func (i *hostImports) Table() map[string]interface{}   { return i.table }
func (i *hostImports) Memory() map[string]interface{}  { return i.memory }

type config struct {
	host HostConfig
	wasi WasiConfig
}

func (c *config) Host() HostConfig { return c.host }
func (c *config) Wasi() WasiConfig { return c.wasi }

func NewConfig(exports []sdk.Function) *config {
	return &config{
		host: &hostConfig{
			imports: &hostImports{
				namespace: "env",
				functions: exports,
			},
		},
		wasi: &wasiConfig{
			stdout: true,
			stderr: true,
		},
	}
}

type Function22 interface {
	Name() string
}

// map of namespace to obj ot imports, globals, etc.

type HostImports interface {
	Namespace() string
	Functions() []sdk.Function
	Globals() map[string]interface{}
	Table() map[string]interface{}
	Memory() map[string]interface{}
}

type HostConfig interface {
	// TODO
	// map of namespaces to functions
	Imports() HostImports
}

type WasiConfig interface {
	Arguments() []string
	CaptureStdout() bool
	CaptureStderr() bool
	Directories() map[string]string
	EnvVars() map[string]string
	Workdir() string
}

type Config interface {
	Host() HostConfig
	Wasi() WasiConfig
}
