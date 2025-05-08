package smali

type Field struct {
	Name       string
	Type       string
	ClassName  string
	Descriptor string
	Value      int64 // NOTE: wrap in some sort of value type wrapping any
}
