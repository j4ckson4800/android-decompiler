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
	StringIDs  Table    // 0x38
	TypeIDs    Table    // 0x40
	ProtoIDs   Table    // 0x48
	FieldIDs   Table    // 0x50
	MethodIDs  Table    // 0x58
	ClassDefs  Table    // 0x60
	Data       Table    // 0x68
} // Size: 0x70

const LEConstant = 0x12345678
const BEConstant = 0x78563412
const DexHeaderSize = 0x70
const Magic = 0x0000000A786564

func NewDexHeader(p Parser) (DexHeader, error) {
	hdr := DexHeader{}
	if err := p.ReadStruct(&hdr); err != nil {
		return DexHeader{}, fmt.Errorf("read struct: %w", err)
	}

	return hdr, nil
}
