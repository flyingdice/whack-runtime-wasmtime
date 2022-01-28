package sdk

import "github.com/google/uuid"

// WRN (Whack Resource Name) uniquely identify Whack resources.
//
// wrn:partition:service:region:account-id:resource-id
// <type>:<nid>:<nss>:
// wrn:pkg:todo:todo
type wrn struct {
	name string
}

func (w *wrn) Name() string {
	return w.name
}

func NewWRN(name string) *wrn {
	return &wrn{
		name: name,
	}
}

func RandomWRN() *wrn {
	return &wrn{
		name: uuid.NewString(),
	}
}
