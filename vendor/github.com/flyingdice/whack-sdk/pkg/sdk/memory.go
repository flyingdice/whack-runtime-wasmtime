package sdk

import (
	"encoding/json"
	"github.com/pkg/errors"
	"log"
)

var _ Memory = (*memory)(nil)
var _ MMU = (*mmu)(nil)

type memory struct {
	address int32
	length  int32
}

func (m *memory) Address() int32 { return m.address }
func (m *memory) Length() int32  { return m.length }

func NewMemory(address, length int32) *memory {
	return &memory{
		address: address,
		length:  length,
	}
}

// mmu TODO
type mmu struct {
	guest Guest
}

func (m *mmu) Free(mem Memory) error               { return m.guest.Free(mem) }
func (m *mmu) Malloc(length int32) (Memory, error) { return m.guest.Malloc(length) }

func (m *mmu) Read(mem Memory) (interface{}, error) {
	bytes, err := m.guest.Read(mem)
	if err != nil {
		return nil, errors.Wrap(err, "failed to write obj memory")
	}

	var obj DeserializedRequest

	err = deserialize(bytes, obj)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deserialize obj from bytes")
	}

	return obj, nil
}

func (m *mmu) Write(obj interface{}) (Memory, error) {
	log.Print("mmu write")
	log.Print(obj)

	bytes, err := serialize(obj)
	if err != nil {
		return nil, errors.Wrap(err, "failed to serialize obj to bytes")
	}

	log.Print("aft serialize")
	log.Print(string(bytes))

	mem, err := m.guest.Malloc(int32(len(bytes)))
	if err != nil {
		return nil, errors.Wrap(err, "failed to malloc obj memory")
	}

	log.Print("done malloc")
	log.Print(mem)
	log.Print(mem.Address())
	log.Print(mem.Length())

	length, err := m.guest.Write(mem.Address(), bytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to write obj bytes to guest memory")
	}

	log.Print("mmu wrote")
	log.Print(length)
	log.Print(err)

	if length != mem.Length() {
		return nil, errors.Errorf("expected to write %d bytes; wrote %d", mem.Length(), length)
	}

	return mem, nil
}

func NewMMU(guest Guest) *mmu {
	return &mmu{
		guest: guest,
	}
}

func serialize(obj interface{}) ([]byte, error) {
	switch o := obj.(type) {
	case []byte:
		return o, nil
	case string:
		return []byte(o), nil
	default:
		return json.Marshal(o)
	}
}

func deserialize(bytes []byte, obj interface{}) error {
	return json.Unmarshal(bytes, obj)
}
