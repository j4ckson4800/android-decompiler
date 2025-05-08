package internal

import (
	"fmt"
)

type Field struct {
	IndexDiff   uint64
	AccessFlags uint64
}

func NewField(p Parser) (Field, error) {
	indexDiff, err := p.ReadULEB128()
	if err != nil {
		return Field{}, fmt.Errorf("read uleb128: %w", err)
	}

	accessFlags, err := p.ReadULEB128()
	if err != nil {
		return Field{}, fmt.Errorf("read uleb128: %w", err)
	}

	return Field{
		IndexDiff:   indexDiff,
		AccessFlags: accessFlags,
	}, nil
}
