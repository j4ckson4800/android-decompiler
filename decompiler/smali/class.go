package smali

type Class struct {
	Name           string
	StaticFields   []Field
	InstanceFields []Field
	Methods        []Method
	SuperClass     string
}

func NewClass(name, superClass string) (Class, error) {
	return Class{
		Name:       name,
		SuperClass: superClass,
	}, nil
}
