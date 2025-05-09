package app

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/j4ckson4800/android-decompiler/cmd/proto-gen/internal/app/defs"
	"github.com/j4ckson4800/android-decompiler/decompiler"
	"github.com/j4ckson4800/android-decompiler/decompiler/smali"
)

var (
	ErrInvalidApkFile         = errors.New("invalid APK file")
	ErrPredefinedMessageClass = errors.New("predefined message class provided")
	ErrInvalidSuperclass      = errors.New("invalid superclass provided")
	ErrEmptyMessageClass      = errors.New("empty message class provided")
)

type codeGenerator interface {
	WritePackage(pkg *defs.ProtoPackage) error
}

type Parser struct {
	apk *decompiler.Apk

	packages  map[string]*defs.ProtoPackage
	generator codeGenerator
}

func NewParser(apkFile string, generator codeGenerator) (*Parser, error) {

	if apkFile == "" {
		return nil, ErrInvalidApkFile
	}

	file, err := os.Open(apkFile)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("file stat: %w", err)
	}

	apk, err := decompiler.NewApk(file, stat.Size())
	if err != nil {
		return nil, fmt.Errorf("new apk: %w", err)
	}

	return &Parser{
		apk:       apk,
		packages:  make(map[string]*defs.ProtoPackage, 16),
		generator: generator,
	}, nil
}

func (p *Parser) Parse() error {
	for _, dex := range p.apk.Dexes {
		for _, cls := range dex.Classes {

			msg, err := p.parseClass(cls)
			if err != nil {
				continue
			}

			packageName := p.getMessagePackageName(msg.Name)
			packageParts := strings.Split(packageName, ".")
			if _, ok := p.packages[packageName]; !ok {
				p.packages[packageName] = &defs.ProtoPackage{
					FileName:      packageParts[len(packageParts)-1],
					PackageName:   packageName,
					GoPackageName: strings.ReplaceAll(packageName, ".", "/"),
					Messages:      make(map[string]*defs.ProtoMessage, 16),
				}
			}

			p.packages[packageName].Messages[msg.Name] = msg
		}
	}

	p.reorganizeAndFixPackages()
	p.reorganizeFields()

	return nil
}

func (p *Parser) GenerateProtoDefs() error {
	for _, pkg := range p.packages {
		if err := p.generator.WritePackage(pkg); err != nil {
			return fmt.Errorf("write package: %w", err)
		}
	}
	return nil
}

func (p *Parser) getMessagePackageName(typename string) string {
	packageParts := strings.Split(typename, "/")
	lastPart := packageParts[len(packageParts)-1]
	fileName := strings.Split(lastPart, "$")[0]
	packageParts[len(packageParts)-1] = fileName

	return strings.Join(packageParts, ".")[1:]
}

func (p *Parser) getMessageName(typename string) string {
	packageName := p.getMessagePackageName(typename)
	packageNameParts := strings.Split(packageName, ".")
	sanitizedTypeName := strings.TrimPrefix(typename, "L"+strings.ReplaceAll(packageName, ".", "/")+"$")

	return strings.TrimSuffix(packageNameParts[len(packageNameParts)-1]+"."+strings.ReplaceAll(sanitizedTypeName, "$", "."), ";")
}

func (p *Parser) parseClass(cls smali.Class) (*defs.ProtoMessage, error) {
	if strings.HasPrefix(cls.Name, defs.ProtobufPackage) {
		return nil, ErrPredefinedMessageClass
	}

	if !strings.HasPrefix(cls.SuperClass, defs.ProtobufPackage) {
		return nil, ErrInvalidSuperclass
	}

	if len(cls.StaticFields) < 2 {
		return nil, ErrEmptyMessageClass
	}

	msg := defs.ProtoMessage{
		Name:     cls.Name,
		IsGlobal: true,
		Fields:   make([]*defs.ProtoField, 0, len(cls.StaticFields)-2), // Omit default instance and parser
	}

	for _, staticField := range cls.StaticFields {
		if !strings.HasSuffix(staticField.Name, defs.ProtobufFieldNumber) {
			continue
		}

		fieldName := strings.TrimSuffix(staticField.Name, defs.ProtobufFieldNumber)
		msg.Fields = append(
			msg.Fields, &defs.ProtoField{
				Name:  strings.ToLower(fieldName),
				Index: int(staticField.Value),
			},
		)
	}

	for _, instanceField := range cls.InstanceFields {
		for _, fieldDef := range msg.Fields {
			if strings.ReplaceAll(fieldDef.Name, "_", "")+"_" == strings.ToLower(instanceField.Name) {
				fieldDef.Type = instanceField.Type

				if defType, ok := defs.JavaDefaultTypes[instanceField.Type]; ok {
					fieldDef.Type = defType
				}
				break
			}
		}
	}

	if err := p.fixOneofTypedFields(cls, &msg); err != nil {
		return nil, fmt.Errorf("fix oneof typed fields: %w", err)
	}

	return &msg, nil
}

