package smali

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

var (
	ErrUnknownOperandType = errors.New("unknown operand type")
	ErrULEB128Overflow    = errors.New("ULEB128 overflow")
	ErrULEB128TooLong     = errors.New("ULEB128 too long")
	ErrInvalidType        = errors.New("invalid type")
)

type parser struct {
	r *bytes.Reader
}

func NewParser(r *bytes.Reader) *parser {
	return &parser{r: r}
}

func (p *parser) ParseInstruction() (Instruction, error) {
	rawOpcode, err := p.r.ReadByte()
	if err != nil {
		return Instruction{}, fmt.Errorf("read opcode: %w", err)
	}

	opcode := Opcode(rawOpcode)
	operandType := getInstructionOperandsType(opcode)

	operands, err := p.tryReadOperands(operandType)
	if err != nil {
		return Instruction{}, fmt.Errorf("read operands: %w", err)
	}

	return Instruction{
		Opcode:      opcode,
		Type:        getInstructionType(opcode),
		Operands:    operands,
		OperandType: operandType,
	}, nil
}

func (p *parser) ReadULEB128() (uint64, error) {
	var result uint64
	var shift uint
	const maxBytes = 10

	for range maxBytes {
		b, err := p.r.ReadByte()
		if err != nil {
			return 0, fmt.Errorf("read byte: %w", err)
		}
		if shift >= 64 && b != 0 {
			return 0, ErrULEB128Overflow
		}
		if shift == 63 && b > 1 {
			return 0, ErrULEB128Overflow
		}

		value := uint64(b & 0x7f)
		if shift >= 64 {
			if value != 0 {
				return 0, ErrULEB128Overflow
			}
			return result, nil
		}

		result |= value << shift
		if b&0x80 == 0 {
			return result, nil
		}

		shift += 7
	}

	return 0, ErrULEB128TooLong
}

func (p *parser) ReadSLEB128() (int64, error) {
	var result int64
	var shift uint
	var lastByte byte
	const maxBytes = 10

	for range maxBytes {
		b, err := p.r.ReadByte()
		if err != nil {
			return 0, fmt.Errorf("read byte: %w", err)
		}
		if shift >= 64 && b != 0 {
			return 0, ErrULEB128Overflow
		}
		if shift == 63 && b > 1 {
			return 0, ErrULEB128Overflow
		}
		lastByte = b

		result |= int64(b&0x7f) << shift
		shift += 7

		if b&0x80 == 0 {

			if shift < 64 && lastByte&0x40 != 0 {
				result |= -1 << shift
			}

			return result, nil
		}
	}

	if shift < 64 && lastByte&0x40 != 0 {
		result |= -1 << shift
	}

	return result, nil
}

func (p *parser) ReadBytes(n int64) ([]byte, error) {
	buf := make([]byte, n)
	if _, err := p.r.Read(buf); err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	return buf, nil
}

func (p *parser) ReadUint64() (uint64, error) {
	buf := make([]byte, 8)
	if _, err := p.r.Read(buf); err != nil {
		return 0, fmt.Errorf("read: %w", err)
	}

	return binary.LittleEndian.Uint64(buf), nil
}

func (p *parser) ReadUint32() (uint32, error) {
	buf := make([]byte, 4)
	if _, err := p.r.Read(buf); err != nil {
		return 0, fmt.Errorf("read: %w", err)
	}

	return binary.LittleEndian.Uint32(buf), nil
}

func (p *parser) ReadUint16() (uint16, error) {
	buf := make([]byte, 2)
	if _, err := p.r.Read(buf); err != nil {
		return 0, fmt.Errorf("read: %w", err)
	}

	return binary.LittleEndian.Uint16(buf), nil
}

func (p *parser) ReadByte() (byte, error) {
	b, err := p.r.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("read: %w", err)
	}

	return b, nil
}

func (p *parser) Pos() int64 {
	return p.r.Size() - int64(p.r.Len())
}

func (p *parser) SkipN(n int64) error {
	if _, err := p.r.Seek(n, io.SeekCurrent); err != nil {
		return fmt.Errorf("seek: %w", err)
	}
	return nil
}

func (p *parser) SetCursorTo(offset int64) error {
	if _, err := p.r.Seek(offset, io.SeekStart); err != nil {
		return fmt.Errorf("seek: %w", err)
	}
	return nil
}

