package defs

type Parser interface {
	ReadULEB128() (uint64, error)
	ReadSLEB128() (int64, error)
	ReadBytes(n int64) ([]byte, error)
	ReadUint64() (uint64, error)
	ReadUint32() (uint32, error)
	ReadUint16() (uint16, error)
	ReadStruct(any) error
}
