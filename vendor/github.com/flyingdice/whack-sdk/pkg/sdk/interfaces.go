package sdk

type MMU interface {
	Free(mem Memory) error
	Malloc(length int32) (Memory, error)

	Read(mem Memory) (interface{}, error)
	Write(obj interface{}) (Memory, error)
}

// Guest is the sdk side implementation of the guest SDK API.
type Guest interface {
	RuntimeInstance

	Free(mem Memory) error
	Main(mem Memory) error
	Malloc(length int32) (Memory, error)
}

type Request interface {
	Bytes() ([]byte, error)
}

type Response interface {
	Success() interface{}
	Error() error
}

type App interface {
	Main(Request) Response
}

type Memory interface {
	Address() int32
	Length() int32
}

// Module is an interface that represents a WASM module load into memory.
type Module interface {
	Wrn() WRN
	Bytes() []byte
}

// WRN (Whack Resource Name)  is an interface uniquely identify whack resources.
type WRN interface {
	Name() string
}

// Runtime represents a WASM runtime implementation.
type Runtime interface {
	New() (Instance, error)
	Get(wrn WRN) (Instance, error)
}

// RuntimeInstance represents a WASM instance for a specific runtime.
type RuntimeInstance interface {
	Call(name string, args ...interface{}) (interface{}, error)
	Close() error
	Read(mem Memory) ([]byte, error)
	Stdout() string
	Stderr() string
	Write(addr int32, bytes []byte) (int32, error)
	WRN() WRN
}

// Instance represents a generic WASM instance.
type Instance interface {
	RuntimeInstance

	GetResult() (interface{}, error)
	SetResult(interface{}, error)
}

// InstancePool represents a pool of WASM instances.
type InstancePool interface {
	Get(wrn WRN) (Instance, error)
	Set(ri RuntimeInstance) (Instance, error)
	Return(i Instance) error
}
