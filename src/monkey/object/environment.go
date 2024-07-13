package object

type Environment struct {
	store     map[string]Object
	enclosing *Environment
}

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object), enclosing: nil}
}

func NewEnclosedEnvironment(enclosing *Environment) *Environment {
	environment := NewEnvironment()
	environment.enclosing = enclosing
	return environment
}

func (environment *Environment) Get(name string) (Object, bool) {
	obj, ok := environment.store[name]
	if !ok && environment.enclosing != nil {
		obj, ok = environment.enclosing.Get(name)
	}
	return obj, ok
}

func (environment *Environment) Set(name string, value Object) Object {
	environment.store[name] = value
	return value
}
