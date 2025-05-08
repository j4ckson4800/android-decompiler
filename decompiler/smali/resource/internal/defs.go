package internal

type ResChunkType uint16

const (
	ResNullType       ResChunkType = 0x0000
	ResStringPoolType ResChunkType = 0x0001
	ResTableType      ResChunkType = 0x0002
	ResXmlType        ResChunkType = 0x0003

	ResXmlFirstChunkType     ResChunkType = 0x0100
	ResXmlStartNamespaceType ResChunkType = 0x0100
	ResXmlEndNamespaceType   ResChunkType = 0x0101
	ResXmlStartElementType   ResChunkType = 0x0102
	ResXmlEndElementType     ResChunkType = 0x0103
	ResXmlCdataType          ResChunkType = 0x0104

	ResXmlLastChunkType   ResChunkType = 0x017f
	ResXmlResourceMapType ResChunkType = 0x0180

	ResTablePackageType       ResChunkType = 0x0200 // 512
	ResTableTypeType          ResChunkType = 0x0201 // 513
	ResTableTypeSpecType      ResChunkType = 0x0202 // 514
	ResTableTypeLibrary       ResChunkType = 0x0203 // 515
	ResTableTypeOverlay       ResChunkType = 0x0204 // 516
	ResTableTypeOverlayPolicy ResChunkType = 0x0205 // 517
	ResTableTypeStagedAlias   ResChunkType = 0x0206 // 518
)

type StringPoolFlag uint32

const (
	UTF8_FLAG StringPoolFlag = (1 << 8)
)

type ResChunkHeader struct {
	Type       ResChunkType
	HeaderSize uint16
	Size       uint32
}

type ResTableHeader struct {
	PackageCount uint32
}

type RestStringPool struct {
	StringCount   uint32
	StyleCount    uint32
	Flags         uint32
	StringsOffset uint32
	StylesOffset  uint32
}

func (s *RestStringPool) IsUTF8() bool {
	return s.Flags&uint32(UTF8_FLAG) != 0
}

type StringSpan struct {
	Name      uint32
	FirstChar uint32
	LastChar  uint32
}

type ResTable struct {
	// https://github.com/iBotPeaches/platform_frameworks_base/blob/cd920c63bf716f17cc91d3e4ad914bfe0ba0871d/libs/androidfw/include/androidfw/ResourceTypes.h#L1640
	PackageID        uint32
	PackageName      [256]byte
	TypeStringOffset uint32
	LastPublicType   uint32
	KeyStringOffset  uint32
	LastPublicKey    uint32
	TypeIDOffset     uint32
}