func (p *parser) ReadStruct(out any) error {
	if err := binary.Read(p.r, binary.LittleEndian, out); err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (p *parser) HasMore() bool {
	_, err := p.r.ReadByte()
	if err == nil {
		_ = p.r.UnreadByte()
	}
	return err == nil
}

func (p *parser) read2ShortRegs() ([]int64, error) {
	operands, err := p.r.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	return []int64{int64(operands & 0x0f), int64(operands >> 4)}, nil
}

func (p *parser) readImm() (int64, error) {
	buf := make([]byte, 2)
	if _, err := p.r.Read(buf); err != nil {
		return 0, fmt.Errorf("read: %w", err)
	}

	return int64(binary.LittleEndian.Uint16(buf)), nil
}

func (p *parser) readImm32() (int64, error) {
	buf := make([]byte, 4)
	if _, err := p.r.Read(buf); err != nil {
		return 0, fmt.Errorf("read: %w", err)
	}

	return int64(binary.LittleEndian.Uint32(buf)), nil
}

func (p *parser) readImm64() (int64, error) {
	buf := make([]byte, 8)
	if _, err := p.r.Read(buf); err != nil {
		return 0, fmt.Errorf("read: %w", err)
	}

	return int64(binary.LittleEndian.Uint64(buf)), nil
}

func (p *parser) read2ShortRegsImm() ([]int64, error) {
	operands, err := p.r.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("read byte: %w", err)
	}

	imm, err := p.readImm()
	if err != nil {
		return nil, fmt.Errorf("read imm16: %w", err)
	}

	return []int64{int64(operands & 0x0f), int64(operands >> 4), imm}, nil
}

func (p *parser) read2Imm() ([]int64, error) {
	imm1, err := p.readImm()
	if err != nil {
		return nil, fmt.Errorf("read imm16: %w", err)
	}

	imm2, err := p.readImm()
	if err != nil {
		return nil, fmt.Errorf("read imm16: %w", err)
	}

	return []int64{imm1, imm2}, nil
}

func (p *parser) readReg() (int64, error) {
	reg, err := p.r.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("read byte: %w", err)
	}

	return int64(reg), nil
}

func (p *parser) read3Regs() ([]int64, error) {
	operands := make([]int64, 3)
	for i := range operands {
		reg, err := p.readReg()
		if err != nil {
			return nil, fmt.Errorf("read reg: %w", err)
		}

		operands[i] = reg
	}

	return operands, nil
}

func (p *parser) readRegImm(size int8) ([]int64, error) {
	reg, err := p.readReg()
	if err != nil {
		return nil, fmt.Errorf("read reg: %w", err)
	}

	var imm int64
	switch size {
	case 16:
		imm, err = p.readImm()
	case 32:
		imm, err = p.readImm32()
	case 64:
		imm, err = p.readImm64()
	}
	if err != nil {
		return nil, fmt.Errorf("read immn: %w", err)
	}

	return []int64{reg, imm}, nil
}

func (p *parser) readRegisterRange() ([]int64, error) {
	count, err := p.readReg()
	if err != nil {
		return nil, fmt.Errorf("read reg: %w", err)
	}

	typeID, err := p.readImm()
	if err != nil {
		return nil, fmt.Errorf("read imm16: %w", err)
	}

	firstReg, err := p.readImm()
	if err != nil {
		return nil, fmt.Errorf("read imm16: %w", err)
	}

	regs := make([]int64, count, count+1)
	for i := range regs {
		regs[i] = firstReg + int64(i)
	}

	return append(regs, typeID), nil
}

func (p *parser) readRegisterArray() ([]int64, error) {
	count, err := p.r.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("read reg: %w", err)
	}

	gReg := count & 0x0f

	count >>= 4
	typeID, err := p.readImm()
	if err != nil {
		return nil, fmt.Errorf("read imm16: %w", err)
	}

	regs := make([]int64, count, count+1)
	reg, err := p.readImm()
	if err != nil {
		return nil, fmt.Errorf("read reg: %w", err)
	}
	for i := range regs {
		if i == 4 {
			// even though len can be at most 7
			// according to docs 5 is the max?
			// ref(35c): https://source.android.com/docs/core/runtime/instruction-formats
			regs[i] = int64(gReg)
			break
		}

		regs[i] = reg & 0x0f
		reg >>= 4
	}

	return append(regs, typeID), nil
}

