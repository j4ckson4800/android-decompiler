package defs

const (
	ProtobufPackage     = "Lcom/google/protobuf/"
	ProtobufFieldNumber = "_FIELD_NUMBER"
)

var (
	JavaDefaultTypes = map[string]string{
		"I":                  "int32",
		"J":                  "int64",
		"F":                  "float",
		"D":                  "double",
		"B":                  "byte",
		"Z":                  "bool",
		"Ljava/lang/String;": "string",
	}
)

type ProtoField struct {
	Name      string
	Index     int
	Type      string
	Qualifier string
}

type ProtoOneof struct {
	Name   string
	Fields []*ProtoField
}

type ProtoEnum struct {
}

type ProtoMessage struct {
	Name        string
	IsGlobal    bool
	Fields      []*ProtoField
	OneOfs      []*ProtoOneof
	SubMessages []*ProtoMessage
	Enums       []ProtoEnum
}

type ProtoPackage struct {
	FileName      string
	PackageName   string
	GoPackageName string
	Messages      map[string]*ProtoMessage
	Enums         []ProtoEnum
	Imports       []string
}
