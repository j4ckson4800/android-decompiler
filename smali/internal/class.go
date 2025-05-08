package internal

import (
	"fmt"

	"github.com/j4ckson4800/android-decompiler/smali/internal/defs"
)

type Class struct {
	StaticFields   []Field
	InstanceFields []Field
	Methods        []Method
	VirtualMethods []Method
	StaticValues   Array

	RawClass defs.ClassDef
}

func NewClass(p Parser, def defs.ClassDef) (Class, error) {
	if def.ClassDataOffset == 0 {
		return Class{}, nil
	}

	if err := p.SetCursorTo(int64(def.ClassDataOffset)); err != nil {
		return Class{}, fmt.Errorf("set cursor: %w", err)
	}

	staticFieldsSize, err := p.ReadULEB128()
	if err != nil {
		return Class{}, fmt.Errorf("read uleb128: %w", err)
	}

	instanceFieldsSize, err := p.ReadULEB128()
	if err != nil {
		return Class{}, fmt.Errorf("read uleb128: %w", err)
	}

	methodsSize, err := p.ReadULEB128()
	if err != nil {
		return Class{}, fmt.Errorf("read uleb128: %w", err)
	}

	virtualMethodsSize, err := p.ReadULEB128()
	if err != nil {
		return Class{}, fmt.Errorf("read uleb128: %w", err)
	}

	cls := Class{
		StaticFields:   make([]Field, staticFieldsSize),
		InstanceFields: make([]Field, instanceFieldsSize),
		Methods:        make([]Method, methodsSize),
		VirtualMethods: make([]Method, virtualMethodsSize),
		RawClass:       def,
	}

	for i := range cls.StaticFields {
		field, err := NewField(p)
		if err != nil {
			return Class{}, fmt.Errorf("new field: %w", err)
		}
		cls.StaticFields[i] = field
	}

	for i := range cls.InstanceFields {
		field, err := NewField(p)
		if err != nil {
			return Class{}, fmt.Errorf("new field: %w", err)
		}
		cls.InstanceFields[i] = field
	}

	for i := range cls.Methods {
		method, err := NewMethod(p)
		if err != nil {
			return Class{}, fmt.Errorf("new method: %w", err)
		}
		cls.Methods[i] = method
	}

	for i := range cls.VirtualMethods {
		method, err := NewMethod(p)
		if err != nil {
			return Class{}, fmt.Errorf("new method: %w", err)
		}
		cls.VirtualMethods[i] = method
	}

	if err := cls.parseMethods(p); err != nil {
		return Class{}, fmt.Errorf("parse methods: %w", err)
	}

	if def.StaticValuesOffset == 0 {
		return cls, nil
	}

	if err := p.SetCursorTo(int64(def.StaticValuesOffset)); err != nil {
		return Class{}, fmt.Errorf("set cursor: %w", err)
	}

	staticValues, err := NewArray(p)
	if err != nil {
		return Class{}, fmt.Errorf("new array: %w", err)
	}

	cls.StaticValues = staticValues
	return cls, nil
}

func (c *Class) parseMethods(p Parser) error {
	for i := range c.Methods {
		if err := c.Methods[i].ParseCode(p); err != nil {
			return fmt.Errorf("parse code: %w", err)
		}
	}

	for i := range c.VirtualMethods {
		if err := c.VirtualMethods[i].ParseCode(p); err != nil {
			return fmt.Errorf("parse code: %w", err)
		}
	}

	return nil
}