func (p *parser) readRegConstWide(size int) ([]int64, error) {
	reg, err := p.readReg()
	if err != nil {
		return nil, err
	}

	const signMask16 = 0x7f00
	signMask := int64(signMask16)

	var imm int64
	switch size {
	case 16:
		imm, err = p.readImm()
	case 32:
		signMask <<= 16
		imm, err = p.readImm32()
	}
	if err != nil {
		return nil, fmt.Errorf("read immn: %w", err)
	}

	const signExtensionMask = int64(0xffffffffff0000)

	sign := (imm & signMask) >> (size - 1)
	return []int64{reg, (signExtensionMask & (^(sign - 1))) | imm}, nil
}

func (p *parser) readRegConstHigh(size int) ([]int64, error) {
	reg, err := p.readReg()
	if err != nil {
		return nil, fmt.Errorf("read reg: %w", err)
	}

	imm, err := p.readImm()
	if err != nil {
		return nil, fmt.Errorf("read imm16: %w", err)
	}

	return []int64{reg, imm << (size - 16)}, nil
}

func (p *parser) tryReadOperands(opType OperandType) ([]int64, error) {
	switch opType {
	case OperandTypeNone, OperandTypeReg:
		reg, err := p.readReg()
		if err != nil {
			return nil, fmt.Errorf("read reg: %w", err)
		}
		return []int64{reg}, nil
	case OperandType2reg:
		operands, err := p.read2ShortRegs()
		if err != nil {
			return nil, fmt.Errorf("read2ShortRegs: %w", err)
		}
		return operands, nil
	case OperandTypeRegShort:
		operands, err := p.readRegImm(16)
		if err != nil {
			return nil, fmt.Errorf("readRegImm16: %w", err)
		}
		return operands, nil
	case OperandType2short:
		_, err := p.readReg()
		if err != nil {
			return nil, fmt.Errorf("read reg: %w", err)
		}
		operands, err := p.read2Imm()
		if err != nil {
			return nil, fmt.Errorf("read2Imm: %w", err)
		}
		return operands, nil
	case OperandTypeRegUint:
		operands, err := p.readRegImm(32)
		if err != nil {
			return nil, fmt.Errorf("readRegImm32: %w", err)
		}
		return operands, nil
	case OperandType2regShort:
		regs, err := p.read2ShortRegsImm()
		if err != nil {
			return nil, fmt.Errorf("read2ShortRegsImm: %w", err)
		}
		return regs, nil
	case OperandType3reg:
		operands, err := p.read3Regs()
		if err != nil {
			return nil, fmt.Errorf("read3Regs: %w", err)
		}
		return operands, nil
	case OperandTypeShort:
		_, err := p.readReg()
		if err != nil {
			return nil, fmt.Errorf("read reg: %w", err)
		}
		imm, err := p.readImm()
		if err != nil {
			return nil, fmt.Errorf("readImm: %w", err)
		}
		return []int64{imm}, nil
	case OperandTypeUint:
		_, err := p.readReg()
		if err != nil {
			return nil, fmt.Errorf("read reg: %w", err)
		}
		imm, err := p.readImm32()
		if err != nil {
			return nil, fmt.Errorf("readRegImm32: %w", err)
		}
		return []int64{imm}, nil
	case OperandTypeRegUlong:
		operands, err := p.readRegImm(64)
		if err != nil {
			return nil, fmt.Errorf("readRegImm64: %w", err)
		}
		return operands, nil
	case OperandRegisterArray:
		operands, err := p.readRegisterArray()
		if err != nil {
			return nil, fmt.Errorf("readRegisterArray: %w", err)
		}
		return operands, nil
	case OperandRegisterArrayRange:
		operands, err := p.readRegisterRange()
		if err != nil {
			return nil, fmt.Errorf("readRegisterRange: %w", err)
		}
		return operands, nil
	case OperandRegHigh32:
		operands, err := p.readRegConstHigh(32)
		if err != nil {
			return nil, fmt.Errorf("readRegConstHigh: %w", err)
		}
		return operands, nil
	case OperandRegHigh64:
		operands, err := p.readRegConstHigh(64)
		if err != nil {
			return nil, fmt.Errorf("readRegConstHigh: %w", err)
		}
		return operands, nil
	case OperandRegWide16:
		operands, err := p.readRegConstWide(16)
		if err != nil {
			return nil, fmt.Errorf("readRegConstWide16: %w", err)
		}
		return operands, nil
	case OperandRegWide32:
		operands, err := p.readRegConstWide(32)
		if err != nil {
			return nil, fmt.Errorf("readRegConstWide32: %w", err)
		}
		return operands, nil
	}

	return nil, ErrUnknownOperandType
}
