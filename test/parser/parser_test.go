package parser

import (
	"fmt"
	"testing"
	"writing-in-interpreter-in-go/src/monkey/ast"
	"writing-in-interpreter-in-go/src/monkey/lexer"
	"writing-in-interpreter-in-go/src/monkey/parser"
	"writing-in-interpreter-in-go/src/monkey/token"
)

func TestLetStatements(testing *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkErrors(testing, p)

		if len(program.Statements) != 1 {
			testing.Fatalf("program.Statement does not contain 3 statements. Got=%d", len(program.Statements))
		}
		statement := program.Statements[0]
		if statement.TokenLiteral() != token.LET {
			testing.Fatalf("Statement is not a let statement. Got %s", statement.TokenLiteral())
		}

		letStatement, ok := statement.(*ast.LetStatement)
		if !ok {
			testing.Fatalf("Statement is not *ast.LetStatement. Got %T", statement)
		}

		if letStatement.Identifier.Value != test.expectedIdentifier {
			testing.Fatalf("Bad identifier. Got %s", letStatement.Identifier.Value)
		}

		if letStatement.Identifier.TokenLiteral() != test.expectedIdentifier {
			testing.Fatalf("Bad identifier. Got %s", letStatement.Identifier.Value)
		}
		testLiteralExpression(testing, letStatement.Value, test.expectedValue)
	}
}

func TestReturnStatement(testing *testing.T) {
	input := `
	return 5;
	return 10;
	return 993322;
	`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	expectedTokens := []struct {
		tokenType string
	}{
		{token.RETURN},
		{token.RETURN},
		{token.RETURN},
	}

	for index, expectedToken := range expectedTokens {
		statement := program.Statements[index]
		if statement.TokenLiteral() != expectedToken.tokenType {
			testing.Fatalf("Statement is not a let statement. Got %d", statement)
		}

		returnStatement, ok := statement.(*ast.ReturnStatement)
		if !ok {
			testing.Fatalf("Statement is not *ast.ReturnStatement. Got %T", returnStatement)
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar"
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	checkErrors(t, p)

	if !(len(program.Statements) == 1) {
		t.Errorf("Expected 1 statement, got %d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Statement is not *ast.Identifier. Got %T", statement)
	}

	testIdentifier(t, statement.Expression, input)
}

func TestNumberExpression(t *testing.T) {
	input := "5"
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	checkErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. Got=%d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is not *ast.ExpressionStatement. Got %T", statement)
	}

	expression, ok := statement.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Statement is not *ast.IntegerLiteral. Got %T", statement)
	}
	if expression.TokenLiteral() != "5" {
		t.Errorf("Expecting token e %s. Got %s", token.IDENTIFIER, expression.TokenLiteral())
	}

	if expression.Value != 5 {
		t.Errorf("Expecting 5. Got %d", expression.Value)
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkErrors(t, p)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}
	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}
	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkErrors(t, p)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}
	testIdentifier(t, indexExp.Left, "myArray")
	testInfixExpression(t, indexExp.Index, 1, "+", 1)
}

func TestPrefixExpression(t *testing.T) {
	input := []struct {
		input    string
		operator string
		value    int64
	}{
		{"-5", "-", 5},
		{"!15", "!", 15},
	}

	for _, element := range input {
		l := lexer.New(element.input)
		p := parser.New(l)
		program := p.ParseProgram()

		checkErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. Got %d", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statement is not *ast.ExpressionStatement. Got %T", statement)
		}

		testPrefixExpression(t, statement.Expression, element.operator, element.value)
	}
}

func TestInfixExpression(t *testing.T) {
	input := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, element := range input {
		l := lexer.New(element.input)
		p := parser.New(l)
		program := p.ParseProgram()

		if !(len(program.Statements) == 1) {
			t.Fatalf("program.Statements does not contain 1 statement. Got %d", len(program.Statements))
		}

		expression, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statement is not *ast.ExpressionStatement. Got %T", expression)
		}

		testInfixExpression(t, expression.Expression, element.leftValue, element.operator, element.rightValue)
	}
}

func TestIfStatement(t *testing.T) {
	input := "if (x < y) { x }"

	l := lexer.New(input)
	p := parser.New(l)
	checkErrors(t, p)
	program := p.ParseProgram()
	statements := program.Statements

	if len(statements) != 1 {
		t.Errorf("Expected 1 statement. Got %d.", len(program.Statements))
	}

	expression, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is not *ast.ExpressionStatement. Got %T", expression)
	}

	ifExpression, ok := expression.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("Operand of not of type %T. Got %T.", ast.IfExpression{}, ifExpression)
	}

	testInfixExpression(t, ifExpression.Condition, "x", "<", "y")

	consequence := ifExpression.Consequence
	if consequence.Token.Type != token.LBRACE {
		t.Fatalf("Expected { but got %s", consequence.Token.Type)
	}

	expressionStatement, ok := consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement. Got %T.", expressionStatement)
	}
	testIdentifier(t, expressionStatement.Expression, "x")
}

