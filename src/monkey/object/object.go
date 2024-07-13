package object

type Type string

const (
	INTEGER  = "INTEGER"
	STRING   = "STRING"
	BOOLEAN  = "BOOLEAN"
	RETURN   = "RETURN"
	FUNCTION = "FUNCTION"
	ARRAY    = "ARRAY"
	BUILT_IN = "BUILT_IN"
	ERROR    = "ERROR"
	NULL_OBJ = "NULL_OBJ"
)

type Object interface {
	Type() Type
	Inspect() string
}
