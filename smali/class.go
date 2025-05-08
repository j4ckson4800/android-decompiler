package smali

type Class struct {
	Name    string
	Fields  []Field
	Methods []Method
}

func NewClass(name string) (Class, error) {
	return Class{
		Name: name,
	}, nil
}
