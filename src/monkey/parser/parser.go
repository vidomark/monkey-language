package parser

import (
	"fmt"
	"strconv"
	"writing-in-interpreter-in-go/src/monkey/ast"
	"writing-in-interpreter-in-go/src/monkey/lexer"
	"writing-in-interpreter-in-go/src/monkey/token"
)

type (
	prefixParseFns func() ast.Expression
	infixParseFns  func(expression ast.Expression) ast.Expression
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken *token.Token
	peekToken    *token.Token
	prefixFns    map[token.Type]prefixParseFns
	infixFns     map[token.Type]infixParseFns
	Errors       []string
}

func New(lexer *lexer.Lexer) *Parser {
	parser := Parser{
		lexer:     lexer,
		prefixFns: make(map[token.Type]prefixParseFns),
		infixFns:  make(map[token.Type]infixParseFns),
		Errors:    []string{},
	}

	parser.registerPrefix(token.IDENTIFIER, parser.parseIdentifier)
	parser.registerPrefix(token.INT, parser.parseIntegerLiteral)
	parser.registerPrefix(token.STRING, parser.parseStringLiteral)
	parser.registerPrefix(token.MINUS, parser.parsePrefixExpression)
	parser.registerPrefix(token.BANG, parser.parsePrefixExpression)
	parser.registerPrefix(token.TRUE, parser.parseBooleanExpression)
	parser.registerPrefix(token.FALSE, parser.parseBooleanExpression)
	parser.registerPrefix(token.LPAREN, parser.parseGroupedExpression)
	parser.registerPrefix(token.LBRACKET, parser.parseArrayLiteral)
	parser.registerPrefix(token.IF, parser.parseIfExpression)
	parser.registerPrefix(token.FUNCTION, parser.parseFunctionLiteral)
	parser.registerPrefix(token.MACRO, parser.parseMacroLiteral)

	parser.registerInfix(token.PLUS, parser.parseInfixExpression)
	parser.registerInfix(token.MINUS, parser.parseInfixExpression)
	parser.registerInfix(token.SLASH, parser.parseInfixExpression)
	parser.registerInfix(token.ASTERISK, parser.parseInfixExpression)
	parser.registerInfix(token.EQ, parser.parseInfixExpression)
	parser.registerInfix(token.NOT_EQ, parser.parseInfixExpression)
	parser.registerInfix(token.LT, parser.parseInfixExpression)
	parser.registerInfix(token.GT, parser.parseInfixExpression)
	parser.registerInfix(token.LPAREN, parser.parseCallExpression)
	parser.registerInfix(token.LBRACKET, parser.parseIndexExpression)

	parser.nextToken()
	parser.nextToken()
	return &parser
}

func (parser *Parser) ParseProgram() ast.Program {
	program := ast.Program{Statements: []ast.Statement{}}

	for parser.currentToken.Type != token.EOF {
		statement := parser.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, *statement)
		}
		parser.nextToken()
	}

	return program
}

func (parser *Parser) parseStatement() *ast.Statement {
	var statement ast.Statement

	switch parser.currentToken.Type {
	case token.LET:
		statement = parser.parseLetStatement()
	case token.RETURN:
		statement = parser.parseReturnStatement()
	default:
		statement = parser.parseExpressionStatement()
	}

	return &statement
}

func (parser *Parser) parseExpression(precedence int) ast.Expression {
	prefixFun := parser.prefixFns[parser.currentToken.Type]
	if prefixFun == nil {
		return nil
	}

	leftExpression := prefixFun()

	for !parser.peekTokenIs(token.SEMICOLON) && precedence < parser.peekPrecedence() {
		infix := parser.infixFns[parser.peekToken.Type]
		if infix == nil {
			return leftExpression
		}
		parser.nextToken()
		leftExpression = infix(leftExpression)
	}

	return leftExpression
}

func (parser *Parser) parseLetStatement() *ast.LetStatement {
	statement := ast.LetStatement{Token: parser.currentToken}
	if !parser.expectPeek(token.IDENTIFIER) {
		return nil
	}

	statement.Identifier = &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}

	if !parser.expectPeek(token.ASSIGN) {
		return nil
	}

	parser.nextToken()
	statement.Value = parser.parseExpression(ast.LOWEST)
	if !parser.expectPeek(token.SEMICOLON) {
		return &statement
	}
	return &statement
}

func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: parser.currentToken}
	parser.nextToken()
	stmt.ReturnValue = parser.parseExpression(ast.LOWEST)
	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}
	return stmt
}

func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := ast.ExpressionStatement{Token: parser.currentToken}
	statement.Expression = parser.parseExpression(ast.LOWEST)
	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}
	return &statement
}

func (parser *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}
}

func (parser *Parser) parseIntegerLiteral() ast.Expression {
	expression := ast.IntegerLiteral{
		Token: parser.currentToken,
	}

	value, err := strconv.ParseInt(parser.currentToken.Literal, 0, 64)
	if err != nil {
		return nil
	}

	expression.Value = value

	return &expression
}

func (parser *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: parser.currentToken, Value: parser.currentToken.Literal}
}

func (parser *Parser) parseBooleanExpression() ast.Expression {
	return &ast.BooleanExpression{Token: parser.currentToken, Value: parser.currentTokenIs(token.TRUE)}
}

func (parser *Parser) parsePrefixExpression() ast.Expression {
	prefixExpression := ast.PrefixExpression{
		Token:    parser.currentToken,
		Operator: parser.currentToken.Literal,
	}
	parser.nextToken()
	prefixExpression.Operand = parser.parseExpression(ast.PREFIX)
	return &prefixExpression
}

