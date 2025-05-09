package internal

type ResChunkType uint16

const (
	ResNullType       ResChunkType = 0x0000
	ResStringPoolType ResChunkType = 0x0001
	ResTableType      ResChunkType = 0x0002
	ResXMLType        ResChunkType = 0x0003

	ResXMLFirstChunkType     ResChunkType = 0x0100
	ResXMLStartNamespaceType ResChunkType = 0x0100
	ResXMLEndNamespaceType   ResChunkType = 0x0101
	ResXMLStartElementType   ResChunkType = 0x0102
	ResXMLEndElementType     ResChunkType = 0x0103
	ResXMLCdataType          ResChunkType = 0x0104

	ResXMLLastChunkType   ResChunkType = 0x017f
	ResXMLResourceMapType ResChunkType = 0x0180

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
	Utf8Flag StringPoolFlag = 1 << 8
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
	return s.Flags&uint32(Utf8Flag) != 0
}

type StringSpan struct {
	Name      uint32
	FirstChar uint32
	LastChar  uint32
}

type ResTable struct {
	// https://github.com/iBotPeaches/platform_frameworks_base/blob/main/libs/androidfw/include/androidfw/ResourceTypes.h#L1640
	PackageID        uint32
	PackageName      [256]byte
	TypeStringOffset uint32
	LastPublicType   uint32
	KeyStringOffset  uint32
	LastPublicKey    uint32
	TypeIDOffset     uint32
}
