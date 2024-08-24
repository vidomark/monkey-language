package compiler

type SymbolScope string

const (
	GlobalScope  SymbolScope = "GLOBAL"
	LocalScope   SymbolScope = "LOCAL"
	BuiltinScope SymbolScope = "BUILTIN"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	Enclosing      *SymbolTable
	store          map[string]Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	store := make(map[string]Symbol)
	return &SymbolTable{store: store, numDefinitions: 0}
}

func NewEnclosedSymbolTable(enclosing *SymbolTable) *SymbolTable {
	symbolTable := NewSymbolTable()
	symbolTable.Enclosing = enclosing
	return symbolTable
}

func (symbolTable *SymbolTable) Define(identifier string) Symbol {
	symbol := Symbol{Name: identifier, Scope: GlobalScope, Index: symbolTable.numDefinitions}
	if symbolTable.Enclosing == nil {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}
	symbolTable.store[identifier] = symbol
	symbolTable.numDefinitions++
	return symbol
}

func (symbolTable *SymbolTable) Resolve(identifier string) (Symbol, bool) {
	symbol, ok := symbolTable.store[identifier]
	if !ok && symbolTable.Enclosing != nil {
		return symbolTable.Enclosing.Resolve(identifier)
	}
	return symbol, ok
}

func (symbolTable *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	symbol := Symbol{Name: name, Index: index, Scope: BuiltinScope}
	symbolTable.store[name] = symbol
	return symbol
}
