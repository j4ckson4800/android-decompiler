package defs

import (
	"fmt"
)

type StringOffset uint32

type StringDef struct {
	Data []byte
}

func NewStringOffset(p Parser) (StringOffset, error) {
	offset, err := p.ReadUint32()
	if err != nil {
		return 0, fmt.Errorf("read uint32: %w", err)
	}
	return StringOffset(offset), nil
}

func NewStringDef(p Parser) (StringDef, error) {
	size, err := p.ReadULEB128()
	if err != nil {
		return StringDef{}, fmt.Errorf("read uleb128: %w", err)
	}

	data, err := p.ReadBytes(int64(size))
	if err != nil {
		return StringDef{}, fmt.Errorf("read bytes: %w", err)
	}
	return StringDef{
		Data: data,
	}, nil
}
