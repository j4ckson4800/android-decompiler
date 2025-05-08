package smali

import (
	"bytes"
	"fmt"

	"github.com/j4ckson4800/android-decompiler/smali/internal"
)

type Method struct {
	Class      string
	Name       string
	ReturnType string

	rawMethod internal.Method
	Body      []Instruction
}

func NewMethod(cls, name, returnType string, m internal.Method) (Method, error) {

	return Method{
		Class:      cls,
		Name:       name,
		ReturnType: returnType,

		rawMethod: m,
	}, nil
}

func (m *Method) ParseCode() error {
	reader := bytes.NewReader(m.rawMethod.CodeItem.Payload)
	p := NewParser(reader)

	m.Body = make([]Instruction, 0, len(m.rawMethod.CodeItem.Payload)/(2*2)) // 2 bytes per word, instruction usually consist of 2 words
	end := len(m.rawMethod.CodeItem.Payload)
	globalOffset := 0
	for p.HasMore() {
		instr, err := p.ParseInstruction()
		if err != nil {
			return fmt.Errorf("read instruction: %w", err)
		}

		// Skip packed payloads which may be undefined but exist at the end of the function
		if instr.Opcode == OpNop && instr.Operands[0] != 0 {
			break
		}

		if instr.Opcode == OpFilledArrayData || instr.Opcode == OpPackedSwitch || instr.Opcode == OpSparseSwitch {
			// Skip filled array data, idc about it
			// TODO: add support for this
			offset := globalOffset + (int(reader.Size()) - reader.Len())
			globalOffset = offset
			payloadOffset := instr.Operands[len(instr.Operands)-1]*2 - 6

			if int64(end) > int64(offset)+payloadOffset {
				end = offset + int(payloadOffset)
			}

			reader = bytes.NewReader(m.rawMethod.CodeItem.Payload[offset:end])
			p = NewParser(reader)
		}

		m.Body = append(m.Body, instr)
	}

	return nil
}
