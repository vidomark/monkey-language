package evaluator

import (
	"writing-in-interpreter-in-go/src/monkey/ast"
	"writing-in-interpreter-in-go/src/monkey/object"
)

func ExpandMacros(node ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(node, func(node ast.Node) ast.Node {
		callExpression, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}

		macro, ok := isMacroCall(callExpression, env)
		if !ok {
			return node
		}

		args := quoteArgs(callExpression)
		evalEnv := extendMacroEnv(macro, args)

		evaluated := Eval(macro.Body, evalEnv)

		quote, ok := evaluated.(*object.Quote)
		if !ok {
			panic("we only support returning AST-nodes from macros")
		}

		return quote.Node
	})
}

func isMacroCall(
	exp *ast.CallExpression,
	env *object.Environment,
) (*object.Macro, bool) {
	identifier, ok := exp.Function.(*ast.Identifier)
	if !ok {
		return nil, false
	}

	obj, ok := env.Get(identifier.Value)
	if !ok {
		return nil, false
	}

	macro, ok := obj.(*object.Macro)
	if !ok {
		return nil, false
	}

	return macro, true
}

func quoteArgs(exp *ast.CallExpression) []*object.Quote {
	var args []*object.Quote

	for _, argument := range exp.Arguments {
		args = append(args, &object.Quote{Node: argument})
	}

	return args
}

func extendMacroEnv(
	macro *object.Macro,
	args []*object.Quote,
) *object.Environment {
	extended := object.NewEnclosedEnvironment(macro.Env)

	for paramIdx, param := range macro.Parameters {
		extended.Set(param.(*ast.Identifier).Value, args[paramIdx])
	}

	return extended
}
