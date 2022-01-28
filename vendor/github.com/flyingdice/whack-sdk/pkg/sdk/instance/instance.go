package instance

import "github.com/flyingdice/whack-sdk/pkg/sdk"

var _ sdk.Instance = (*instance)(nil)

// instance wraps a runtime specific instance and leverages result channels for
// interop between golang and WASM.
type instance struct {
	instance sdk.RuntimeInstance

	success chan interface{}
	error   chan error
}

func (i *instance) GetResult() (interface{}, error) {
	select {
	case s := <-i.success:
		return s, nil
	case e := <-i.error:
		return nil, e
	default:
		return nil, nil
	}
}

func (i *instance) SetResult(result interface{}, err error) {
	if result != nil {
		i.success <- result
	} else if err != nil {
		i.error <- err
	}
}

func (i *instance) Call(name string, args ...interface{}) (interface{}, error) {
	return i.instance.Call(name, args...)
}
func (i *instance) Close() error                        { return i.instance.Close() }
func (i *instance) Read(mem sdk.Memory) ([]byte, error) { return i.instance.Read(mem) }
func (i *instance) Stdout() string                      { return i.instance.Stdout() }
func (i *instance) Stderr() string                      { return i.instance.Stderr() }
func (i *instance) Write(addr int32, bytes []byte) (int32, error) {
	return i.instance.Write(addr, bytes)
}
func (i *instance) WRN() sdk.WRN { return i.instance.WRN() }

func NewInstance(ri sdk.RuntimeInstance) (*instance, error) {
	return &instance{
		instance: ri,
		success:  make(chan interface{}, 1),
		error:    make(chan error, 1),
	}, nil
}
