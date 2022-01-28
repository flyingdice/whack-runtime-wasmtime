package sdk

import (
	"github.com/pkg/errors"
)

type app struct {
	mmu      MMU
	guest    Guest
	instance Instance
}

func NewApp(instance Instance) *app {
	guest := NewGuest(instance)
	mmu := NewMMU(guest)

	return &app{
		guest:    guest,
		instance: instance,
		mmu:      mmu,
	}
}

func (a *app) Main(req Request) (res Response, err error) {
	var mem Memory

	defer func() {
		if mem != nil {
			e := a.mmu.Free(mem)
			if e != nil {
				err = errors.Wrap(e, "failed to free request memory")
			}
		}
	}()

	bytes, err := req.Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "failed to serialize request to bytes")
	}

	mem, err = a.mmu.Write(bytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to write request to memory")
	}

	err = a.guest.Main(mem)
	if err != nil {
		return nil, errors.Wrap(err, "failed to run app main")
	}

	result, err := a.instance.GetResult()
	if err != nil {
		res = Error(err)
	} else if result != nil {
		res = Success(result)
	}

	return
}