func TestIfElseStatement(t *testing.T) {
	input := "if (x < y) { x } else { y }"

	l := lexer.New(input)
	p := parser.New(l)
	checkErrors(t, p)
	program := p.ParseProgram()
	statements := program.Statements

	if len(statements) != 1 {
		t.Errorf("Expected 1 statement. Got %d.", len(program.Statements))
	}

	expression, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is not *ast.ExpressionStatement. Got %T", expression)
	}

	ifExpression, ok := expression.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("Operand of not of type %T. Got %T.", ast.IfExpression{}, ifExpression)
	}

	testInfixExpression(t, ifExpression.Condition, "x", "<", "y")

	consequence := ifExpression.Consequence
	if consequence.Token.Type != token.LBRACE {
		t.Fatalf("Expected { but got %s", consequence.Token.Type)
	}
	expressionStatement, ok := consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement. Got %T.", expressionStatement)
	}
	testIdentifier(t, expressionStatement.Expression, "x")

	alternative := ifExpression.Alternative
	if alternative.Token.Type != token.LBRACE {
		t.Fatalf("Expected { but got %s", consequence.Token.Type)
	}
	alternativeStatement, ok := alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement. Got %T.", alternativeStatement)
	}
	testIdentifier(t, alternativeStatement.Expression, "y")
}

func TestFunctionLiteral(t *testing.T) {
	input := `fn(x, y) { x + y; }`
	l := lexer.New(input)
	p := parser.New(l)
	checkErrors(t, p)

	program := p.ParseProgram()

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement. Got: %d.", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement. Got %T.", statement)
	}

	functionLiteral, ok := statement.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("Expected FunctionLiteral. Got %T.", functionLiteral)
	}

	if functionLiteral.TokenLiteral() != token.FUNCTION {
		t.Fatalf("Expected token.FUNCTION. Got %s.", functionLiteral.TokenLiteral())
	}
	parameters := functionLiteral.Parameters
	if len(parameters) != 2 {
		t.Fatalf("Expected 2 parameters. Got %d.", len(parameters))
	}
	testIdentifier(t, parameters[0], "x")
	testIdentifier(t, parameters[1], "y")

	body := functionLiteral.Body
	if body.TokenLiteral() != token.LBRACE {
		t.Fatalf("Expected {. Got %s.", body.TokenLiteral())
	}
	if len(body.Statements) != 1 {
		t.Fatalf("Expected 1 statement. Got %d.", len(body.Statements))
	}

	expressionStatement, ok := body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement. Got %T.", expressionStatement)
	}
	testInfixExpression(t, expressionStatement.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkErrors(t, p)
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)
		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(function.Parameters))
		}
		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Operand is not ast.CallExpression. got=%T", stmt.Expression)
	}
	testIdentifier(t, exp.Function, "add")
	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}
	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestMacroLiteralParsing(t *testing.T) {
	input := `macro(x, y) { x + y; }`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	macro, ok := stmt.Expression.(*ast.MacroLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.MacroLiteral. got=%T",
			stmt.Expression)
	}

	if len(macro.Parameters) != 2 {
		t.Fatalf("macro literal parameters wrong. want 2, got=%d\n",
			len(macro.Parameters))
	}

	testLiteralExpression(t, macro.Parameters[0], "x")
	testLiteralExpression(t, macro.Parameters[1], "y")

	if len(macro.Body.Statements) != 1 {
		t.Fatalf("macro.Body.Statements has not 1 statements. got=%d\n",
			len(macro.Body.Statements))
	}

	bodyStmt, ok := macro.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("macro body stmt is not ast.ExpressionStatement. got=%T",
			macro.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		}, {
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func testIntegerLiteral(t *testing.T, expression ast.Expression, value int64) {
	integerLiteral, ok := expression.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("Operand is not ast.IntegerLiteral. Got %T", expression)
	}

	if integerLiteral.Value != value {
		t.Errorf("Expecting %d. Got %d", value, integerLiteral.Value)
	}

	if integerLiteral.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, expression.TokenLiteral())
	}
}

func testBooleanLiteral(t *testing.T, expression ast.Expression, value bool) {
	booleanExpression, ok := expression.(*ast.BooleanExpression)
	if !ok {
		t.Errorf("Operand is not ast.BooleanExpression. Got %T", expression)
	}

	if booleanExpression.Value != value {
		t.Errorf("Expecting %t. Got %t", value, booleanExpression.Value)
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
	}
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) {
	switch v := expected.(type) {
	case int:
		testIntegerLiteral(t, exp, int64(v))
	case int64:
		testIntegerLiteral(t, exp, v)
	case string:
		testIdentifier(t, exp, v)
	case bool:
		testBooleanLiteral(t, exp, v)
	}
}

func testPrefixExpression(
	t *testing.T,
	expression ast.Expression,
	operator string,
	right interface{},
) {
	prefixExpression, ok := expression.(*ast.PrefixExpression)
	if !ok {
		t.Errorf("Operand is not ast.PrefixExpression. Got %T", expression)
	}

	if prefixExpression.Operator != operator {
		t.Errorf("Operator not %s. got=%s", operator, prefixExpression.Operator)
	}

	testLiteralExpression(t, prefixExpression.Operand, right)
}

func testInfixExpression(
	t *testing.T,
	expression ast.Expression,
	left interface{},
	operator string,
	right interface{},
) {
	opExp, ok := expression.(*ast.InfixExpression)
	if !ok {
		t.Errorf("expression is not ast.OperatorExpression. got=%T(%s)", expression, expression)
	}

	testLiteralExpression(t, opExp.Left, left)

	if opExp.Operator != operator {
		t.Errorf("expression.Operator is not '%s'. got=%q", operator, opExp.Operator)
	}

	testLiteralExpression(t, opExp.Right, right)
}

func checkErrors(t *testing.T, parser *parser.Parser) {
	for _, message := range parser.Errors {
		t.Errorf("parser error: %q", message)
	}
}
