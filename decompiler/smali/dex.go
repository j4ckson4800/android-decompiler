package smali

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/j4ckson4800/android-decompiler/decompiler/smali/internal"
)

const ResolveResource = "{{resolve_from_resource}}"

type Dex struct {
	rawDex   internal.Dex
	Filename string

	Classes map[string]Class
	Methods map[string]Method
	Fields  map[string]Field

	MethodsByIndex map[int]string
	FieldsByIndex  map[int]string
}

func NewDex(r *bytes.Reader, cfg Config) (Dex, error) {
	parser := NewParser(r)
	rawDex, err := internal.NewDex(parser)
	if err != nil {
		return Dex{}, fmt.Errorf("new dex: %w", err)
	}

	outDex := Dex{
		rawDex:  rawDex,
		Classes: make(map[string]Class, len(rawDex.ClassDefs)),
		Methods: make(map[string]Method, len(rawDex.MethodDefs)),
		Fields:  make(map[string]Field, len(rawDex.FieldDefs)),
		// NOTE: we usually have a lot of functions from android sdk
		// we don't need to allocate space for them because they don't have impl
		MethodsByIndex: make(map[int]string, len(rawDex.MethodDefs)/10),
		// NOTE: we only need static fields, so we'll reduce the size by some magic value
		FieldsByIndex: make(map[int]string, len(rawDex.FieldDefs)/10),
	}

	for _, classDef := range outDex.rawDex.ClassDefs {
		lowLevelClass, err := internal.NewClass(parser, classDef)
		if err != nil {
			return Dex{}, fmt.Errorf("new internal class: %w", err)
		}

		className := string(outDex.rawDex.StringDefs[outDex.rawDex.TypeIDs[classDef.Index]].Data)
		superClassName := ""
		if classDef.Super > 0 && classDef.Super < uint32(len(outDex.rawDex.TypeIDs)) {
			superClassName = string(outDex.rawDex.StringDefs[outDex.rawDex.TypeIDs[classDef.Super]].Data)
		}
		class, err := NewClass(className, superClassName)
		if err != nil {
			return Dex{}, fmt.Errorf("new class: %w", err)
		}

		class.Methods = make([]Method, 0, len(lowLevelClass.Methods)+len(lowLevelClass.VirtualMethods))
		class.Methods, err = outDex.parseClassMethods(class.Methods, lowLevelClass.Methods, className)
		if err != nil {
			return Dex{}, fmt.Errorf("parse simple methods: %w", err)
		}

		class.Methods, err = outDex.parseClassMethods(class.Methods, lowLevelClass.VirtualMethods, className)
		if err != nil {
			return Dex{}, fmt.Errorf("parse virtual methods: %w", err)
		}

		fieldIdx := 0
		class.StaticFields = make([]Field, 0, len(lowLevelClass.StaticFields))
		class.InstanceFields = make([]Field, 0, len(lowLevelClass.InstanceFields))
		for i, staticField := range lowLevelClass.StaticFields {
			fieldIdx += int(staticField.IndexDiff)
			def := outDex.rawDex.FieldDefs[fieldIdx]
			fieldName := string(outDex.rawDex.StringDefs[def.Name].Data)
			fieldType := string(outDex.rawDex.StringDefs[outDex.rawDex.TypeIDs[def.Type]].Data)

			value := int64(0)
			if len(lowLevelClass.StaticValues.Values) > i {
				value = lowLevelClass.StaticValues.Values[i].Value
			}

			sb := strings.Builder{}
			sb.WriteString(className)
			sb.WriteString("->")
			sb.WriteString(fieldName)
			sb.WriteString(":")
			sb.WriteString(fieldType)

			descriptor := sb.String()
			field := Field{
				DefIdx:     fieldIdx,
				Name:       fieldName,
				Type:       fieldType,
				ClassName:  className,
				Value:      value,
				Descriptor: descriptor,
			}

			outDex.Fields[descriptor] = field
			outDex.FieldsByIndex[fieldIdx] = descriptor
			class.StaticFields = append(class.StaticFields, field)
		}
		fieldIdx = 0
		for i, instanceField := range lowLevelClass.InstanceFields {
			fieldIdx += int(instanceField.IndexDiff)
			def := outDex.rawDex.FieldDefs[fieldIdx]
			fieldName := string(outDex.rawDex.StringDefs[def.Name].Data)
			fieldType := string(outDex.rawDex.StringDefs[outDex.rawDex.TypeIDs[def.Type]].Data)

			value := int64(0)
			if len(lowLevelClass.StaticValues.Values) > i {
				value = lowLevelClass.StaticValues.Values[i].Value
			}

			sb := strings.Builder{}
			sb.WriteString(className)
			sb.WriteString("->")
			sb.WriteString(fieldName)
			sb.WriteString(":")
			sb.WriteString(fieldType)

			descriptor := sb.String()
			field := Field{
				DefIdx:     fieldIdx,
				Name:       fieldName,
				Type:       fieldType,
				ClassName:  className,
				Value:      value,
				Descriptor: descriptor,
			}

			outDex.Fields[descriptor] = field
			outDex.FieldsByIndex[fieldIdx] = descriptor
			class.InstanceFields = append(class.InstanceFields, field)
		}

		outDex.Classes[className] = class
	}

	for i, methodDef := range outDex.rawDex.MethodDefs {
		if _, ok := outDex.MethodsByIndex[i]; ok {
			continue
		}

		proto := outDex.rawDex.MethodProtoDefs[methodDef.Type]

		className := string(outDex.rawDex.StringDefs[outDex.rawDex.TypeIDs[methodDef.Class]].Data)
		methodName := string(outDex.rawDex.StringDefs[methodDef.Name].Data)
		returnType := string(outDex.rawDex.StringDefs[outDex.rawDex.TypeIDs[proto.ReturnTypeIdx]].Data)

		classMethod, err := NewMethod(className, methodName, returnType, proto.ParamsString, internal.Method{})
		if err != nil {
			return Dex{}, fmt.Errorf("new method: %w", err)
		}

		methodSignature := outDex.getMethodSignature(className, i)
		outDex.Methods[methodSignature] = classMethod
		outDex.MethodsByIndex[i] = methodSignature
	}

	if cfg.SanitizeAnnotations {
		if err := rawDex.SanitizeAnnotations(); err != nil {
			return Dex{}, fmt.Errorf("sanitize annotations: %w", err)
		}
	}

	return outDex, nil
}

