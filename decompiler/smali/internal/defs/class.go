package defs

import (
	"fmt"
)

type ClassDef struct {
	Index              uint32
	AccessFlags        uint32
	Super              uint32
	InterfacesOffset   uint32
	SourceFileIndex    uint32
	AnnotationsOffset  uint32
	ClassDataOffset    uint32
	StaticValuesOffset uint32
} // Size: 0x20

func NewClassDef(p Parser) (ClassDef, error) {
	class := ClassDef{}
	if err := p.ReadStruct(&class); err != nil {
		return ClassDef{}, fmt.Errorf("read struct: %w", err)
	}

	return class, nil
}