func (parser *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := ast.InfixExpression{
		Token:    parser.currentToken,
		Left:     left,
		Operator: parser.currentToken.Literal,
	}

	precedence := parser.currentPrecedence()
	parser.nextToken()
	expression.Right = parser.parseExpression(precedence)

	return &expression
}

func (parser *Parser) parseGroupedExpression() ast.Expression {
	parser.nextToken()
	expression := parser.parseExpression(ast.LOWEST)
	for !parser.expectPeek(token.RPAREN) {
		return nil
	}
	return expression
}

func (parser *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: parser.currentToken}
	array.Elements = parser.parseExpressionList(token.RBRACKET)
	return array
}

func (parser *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: parser.currentToken, Left: left}
	parser.nextToken()
	exp.Index = parser.parseExpression(ast.LOWEST)
	if !parser.expectPeek(token.RBRACKET) {
		return nil
	}
	return exp
}

func (parser *Parser) parseExpressionList(end token.Type) []ast.Expression {
	list := []ast.Expression{}
	if parser.peekTokenIs(end) {
		parser.nextToken()
		return list
	}
	parser.nextToken()
	list = append(list, parser.parseExpression(ast.LOWEST))
	for parser.peekTokenIs(token.COMMA) {
		parser.nextToken()
		parser.nextToken()
		list = append(list, parser.parseExpression(ast.LOWEST))
	}
	if !parser.expectPeek(end) {
		return nil
	}
	return list
}

func (parser *Parser) parseIfExpression() ast.Expression {
	expression := ast.IfExpression{Token: parser.currentToken}

	if !parser.expectPeek(token.LPAREN) {
		return nil
	}
	parser.nextToken()

	expression.Condition = parser.parseExpression(ast.LOWEST)

	if !parser.expectPeek(token.RPAREN) {
		return nil
	}

	if !parser.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = parser.parseBlockStatement()

	if parser.expectPeek(token.ELSE) {
		if !parser.expectPeek(token.LBRACE) {
			return nil
		}
		expression.Alternative = parser.parseBlockStatement()
	}

	return &expression
}

func (parser *Parser) parseFunctionLiteral() ast.Expression {
	expression := ast.FunctionLiteral{Token: *parser.currentToken}
	if !parser.expectPeek(token.LPAREN) {
		return nil
	}

	expression.Parameters = parser.parseFunctionParameters()

	if !parser.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Body = parser.parseBlockStatement()

	return &expression
}

func (parser *Parser) parseFunctionParameters() []ast.Expression {
	expressions := []ast.Expression{}

	if parser.peekTokenIs(token.RPAREN) {
		parser.nextToken()
		return expressions
	}

	parser.nextToken()
	expressions = append(expressions, parser.parseExpression(ast.LOWEST))

	if parser.peekTokenIs(token.RPAREN) {
		parser.nextToken()
		return expressions
	}

	for parser.peekTokenIs(token.COMMA) {
		parser.nextToken()
		parser.nextToken()
		expressions = append(expressions, parser.parseExpression(ast.LOWEST))
	}

	if !parser.expectPeek(token.RPAREN) {
		return nil
	}

	return expressions
}

func (parser *Parser) parseBlockStatement() *ast.BlockStatement {
	expression := ast.BlockStatement{Token: *parser.currentToken}
	parser.nextToken()
	for !parser.currentTokenIs(token.RBRACE) {
		expression.Statements = append(expression.Statements, *parser.parseStatement())
		parser.nextToken()
	}
	return &expression
}

func (parser *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expression := ast.CallExpression{
		Token:    *parser.currentToken,
		Function: function,
	}
	expression.Arguments = parser.parseExpressionList(token.RPAREN)
	return &expression
}

func (parser *Parser) parseCallArguments() []ast.Expression {
	arguments := []ast.Expression{}
	if parser.peekTokenIs(token.LPAREN) {
		parser.nextToken()
		return arguments
	}
	parser.nextToken()
	arguments = append(arguments, parser.parseExpression(ast.LOWEST))
	for parser.peekTokenIs(token.COMMA) {
		parser.nextToken()
		parser.nextToken()
		arguments = append(arguments, parser.parseExpression(ast.LOWEST))
	}
	if !parser.expectPeek(token.RPAREN) {
		return nil
	}
	return arguments
}

func (parser *Parser) parseMacroLiteral() ast.Expression {
	lit := &ast.MacroLiteral{Token: parser.currentToken}

	if !parser.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = parser.parseFunctionParameters()

	if !parser.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = parser.parseBlockStatement()

	return lit
}

func (parser *Parser) nextToken() {
	parser.currentToken = parser.peekToken
	parser.peekToken = parser.lexer.NextToken()
}

func (parser *Parser) expectPeek(tokenType token.Type) bool {
	if parser.peekTokenIs(tokenType) {
		parser.nextToken()
		return true
	}

	parser.peekError(tokenType)
	return false
}

func (parser *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, parser.peekToken.Type)
	parser.Errors = append(parser.Errors, msg)
}

func (parser *Parser) currentTokenIs(tokenType token.Type) bool {
	return parser.currentToken.Type == tokenType
}

func (parser *Parser) peekTokenIs(tokenType token.Type) bool {
	return parser.peekToken.Type == tokenType
}

func (parser *Parser) registerPrefix(tokenType token.Type, fn prefixParseFns) {
	parser.prefixFns[tokenType] = fn
}

func (parser *Parser) registerInfix(tokenType token.Type, fn infixParseFns) {
	parser.infixFns[tokenType] = fn
}

func (parser *Parser) currentPrecedence() int {
	precedence, ok := ast.Precedences[parser.currentToken.Type]
	if ok {
		return precedence
	}

	return ast.LOWEST
}

func (parser *Parser) peekPrecedence() int {
	precedence, ok := ast.Precedences[parser.peekToken.Type]
	if ok {
		return precedence
	}

	return ast.LOWEST
}
