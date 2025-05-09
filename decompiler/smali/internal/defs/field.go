package defs

import (
	"fmt"
)

type FieldDef struct {
	Class uint16
	Type  uint16
	Name  uint32
}

func NewFieldDef(p Parser) (FieldDef, error) {
	field := FieldDef{}
	if err := p.ReadStruct(&field); err != nil {
		return FieldDef{}, fmt.Errorf("read struct: %w", err)
	}

	return field, nil
}
