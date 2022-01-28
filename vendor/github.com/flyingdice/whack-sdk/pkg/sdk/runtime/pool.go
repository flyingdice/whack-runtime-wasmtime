package runtime

import (
	"github.com/flyingdice/whack-sdk/pkg/sdk"
	"github.com/flyingdice/whack-sdk/pkg/sdk/instance"
	"github.com/pkg/errors"
	"sync"
)

var _ sdk.InstancePool = (*pool)(nil)

type pool struct {
	instances map[sdk.WRN]sdk.Instance
	mu        *sync.Mutex
}

func NewPool() (*pool, error) {
	return &pool{
		instances: make(map[sdk.WRN]sdk.Instance),
		mu:        &sync.Mutex{},
	}, nil
}

func (p *pool) Set(ri sdk.RuntimeInstance) (sdk.Instance, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	i, err := instance.NewInstance(ri)
	if err != nil {
		return nil, err
	}

	p.instances[i.WRN()] = i
	return i, nil
}

func (p *pool) Get(wrn sdk.WRN) (sdk.Instance, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	i, ok := p.instances[wrn]
	if !ok {
		return nil, errors.Errorf("instance for wrn=%p not found", wrn)
	}

	return i, nil
}

func (p *pool) Return(i sdk.Instance) error {
	return nil
}
