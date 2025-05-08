package internal

import (
	"fmt"
)

type Array struct {
	Size   uint64
	Values []Value
}

func NewArray(p Parser) (Array, error) {
	size, err := p.ReadULEB128()
	if err != nil {
		return Array{}, fmt.Errorf("read uleb128: %w", err)
	}

	values := make([]Value, size)
	for i := range values {
		value, err := NewValue(p)
		if err != nil {
			return Array{}, fmt.Errorf("new value: %w", err)
		}
		values[i] = value
	}

	return Array{
		Size:   size,
		Values: values,
	}, nil
}
