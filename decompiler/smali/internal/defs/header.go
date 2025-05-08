package defs

import (
	"fmt"
)

type Table struct {
	Size   uint32
	Offset uint32
} // Size: 0x8

type DexHeader struct {
	Magic      uint64   // 0x0
	Checksum   uint32   // 0x8
	Signature  [20]byte // 0xc
	FileSize   uint32   // 0x20
	HeaderSize uint32   // 0x24
	EndianTag  uint32   // 0x28
	Links      Table    // 0x2c
	MapOff     uint32   // 0x34
	StringIds  Table    // 0x38
	TypeIds    Table    // 0x40
	ProtoIds   Table    // 0x48
	FieldIds   Table    // 0x50
	MethodIds  Table    // 0x58
	ClassDefs  Table    // 0x60
	Data       Table    // 0x68
} // Size: 0x70

const LEConstant = 0x12345678
const BEConstant = 0x78563412
const DexHeaderSize = 0x70
const Magic = 0x0000000A786564

func NewDexHeader(r Parser) (DexHeader, error) {
	hdr := DexHeader{}
	if err := r.ReadStruct(&hdr); err != nil {
		return DexHeader{}, fmt.Errorf("read struct: %w", err)
	}

	return hdr, nil
}