func (p *Parser) makeSnakeCase(name string) string {
	snakeCase := make([]rune, 0, len(name)+len(name)/2)
	for _, char := range name {

		if char == '_' {
			break
		}

		if char >= 'A' && char <= 'Z' {
			if len(snakeCase) > 0 {
				snakeCase = append(snakeCase, '_')
			}
			snakeCase = append(snakeCase, char+32)
		} else {
			snakeCase = append(snakeCase, char)
		}

	}
	return string(snakeCase)
}

func (p *Parser) fixOneofTypedFields(cls smali.Class, msg *defs.ProtoMessage) error {

	setters := make(map[int]smali.Method, len(cls.Methods)/2)
	for _, method := range cls.Methods {
		if !strings.HasPrefix(method.Name, "set") {
			continue
		}

		untypedFieldIdx := slices.IndexFunc(
			msg.Fields, func(field *defs.ProtoField) bool {
				return field.Type == "" && strings.ReplaceAll(field.Name, "_", "") == strings.ToLower(method.Name[3:])
			},
		)
		if untypedFieldIdx == -1 {
			continue
		}

		setters[untypedFieldIdx] = method
	}

	if len(setters) == 0 {
		return nil
	}

	newFields := make([]*defs.ProtoField, 0, len(msg.Fields))
	oneOfs := make(map[int]*defs.ProtoOneof, len(setters)/2)
	for idx := range msg.Fields {
		if _, ok := setters[idx]; ok {
			continue
		}
		newFields = append(newFields, msg.Fields[idx])
	}

	for fieldIdx, setter := range setters {
		if err := setter.ParseCode(); err != nil {
			return fmt.Errorf("parse setter code: %w", err)
		}

		field := msg.Fields[fieldIdx]
		for _, instr := range setter.Body {
			if instr.Type != smali.TypeInstanceOp {
				continue
			}
			fieldNameIdx := instr.Operands[len(instr.Operands)-1]

			clsFieldIdx := slices.IndexFunc(
				cls.InstanceFields, func(clsField smali.Field) bool {
					return clsField.DefIdx == int(fieldNameIdx)
				},
			)

			if clsFieldIdx == -1 {
				continue
			}

			clsField := cls.InstanceFields[clsFieldIdx]
			if strings.HasSuffix(clsField.Name, "Case_") {
				continue
			}

			if _, ok := oneOfs[clsFieldIdx]; !ok {
				oneOfs[clsFieldIdx] = &defs.ProtoOneof{
					Name: p.makeSnakeCase(clsField.Name),
				}
			}

			field.Type = setter.ArgumentsSignature
			oneOfs[clsFieldIdx].Fields = append(oneOfs[clsFieldIdx].Fields, field)
		}
	}

	msg.OneOfs = make([]*defs.ProtoOneof, 0, len(oneOfs))
	for _, oneOf := range oneOfs {
		if len(oneOf.Fields) == 0 {
			continue
		}

		msg.OneOfs = append(msg.OneOfs, oneOf)
	}
	msg.Fields = newFields

	return nil
}

