package internal

import (
	"errors"
	"fmt"
	"strings"

	"github.com/j4ckson4800/android-decompiler/smali/internal/defs"
)

type defCreator[T any] func(p defs.Parser) (T, error)

type Dex struct {
	Header           defs.DexHeader
	StringDefs       []defs.StringDef
	MethodProtoDefs  []defs.MethodProtoDef
	MethodDefs       []defs.MethodDef
	ClassDefs        []defs.ClassDef
	FieldDefs        []defs.FieldDef
	TypeIDs          []uint32
	parser           Parser
	AuxiliaryStrings map[int]struct{}
}

func NewDex(r Parser) (Dex, error) {
	header, err := defs.NewDexHeader(r)
	if err != nil {
		return Dex{}, err
	}
	if header.HeaderSize != defs.DexHeaderSize {
		return Dex{}, errors.New("invalid header size")
	}

	dex := Dex{
		Header: header,
		parser: r,
	}

	if err := dex.parseStrings(); err != nil {
		return Dex{}, fmt.Errorf("parse strings: %w", err)
	}
	if err := dex.parseTypeIDs(); err != nil {
		return Dex{}, fmt.Errorf("parse type ids: %w", err)
	}
	if err := dex.parseMethodProtoDefs(); err != nil {
		return Dex{}, fmt.Errorf("parse method proto defs: %w", err)
	}
	if err := dex.parseMethodDefs(); err != nil {
		return Dex{}, fmt.Errorf("parse method defs: %w", err)
	}
	if err := dex.parseClassDefs(); err != nil {
		return Dex{}, fmt.Errorf("parse class defs: %w", err)
	}
	if err := dex.parseFieldDefs(); err != nil {
		return Dex{}, fmt.Errorf("parse field defs: %w", err)
	}

	return dex, nil
}

func (d *Dex) parseStrings() error {
	stringOffs, err := parseDef(d.parser, defs.NewStringOffset, d.Header.StringIds)
	if err != nil {
		return fmt.Errorf("parse defs: %w", err)
	}

	d.StringDefs = make([]defs.StringDef, 0, d.Header.StringIds.Size)
	for _, off := range stringOffs {
		if err := d.parser.SetCursorTo(int64(off)); err != nil {
			return fmt.Errorf("set cursor to: %w", err)
		}

		stringDef, err := defs.NewStringDef(d.parser)
		if err != nil {
			return fmt.Errorf("new string def: %w", err)
		}

		d.StringDefs = append(d.StringDefs, stringDef)
	}

	return nil
}

func (d *Dex) parseTypeIDs() error {
	if err := d.parser.SetCursorTo(int64(d.Header.TypeIds.Offset)); err != nil {
		return fmt.Errorf("set cursor to: %w", err)
	}

	typeIds := make([]uint32, 0, d.Header.TypeIds.Size)
	d.AuxiliaryStrings = make(map[int]struct{}, d.Header.TypeIds.Size)
	for i := 0; i < int(d.Header.TypeIds.Size); i++ {
		typeId, err := d.parser.ReadUint32()
		if err != nil {
			return fmt.Errorf("read uint32: %w", err)
		}

		d.AuxiliaryStrings[int(typeId)] = struct{}{}
		typeIds = append(typeIds, typeId)
	}

	d.TypeIDs = typeIds
	return nil
}

func (d *Dex) parseMethodProtoDefs() error {
	protoDefs, err := parseDef(d.parser, defs.NewProtoDef, d.Header.ProtoIds)
	if err != nil {
		return fmt.Errorf("parse defs: %w", err)
	}

	methodProtoDefs := make([]defs.MethodProtoDef, 0, d.Header.ProtoIds.Size)
	for i := 0; i < int(d.Header.ProtoIds.Size); i++ {
		methodProto := defs.MethodProtoDef{}
		if protoDefs[i].ParamsOffset != 0 {
			if err := d.parser.SetCursorTo(int64(protoDefs[i].ParamsOffset)); err != nil {
				return fmt.Errorf("set cursor to: %w", err)
			}

			methodDef, err := defs.NewMethodProtoDef(d.parser)
			if err != nil {
				return fmt.Errorf("new method proto def: %w", err)
			}
			methodProto = methodDef
		}

		methodProto.ReturnTypeIdx = protoDefs[i].Return
		methodProto.Shorty = protoDefs[i].Shorty
		d.AuxiliaryStrings[int(methodProto.Shorty)] = struct{}{}

		sb := strings.Builder{}
		for _, param := range methodProto.Params {
			_, _ = sb.Write(d.StringDefs[d.TypeIDs[param]].Data)
		}

		methodProto.ParamsString = sb.String()
		methodProtoDefs = append(methodProtoDefs, methodProto)
	}

	d.MethodProtoDefs = methodProtoDefs
	return nil
}

func (d *Dex) parseFieldDefs() error {
	fieldDefs, err := parseDef(d.parser, defs.NewFieldDef, d.Header.FieldIds)
	if err != nil {
		return fmt.Errorf("parse defs: %w", err)
	}
	for _, field := range fieldDefs {
		d.AuxiliaryStrings[int(field.Name)] = struct{}{}
	}
	d.FieldDefs = fieldDefs
	return nil
}

func (d *Dex) parseMethodDefs() error {
	methodDefs, err := parseDef(d.parser, defs.NewMethodDef, d.Header.MethodIds)
	if err != nil {
		return fmt.Errorf("parse defs: %w", err)
	}
	for _, method := range methodDefs {
		d.AuxiliaryStrings[int(method.Name)] = struct{}{}
	}

	d.MethodDefs = methodDefs
	return nil
}

