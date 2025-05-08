package internal

import (
	"fmt"
)

type ValueType int8

const (
	ValueTypeByte ValueType = iota
	_
	ValueTypeShort
	ValueTypeChar
	ValueTypeInt
	_
	ValueTypeLong
	_
	_
	_
	_
	_
	_
	_
	_
	_
	ValueTypeFloat
	ValueTypeDouble
	_
	_
	_
	ValueTypeMethodType
	ValueTypeMethodHandle
	ValueTypeString
	ValueTypeType
	ValueTypeField
	ValueTypeMethod
	ValueTypeEnum
	ValueTypeArray
	ValueTypeAnnotation
	ValueTypeNull
	ValueTypeBoolean
)

type Value struct {
	Type            ValueType
	Size            byte
	Pad             int16
	Value           int64
	ArrayValue      *Array
	AnnotationValue *AnnotationValue
}

func NewValue(p Parser) (Value, error) {
	b, err := p.ReadByte()
	if err != nil {
		return Value{}, fmt.Errorf("read byte: %w", err)
	}

	val := Value{
		Type: ValueType(b & 0x1f),
		Size: (b >> 5) & 0x7,
	}

	switch val.Type {
	case ValueTypeByte:
		v, err := p.ReadByte()
		if err != nil {
			return Value{}, fmt.Errorf("read byte: %w", err)
		}
		val.Value = int64(v)
	case ValueTypeLong, ValueTypeDouble,
		ValueTypeShort,
		ValueTypeChar, ValueTypeInt, ValueTypeFloat,
		ValueTypeMethodType, ValueTypeMethodHandle,
		ValueTypeString, ValueTypeType, ValueTypeField, ValueTypeMethod, ValueTypeEnum:
		for i := 0; i <= int(val.Size); i++ {
			v, err := p.ReadByte()
			if err != nil {
				return Value{}, fmt.Errorf("read byte: %w", err)
			}
			val.Value |= int64(v) << (i * 8)
		}
	case ValueTypeBoolean:
		val.Value = int64(val.Size & 0x1)
	case ValueTypeArray:
		arr, err := NewArray(p)
		if err != nil {
			return Value{}, fmt.Errorf("new array: %w", err)
		}
		val.ArrayValue = &arr
	case ValueTypeAnnotation:
		annotation, err := newAnnotationValue(p)
		if err != nil {
			return Value{}, fmt.Errorf("new annotation value: %w", err)
		}
		val.AnnotationValue = &annotation
	case ValueTypeNull:
	}

	return val, nil
}
