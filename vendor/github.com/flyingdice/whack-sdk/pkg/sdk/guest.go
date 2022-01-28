package sdk

import (
	"github.com/pkg/errors"
)

const (
	FreeExport   = "whack_free"
	MainExport   = "whack_main"
	MallocExport = "whack_malloc"
)

// interface to manage memory better??
// allocator
//
// Malloc, Free
// Store (interface{}, error)
// Get(addr) (interface{}, error)

type guest struct {
	instance Instance
}

// Main invokes the whack main function inside the instance.
func (g *guest) Main(mem Memory) error {
	_, err := g.instance.Call(MainExport, mem.Address(), mem.Length())
	if err != nil {
		return errors.Wrapf(err, "failed to call %s", MainExport)
	}
	return nil
}

// Malloc allocates memory inside the instance.
func (g *guest) Malloc(length int32) (Memory, error) {
	ptr, err := g.instance.Call(MallocExport, length)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to call %s", MallocExport)
	}
	return NewMemory(ptr.(int32), length), nil
}

// Free deallocates memory inside the instance.
func (g *guest) Free(mem Memory) error {
	_, err := g.instance.Call(FreeExport, mem.Address(), mem.Length())
	if err != nil {
		return errors.Wrapf(err, "failed to call %s", FreeExport)
	}
	return nil
}

func (g *guest) Call(name string, args ...interface{}) (interface{}, error) {
	return g.instance.Call(name, args...)
}
func (g *guest) Close() error                    { return g.instance.Close() }
func (g *guest) Read(mem Memory) ([]byte, error) { return g.instance.Read(mem) }
func (g *guest) Stdout() string                  { return g.instance.Stdout() }
func (g *guest) Stderr() string                  { return g.instance.Stderr() }
func (g *guest) Write(addr int32, bytes []byte) (int32, error) {
	return g.instance.Write(addr, bytes)
}
func (g *guest) WRN() WRN { return g.instance.WRN() }

func NewGuest(instance Instance) *guest {
	return &guest{
		instance: instance,
	}
}
