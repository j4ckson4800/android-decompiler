package resource

import (
	"fmt"
	"unsafe"

	"github.com/j4ckson4800/android-decompiler/decompiler/smali/resource/internal"
)

type StringPool struct {
	rawPool      internal.RestStringPool
	stringOffset int64
	strings      []uint32
}

func NewStringPool(p internal.Parser) (StringPool, error) {
	pool := StringPool{}
	stringPoolOffset := p.Pos() - int64(unsafe.Sizeof(internal.ResChunkHeader{}))
	if err := p.ReadStruct(&pool.rawPool); err != nil {
		return pool, fmt.Errorf("read header: %w", err)
	}

	strIndices := make([]uint32, pool.rawPool.StringCount)
	if err := p.ReadStruct(&strIndices); err != nil {
		return pool, fmt.Errorf("read strings: %w", err)
	}

	pool.stringOffset = stringPoolOffset + int64(pool.rawPool.StringsOffset)
	pool.strings = strIndices
	return pool, nil
}

func (s *StringPool) GetString(p internal.Parser, index uint32) (string, error) {
	if index >= s.rawPool.StringCount {
		return "", nil
	}

	offset := s.stringOffset + int64(s.strings[index])
	if err := p.SetCursorTo(offset); err != nil {
		return "", fmt.Errorf("set cursor: %w", err)
	}

	lenValue, err := p.ReadUint16()
	if err != nil {
		return "", fmt.Errorf("read len: %w", err)
	}

	length := int(lenValue) & 0xff
	if lenValue&0x8000 != 0 {
		// skip long strings
		return "", nil
	}

	if !s.rawPool.IsUTF8() {
		out := make([]byte, 0, length)
		for range length {
			char, err := p.ReadByte()
			if err != nil {
				return "", fmt.Errorf("read byte: %w", err)
			}

			if _, err := p.ReadByte(); err != nil {
				return "", fmt.Errorf("read byte: %w", err)
			}
			out = append(out, char)
		}
		return string(out), nil
	}

	str, err := p.ReadBytes(int64(length))
	if err != nil {
		return "", fmt.Errorf("read bytes: %w", err)
	}

	return string(str), nil
}
