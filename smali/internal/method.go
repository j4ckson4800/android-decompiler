package internal

import (
	"fmt"

	"github.com/j4ckson4800/android-decompiler/smali/internal/defs"
)

type Method struct {
	IndexDiff   uint64
	AccessFlags uint64
	codeOffset  uint64
	CodeItem    defs.CodeItem
}

func NewMethod(p Parser) (Method, error) {
	indexDiff, err := p.ReadULEB128()
	if err != nil {
		return Method{}, fmt.Errorf("read uleb128: %w", err)
	}

	accessFlags, err := p.ReadULEB128()
	if err != nil {
		return Method{}, fmt.Errorf("read uleb128: %w", err)
	}

	codeOffset, err := p.ReadULEB128()
	if err != nil {
		return Method{}, fmt.Errorf("read uleb128: %w", err)
	}

	return Method{
		IndexDiff:   indexDiff,
		AccessFlags: accessFlags,
		codeOffset:  codeOffset,
	}, nil
}

func (m *Method) ParseCode(p Parser) error {
	if m.codeOffset == 0 {
		return nil
	}

	if err := p.SetCursorTo(int64(m.codeOffset)); err != nil {
		return fmt.Errorf("set cursor: %w", err)
	}

	codeItem, err := defs.NewCodeItem(p)
	if err != nil {
		return fmt.Errorf("new code item: %w", err)
	}

	m.CodeItem = codeItem
	return nil
}
