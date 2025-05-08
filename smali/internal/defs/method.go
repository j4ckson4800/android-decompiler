package defs

import (
	"fmt"
)

type MethodDef FieldDef

type MethodProtoDef struct {
	Shorty        uint32
	ReturnTypeIdx uint32
	Params        []uint16
	ParamsString  string
}

type ProtoDef struct {
	Shorty       uint32
	Return       uint32
	ParamsOffset uint32
}

func NewMethodDef(r Parser) (MethodDef, error) {
	field, err := NewFieldDef(r)
	if err != nil {
		return MethodDef{}, fmt.Errorf("new field def: %w", err)
	}

	return MethodDef(field), nil
}

func NewMethodProtoDef(r Parser) (MethodProtoDef, error) {
	count, err := r.ReadUint32()
	if err != nil {
		return MethodProtoDef{}, fmt.Errorf("read uint32: %w", err)
	}

	params := make([]uint16, 0, count)
	for i := 0; i < int(count); i++ {
		param, err := r.ReadUint16()
		if err != nil {
			return MethodProtoDef{}, fmt.Errorf("read uint16: %w", err)
		}

		params = append(params, param)
	}

	return MethodProtoDef{
		Params: params,
	}, nil
}

func NewProtoDef(r Parser) (ProtoDef, error) {
	proto := ProtoDef{}
	if err := r.ReadStruct(&proto); err != nil {
		return ProtoDef{}, fmt.Errorf("read struct: %w", err)
	}

	return proto, nil
}
