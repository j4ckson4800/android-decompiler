package defs

type codeItem struct {
	RegisterSize uint16
	InsSize      uint16
	OutsSize     uint16
	TriesSize    uint16
	DebugInfoOff uint32
	InsnsSize    uint32
}

type CodeItem struct {
	rawCodeItem codeItem
	Payload     []byte
	// NOTE: store exc handlers if we need it later
}

func NewCodeItem(r Parser) (CodeItem, error) {
	code := codeItem{}
	if err := r.ReadStruct(&code); err != nil {
		return CodeItem{}, err
	}

	data, err := r.ReadBytes(int64(code.InsnsSize * 2))
	if err != nil {
		return CodeItem{}, err
	}

	return CodeItem{
		rawCodeItem: code,
		Payload:     data,
	}, nil
}
