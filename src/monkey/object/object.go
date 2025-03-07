package object

type Type string

const (
	INTEGER           = "INTEGER"
	STRING            = "STRING"
	BOOLEAN           = "BOOLEAN"
	RETURN            = "RETURN"
	FUNCTION          = "FUNCTION"
	ARRAY             = "ARRAY"
	COMPILED_FUNCTION = "COMPILED_FUNCTION"
	CLOSURE           = "CLOSURE"
	BUILT_IN          = "BUILT_IN"
	QUOTE             = "QUOTE"
	MACRO             = "MACRO"
	ERROR             = "ERROR"
	NULL_OBJ          = "NULL_OBJ"
	VOID_OBJ          = "VOID"
)

type Object interface {
	Type() Type
	Inspect() string
}
