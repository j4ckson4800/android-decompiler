package defs

import (
	"fmt"
)

type AnnotationsDirectory struct {
	ClassAnnotations uint32 // -> AnnotationSet
	FieldsSize       uint32
	MethodsSize      uint32
	ParametersSize   uint32
}

type AnnotationTables struct {
	Fields     []AnnotationTable
	Methods    []AnnotationTable
	Parameters []AnnotationTable
}

type AnnotationTable struct {
	Index  uint32 // unused
	Offset uint32 // -> AnnotationSet
}

type AnnotationSet struct {
	Size uint32
}

type AnnotationDef struct {
	Dir    AnnotationsDirectory
	Tables AnnotationTables
}

type AnnotationSetDef struct {
	Set     AnnotationSet
	Offsets []uint32
}

func NewAnnotationDef(p Parser) (AnnotationDef, error) {
	dir := AnnotationsDirectory{}
	if err := p.ReadStruct(&dir); err != nil {
		return AnnotationDef{}, fmt.Errorf("read struct: %w", err)
	}

	tables := AnnotationTables{
		Fields:     make([]AnnotationTable, dir.FieldsSize),
		Methods:    make([]AnnotationTable, dir.MethodsSize),
		Parameters: make([]AnnotationTable, dir.ParametersSize),
	}

	if err := p.ReadStruct(&tables.Fields); err != nil {
		return AnnotationDef{}, fmt.Errorf("read fields: %w", err)
	}
	if err := p.ReadStruct(&tables.Methods); err != nil {
		return AnnotationDef{}, fmt.Errorf("read methods: %w", err)
	}
	if err := p.ReadStruct(&tables.Parameters); err != nil {
		return AnnotationDef{}, fmt.Errorf("read parameters: %w", err)
	}

	return AnnotationDef{
		Dir:    dir,
		Tables: tables,
	}, nil
}

func NewAnnotationSetDef(p Parser) (AnnotationSetDef, error) {
	set := AnnotationSet{}
	if err := p.ReadStruct(&set); err != nil {
		return AnnotationSetDef{}, fmt.Errorf("read struct: %w", err)
	}

	offsets := make([]uint32, set.Size)
	if err := p.ReadStruct(&offsets); err != nil {
		return AnnotationSetDef{}, fmt.Errorf("read struct: %w", err)
	}

	return AnnotationSetDef{
		Set:     set,
		Offsets: offsets,
	}, nil
}