func (p *Parser) guessInternalProtoType(messageName string, field *defs.ProtoField) bool {
	if !strings.HasPrefix(field.Type, "Lcom/google/protobuf/") {
		return true
	}

	for _, dex := range p.apk.Dexes {
		cls, ok := dex.Classes[field.Type]
		if !ok {
			continue
		}

		// first of all check for `bytes` since it's ByteBuffer, not List
		if slices.ContainsFunc(
			cls.Methods, func(method smali.Method) bool {
				return method.Name == "byteAt"
			},
		) {
			field.Type = "bytes"
			return true
		}

		if !slices.ContainsFunc(
			cls.Methods, func(method smali.Method) bool {
				return method.Name == "isModifiable"
			},
		) {
			continue
		}

		// We have list over here
		field.Qualifier = "repeated"
		field.Type = "string"

		// Try to find if it's list of ints
		if slices.ContainsFunc(
			cls.Methods, func(method smali.Method) bool {
				return method.Name == "addInt"
			},
		) {
			field.Type = "int32"
			return true
		}

		fmt.Printf("Field %s of message %s has unknown list type, `string` assumed\n", field.Name, messageName)
		return true
	}

	fmt.Printf("Field %s of message %s has no type, probably because of oneof or any abuse\n", field.Name, messageName)
	return false
}

func (p *Parser) fixFieldTypes(pkg *defs.ProtoPackage, msg *defs.ProtoMessage, fields []*defs.ProtoField) {
	for _, field := range fields {

		if field.Type == "" {
			continue
		}

		if !p.guessInternalProtoType(msg.Name, field) {
			continue
		}

		if field.Type[0] != 'L' || field.Type[len(field.Type)-1] != ';' {
			continue
		}
		packageName := p.getMessagePackageName(field.Type)
		field.Type = p.getMessageName(field.Type)

		if packageName != pkg.PackageName {
			packageNameParts := strings.Split(packageName, ".")
			pkg.Imports = append(pkg.Imports, packageNameParts[len(packageNameParts)-1])
			continue
		}
	}
}

func (p *Parser) reorganizeFields() {
	for _, pkg := range p.packages {
		for _, msg := range pkg.Messages {

			slices.SortFunc(
				msg.Fields, func(a, b *defs.ProtoField) int {
					return a.Index - b.Index
				},
			)

			for _, oneof := range msg.OneOfs {
				slices.SortFunc(
					oneof.Fields, func(a, b *defs.ProtoField) int {
						return a.Index - b.Index
					},
				)
				p.fixFieldTypes(pkg, msg, oneof.Fields)
			}

			p.fixFieldTypes(pkg, msg, msg.Fields)
		}
	}
}

func (p *Parser) reorganizeAndFixPackages() {
	for _, pkg := range p.packages {
		newPkgMessages := make(map[string]*defs.ProtoMessage, len(pkg.Messages))
		for _, msg := range pkg.Messages {
			if strings.Contains(msg.Name, "/") {
				// Trim package name
				msg.Name = strings.TrimPrefix(msg.Name, "L"+pkg.GoPackageName+"$")
			}

			if strings.Contains(msg.Name, "$") {
				msg.IsGlobal = false
			}

			messageHierarchy := strings.Split(strings.TrimSuffix(msg.Name, ";"), "$")
			var previousMessage *defs.ProtoMessage

			for partIdx := len(messageHierarchy); partIdx > 0; partIdx-- {
				pkgName := strings.Join(messageHierarchy[:partIdx], ".")
				part := pkg.FileName + "." + pkgName

				if _, ok := newPkgMessages[part]; ok {
					if previousMessage != nil {
						newPkgMessages[part].SubMessages = append(newPkgMessages[part].SubMessages, previousMessage)
					}
					break
				}

				newPkgMessages[part] = &defs.ProtoMessage{
					Name:     messageHierarchy[partIdx-1],
					IsGlobal: partIdx == 1,
				}

				if previousMessage != nil {
					newPkgMessages[part].SubMessages = append(newPkgMessages[part].SubMessages, previousMessage)
				}
				previousMessage = newPkgMessages[part]
			}

			curMsg := newPkgMessages[pkg.FileName+"."+strings.Join(messageHierarchy, ".")]
			curMsg.Fields = msg.Fields
			curMsg.Enums = msg.Enums
			curMsg.OneOfs = msg.OneOfs
		}
		pkg.Messages = newPkgMessages
	}
}
