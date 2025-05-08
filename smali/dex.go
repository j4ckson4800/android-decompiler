package smali

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/j4ckson4800/android-decompiler/smali/internal"
)

const ResolveResource = "{{resolve_from_resource}}"
const nonEmptyGibberishValue = "{{DIRTY}}"

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
		class, err := NewClass(className)
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
		class.Fields = make([]Field, 0, len(lowLevelClass.StaticFields))
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
				Name:       fieldName,
				Type:       fieldType,
				ClassName:  className,
				Value:      value,
				Descriptor: descriptor,
			}

			outDex.Fields[descriptor] = field
			outDex.FieldsByIndex[fieldIdx] = descriptor
			class.Fields = append(class.Fields, field)
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

		classMethod, err := NewMethod(className, methodName, returnType, internal.Method{})
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

		classMethod, err := NewMethod(className, methodName, returnType, method)
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

func (d *Dex) findFirstRegUpdateValue(method *Method, from int, regs map[int]string) map[int]string {
	for j := from - 1; j >= 0; j-- {
		prevInst := method.Body[j]
		if len(prevInst.Operands) != 2 && prevInst.Type != TypeMoveResult {
			continue
		}
		if v, ok := regs[int(prevInst.Operands[0])]; !ok || v != "" {
			continue
		}

		switch prevInst.Type {
		case TypeConst:
			switch prevInst.Opcode {
			case OpConstString, OpConstStringJumbo:
				regs[int(prevInst.Operands[0])] = string(d.rawDex.StringDefs[prevInst.Operands[1]].Data)
			case OpConstClass, OpConstMethodHandle, OpConstMethodType:
				// skip this shit, but mark with some magic empty value
				// NOTE: we don't want to return with error over here
				// because we don't drop instance register because we always may face `-static` instr
				regs[int(prevInst.Operands[0])] = nonEmptyGibberishValue
			default:
				regs[int(prevInst.Operands[0])] = strconv.FormatInt(prevInst.Operands[1], 10)
			}
		case TypeStaticOp:
			switch prevInst.Opcode {
			case OpSgetObject:
				field := d.Fields[d.FieldsByIndex[int(prevInst.Operands[1])]]
				if field.Type != "Ljava/lang/String;" {
					// skip non-string objects0
					regs[int(prevInst.Operands[0])] = nonEmptyGibberishValue
					continue
				}

				regs[int(prevInst.Operands[0])] = string(d.rawDex.StringDefs[field.Value].Data)
			case OpSget, OpSgetBoolean, OpSgetByte, OpSgetChar, OpSgetShort, OpSgetWide:
				field := d.Fields[d.FieldsByIndex[int(prevInst.Operands[1])]]
				regs[int(prevInst.Operands[0])] = strconv.FormatInt(field.Value, 10)
			default:
				// skip this shit
				continue
			}
		case TypeMoveResult:
			switch prevInst.Opcode {
			case OpMoveResultObject:
				if j == 0 {
					continue
				}
				callInst := method.Body[j-1]
				if callInst.Type != TypeInvocation {
					continue
				}
				methodIdx := int(callInst.Operands[len(callInst.Operands)-1])
				sig, hasMethodIdx := d.MethodsByIndex[methodIdx]
				if sig == "" || !hasMethodIdx {
					continue
				}

				if sig != "Landroid/content/Context;->getString(I)Ljava/lang/String;" {
					continue
				}

				if len(callInst.Operands) < 3 {
					continue
				}

				lookForReg := make(map[int]string)
				lookForReg[int(callInst.Operands[1])] = ""

				// not sure if we want to fail here
				vals := d.findFirstRegUpdateValue(method, j-1, lookForReg)
				regs[int(prevInst.Operands[0])] = ResolveResource + ":" + vals[int(callInst.Operands[1])]
			default:
				regs[int(prevInst.Operands[0])] = nonEmptyGibberishValue
			}
		case TypeArithmetics, TypeComparison:
			regs[int(prevInst.Operands[0])] = nonEmptyGibberishValue
		case TypeInstanceOp, TypeArrayOp:
			switch prevInst.Opcode {
			case OpIget, OpIgetBoolean, OpIgetByte, OpIgetChar, OpIgetShort, OpIgetWide, OpIgetObject, OpNewInstance, OpInstanceOf,
				OpArrayLength, OpAget, OpAgetWide, OpAgetObject, OpAgetBoolean, OpAgetByte, OpAgetShort, OpNewArray:
				regs[int(prevInst.Operands[0])] = nonEmptyGibberishValue
			default:
				// skip this shit
				continue
			}
		default:
			continue
		}
	}

	return regs
}

func (d *Dex) MethodArguments(signature *regexp.Regexp) ([][]string, error) {
	out := [][]string{}
	outCache := make(map[string]struct{})
	for _, method := range d.Methods {
		if err := method.ParseCode(); err != nil {
			return nil, fmt.Errorf("parse code: %w", err)
		}

		for i, instr := range method.Body {
			if instr.Type != TypeInvocation {
				continue
			}

			methodIdx := int(instr.Operands[len(instr.Operands)-1])
			sig, hasMethodIdx := d.MethodsByIndex[methodIdx]
			if sig == "" || !hasMethodIdx {
				continue
			}
			if !signature.MatchString(sig) {
				continue
			}

			instr := instr
			start := 1
			if instr.Opcode == OpInvokeStatic || instr.Opcode == OpInvokeStaticRange {
				start = 0
			}

			regs := make(map[int]string, len(instr.Operands)-1-start)
			for _, reg := range instr.Operands[start : len(instr.Operands)-1] {
				regs[int(reg)] = "" // reg can't exceed byte
			}

			regValues := d.findFirstRegUpdateValue(&method, i, regs)
			regs = regValues
			allEmpty := true
			regStates := make([]string, 0, len(regs))
			for k := len(instr.Operands) - 2; k >= start; k-- {
				if regs[int(instr.Operands[k])] == nonEmptyGibberishValue {
					regs[int(instr.Operands[k])] = ""
				}
				regStates = append(regStates, regs[int(instr.Operands[k])])
				if regs[int(instr.Operands[k])] != "" {
					allEmpty = false
				}
			}
			if allEmpty {
				continue
			}

			cacheKey := strings.Join(regStates, "|")
			if _, ok := outCache[cacheKey]; ok {
				continue
			}

			outCache[cacheKey] = struct{}{}
			out = append(out, regStates)
		}
	}

	return out, nil
}

func (d *Dex) MethodStrings(m *Method) []string {
	var out []string
	for _, instr := range m.Body {
		if instr.Opcode == OpConstString {
			out = append(out, string(d.rawDex.StringDefs[instr.Operands[1]].Data))
		}
	}

	return out
}

func (d *Dex) GetConstStrings() []string {
	// NOTE: remove annotations from string set
	stringSet := make(map[string]struct{}, len(d.rawDex.StringDefs))
	for i, str := range d.rawDex.StringDefs {

		if _, ok := d.rawDex.AuxiliaryStrings[i]; ok {
			continue
		}

		stringSet[string(str.Data)] = struct{}{}
	}

	out := make([]string, 0, len(stringSet))
	for str := range stringSet {
		out = append(out, str)
	}

	return out
}