func (d *Dex) visitArray(array *Array) {
	for _, value := range array.Values {
		switch value.Type {
		case ValueTypeArray:
			d.visitArray(value.ArrayValue)
		case ValueTypeString:
			d.AuxiliaryStrings[int(value.Value)] = struct{}{}
		case ValueTypeAnnotation:
			d.visitAnnotationValue(value.AnnotationValue)
		default:
			continue
		}
	}
}

func (d *Dex) visitAnnotationValue(annotationValue *AnnotationValue) {
	for _, element := range annotationValue.Elements {
		d.AuxiliaryStrings[int(element.NameID)] = struct{}{}

		switch element.Value.Type {
		case ValueTypeArray:
			d.visitArray(element.Value.ArrayValue)
		case ValueTypeString:
			d.AuxiliaryStrings[int(element.Value.Value)] = struct{}{}
		case ValueTypeAnnotation:
			d.visitAnnotationValue(element.Value.AnnotationValue)
		default:
			continue
		}
	}
}

func (d *Dex) parseAnnotationSet(annotationSet *defs.AnnotationSetDef) error {
	for _, offset := range annotationSet.Offsets {
		if offset == 0 {
			continue
		}

		if err := d.parser.SetCursorTo(int64(offset)); err != nil {
			return fmt.Errorf("set cursor to: %w", err)
		}

		annotation, err := NewAnnotation(d.parser)
		if err != nil {
			return fmt.Errorf("new annotation: %w", err)
		}

		d.visitAnnotationValue(&annotation.AnnotationValue)
	}

	return nil
}

func (d *Dex) parseAnnotationTables(tables []defs.AnnotationTable) error {
	for _, table := range tables {
		if table.Offset == 0 {
			continue
		}

		if err := d.parser.SetCursorTo(int64(table.Offset)); err != nil {
			return fmt.Errorf("set cursor to: %w", err)
		}

		annotationSet, err := defs.NewAnnotationSetDef(d.parser)
		if err != nil {
			return fmt.Errorf("new annotation table: %w", err)
		}

		for _, offset := range annotationSet.Offsets {
			if offset == 0 {
				continue
			}

			if err := d.parser.SetCursorTo(int64(offset)); err != nil {
				return fmt.Errorf("set cursor to: %w", err)
			}

			if err := d.parseAnnotationSet(&annotationSet); err != nil {
				return fmt.Errorf("parse annotation set: %w", err)
			}
		}
	}

	return nil
}

func (d *Dex) parseAnnotations(classDef defs.ClassDef) error {
	if classDef.AnnotationsOffset == 0 {
		return nil
	}

	if err := d.parser.SetCursorTo(int64(classDef.AnnotationsOffset)); err != nil {
		return fmt.Errorf("set cursor to: %w", err)
	}

	annotationsDirectory, err := defs.NewAnnotationDef(d.parser)
	if err != nil {
		return fmt.Errorf("new annotations: %w", err)
	}

	if err := d.parseAnnotationTables(annotationsDirectory.Tables.Methods); err != nil {
		return fmt.Errorf("parse method annotations: %w", err)
	}
	if err := d.parseAnnotationTables(annotationsDirectory.Tables.Fields); err != nil {
		return fmt.Errorf("parse field annotations: %w", err)
	}
	if err := d.parseAnnotationTables(annotationsDirectory.Tables.Parameters); err != nil {
		return fmt.Errorf("parse parameter annotations: %w", err)
	}

	if annotationsDirectory.Dir.ClassAnnotations == 0 {
		return nil
	}

	if err := d.parser.SetCursorTo(int64(annotationsDirectory.Dir.ClassAnnotations)); err != nil {
		return fmt.Errorf("set cursor to: %w", err)
	}

	annotationSet, err := defs.NewAnnotationSetDef(d.parser)
	if err != nil {
		return fmt.Errorf("new annotation set: %w", err)
	}
	if err := d.parseAnnotationSet(&annotationSet); err != nil {
		return fmt.Errorf("parse annotation set: %w", err)
	}

	return nil
}

func (d *Dex) SanitizeAnnotations() error {
	const noIndex = ^uint32(0)

	for _, classDef := range d.ClassDefs {

		if classDef.SourceFileIndex != noIndex {
			d.AuxiliaryStrings[int(classDef.SourceFileIndex)] = struct{}{}
		}

		if err := d.parseAnnotations(classDef); err != nil {
			return fmt.Errorf("parse annotations: %w", err)
		}
	}

	return nil
}

func (d *Dex) parseClassDefs() error {
	classDefs, err := parseDef(d.parser, defs.NewClassDef, d.Header.ClassDefs)
	if err != nil {
		return fmt.Errorf("parse defs: %w", err)
	}

	d.ClassDefs = classDefs
	return nil
}

func parseDef[T any](r Parser, c defCreator[T], table defs.Table) ([]T, error) {
	if err := r.SetCursorTo(int64(table.Offset)); err != nil {
		return nil, fmt.Errorf("set cursor to: %w", err)
	}

	definitions := make([]T, 0, table.Size)
	for i := 0; i < int(table.Size); i++ {
		def, err := c(r)
		if err != nil {
			return nil, fmt.Errorf("new class: %w", err)
		}

		definitions = append(definitions, def)
	}

	return definitions, nil
}
