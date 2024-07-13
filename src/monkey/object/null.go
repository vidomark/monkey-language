package object

var (
	NULL = &Null{}
)

type Null struct{}

func (n *Null) Type() Type { return NULL_OBJ }

func (n *Null) Inspect() string { return "null" }
