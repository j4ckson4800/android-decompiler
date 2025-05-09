package defs

type MutatorFunc func(field *ProtoField) bool

type Mutator struct {
	MethodName string
	Mutator    MutatorFunc
}

func NewMutatorList() []Mutator {
	return []Mutator{
		// firstly check for bytes, then for lists
		{
			MethodName: "byteAt",
			Mutator:    BytesMutator,
		},
		{
			MethodName: "addInt",
			Mutator:    IntListMutator,
		},
		{
			MethodName: "isModifiable",
			Mutator:    ListMutator,
		},
	}
}

func BytesMutator(field *ProtoField) bool {
	field.Type = "bytes"
	return true
}

func ListMutator(field *ProtoField) bool {
	field.Qualifier = "repeated"
	// let's assume we have string type, it should work most of the time
	field.Type = "string"
	return true
}

func IntListMutator(field *ProtoField) bool {
	ListMutator(field)
	field.Type = "int32"
	return true
}
