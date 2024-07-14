package object

var (
	VOID = Void{}
)

type Void struct{}

func (void *Void) Type() Type {
	return VOID_OBJ
}

func (void *Void) Inspect() string {
	return ""
}