func (d *Dex) parseClassMethods(classMethods []Method, methods []internal.Method, className string) ([]Method, error) {
	methodIdx := 0
	for _, method := range methods {
		methodIdx += int(method.IndexDiff)
		def := d.rawDex.MethodDefs[methodIdx]
		proto := d.rawDex.MethodProtoDefs[def.Type]

		methodName := string(d.rawDex.StringDefs[def.Name].Data)
		returnType := string(d.rawDex.StringDefs[d.rawDex.TypeIDs[proto.ReturnTypeIdx]].Data)

		classMethod, err := NewMethod(className, methodName, returnType, d.rawDex.MethodProtoDefs[def.Type].ParamsString, method)
		if err != nil {
			return nil, fmt.Errorf("new method: %w", err)
		}

		methodSignature := d.getMethodSignature(className, methodIdx)
		classMethods = append(classMethods, classMethod)
		d.Methods[methodSignature] = classMethod
		d.MethodsByIndex[methodIdx] = methodSignature
	}

	return classMethods, nil
}

func (d *Dex) getMethodSignature(className string, methodIdx int) string {
	def := d.rawDex.MethodDefs[methodIdx]
	proto := d.rawDex.MethodProtoDefs[def.Type]

	methodName := string(d.rawDex.StringDefs[def.Name].Data)
	returnType := string(d.rawDex.StringDefs[d.rawDex.TypeIDs[proto.ReturnTypeIdx]].Data)

	sb := strings.Builder{}
	sb.WriteString(className)
	sb.WriteString("->")
	sb.WriteString(methodName)
	sb.WriteString("(")
	sb.WriteString(proto.ParamsString)
	sb.WriteString(")")
	sb.WriteString(returnType)

	return sb.String()
}
