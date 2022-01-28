package internal

import "github.com/bytecodealliance/wasmtime-go"

type global struct {
	value interface{}
}

func (g *global) Bind(namespace string, linker *wasmtime.Linker) error {
	return nil
}

func NewGlobal(value interface{}) (*global, error) {
	return &global{
		value: value,
	}, nil
}
