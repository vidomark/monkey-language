package compiler

type SymbolScope string

const (
	GlobalScope   SymbolScope = "GLOBAL"
	LocalScope    SymbolScope = "LOCAL"
	BuiltinScope  SymbolScope = "BUILTIN"
	FreeScope     SymbolScope = "FREE"
	FunctionScope SymbolScope = "FUNCTION"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	Enclosing      *SymbolTable
	FreeSymbols    []Symbol
	store          map[string]Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	store := make(map[string]Symbol)
	return &SymbolTable{FreeSymbols: []Symbol{}, store: store, numDefinitions: 0}
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

func (symbolTable *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	symbol := Symbol{Name: name, Index: index, Scope: BuiltinScope}
	symbolTable.store[name] = symbol
	return symbol
}

func (symbolTable *SymbolTable) DefineFree(symbol Symbol) Symbol {
	symbolTable.FreeSymbols = append(symbolTable.FreeSymbols, symbol)
	freeSymbol := Symbol{
		Name:  symbol.Name,
		Scope: FreeScope,
		Index: len(symbolTable.FreeSymbols) - 1,
	}
	return freeSymbol
}

func (symbolTable *SymbolTable) Resolve(identifier string) (Symbol, bool) {
	symbol, ok := symbolTable.store[identifier]
	if !ok && symbolTable.Enclosing != nil {
		obj, ok := symbolTable.Enclosing.Resolve(identifier)
		if !ok {
			return obj, ok
		}

		if obj.Scope == GlobalScope || obj.Scope == BuiltinScope {
			return obj, ok
		}

		return symbolTable.DefineFree(obj), true
	}
	return symbol, ok
}

func (symbolTable *SymbolTable) DefineFunctionName(identifier string) Symbol {
	symbol := Symbol{Name: identifier, Index: 0, Scope: FunctionScope}
	symbolTable.store[identifier] = symbol
	return symbol
}
