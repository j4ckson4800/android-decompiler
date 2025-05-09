package internal

import (
	"fmt"
)

type EncodedAnnotationHeader struct {
	TypeID uint64
	Size   uint64
}

type AnnotationElement struct {
	NameID uint64
	Value  Value
}

type AnnotationValue struct {
	Header   EncodedAnnotationHeader
	Elements []AnnotationElement
}

type Annotation struct {
	Visibility byte
	AnnotationValue
}

func newAnnotationValue(p Parser) (AnnotationValue, error) {
	typeID, err := p.ReadULEB128()
	if err != nil {
		return AnnotationValue{}, fmt.Errorf("read uleb128: %w", err)
	}

	header := EncodedAnnotationHeader{}

	header.TypeID = typeID

	size, err := p.ReadULEB128()
	if err != nil {
		return AnnotationValue{}, fmt.Errorf("read uleb128: %w", err)
	}

	header.Size = size

	elements := make([]AnnotationElement, header.Size)
	for i := range elements {
		element := AnnotationElement{}

		nameID, err := p.ReadULEB128()
		if err != nil {
			return AnnotationValue{}, fmt.Errorf("read uleb128: %w", err)
		}

		element.NameID = nameID
		val, err := NewValue(p)
		if err != nil {
			return AnnotationValue{}, fmt.Errorf("new value: %w", err)
		}

		element.Value = val
		elements[i] = element
	}

	return AnnotationValue{
		Header:   header,
		Elements: elements,
	}, nil
}

func NewAnnotation(p Parser) (Annotation, error) {
	visibility, err := p.ReadByte()
	if err != nil {
		return Annotation{}, fmt.Errorf("read byte: %w", err)
	}

	val, err := newAnnotationValue(p)
	if err != nil {
		return Annotation{}, fmt.Errorf("new annotation value: %w", err)
	}

	return Annotation{
		Visibility:      visibility,
		AnnotationValue: val,
	}, nil
}
