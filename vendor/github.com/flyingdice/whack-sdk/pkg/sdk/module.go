package sdk

import (
	"github.com/pkg/errors"
	"os"
)

var _ Module = (*module)(nil)

var ErrEmptyPath = errors.New("module file path must be set")

type module struct {
	wrn   WRN
	bytes []byte
}

func (m *module) Wrn() WRN      { return m.wrn }
func (m *module) Bytes() []byte { return m.bytes }

// NewModule creates a new Module from the content at the given file path.
func NewModule(wrn WRN, path string) (*module, error) {
	if path == "" {
		return nil, ErrEmptyPath
	}
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read module path %s", path)
	}
	return &module{
		wrn:   wrn,
		bytes: bytes,
	}, nil
}
