package lexer

import (
	"unicode"
	"writing-in-interpreter-in-go/src/monkey/token"
)

type Lexer struct {
	input           string
	currentPosition int
	readPosition    int
	character       byte
}

func (lexer *Lexer) readCharacter() {
	if lexer.readPosition >= len(lexer.input) {
		lexer.character = 0
	} else {
		lexer.character = lexer.input[lexer.readPosition]
	}
	lexer.currentPosition = lexer.readPosition
	lexer.readPosition++
}

func (lexer *Lexer) Peek(count int) byte {
	if lexer.currentPosition+count > len(lexer.input) {
		return 0
	}
	return lexer.input[lexer.currentPosition+count]
}

func (lexer *Lexer) NextToken() *token.Token {
	var nextToken *token.Token
	lexer.skipWhitespace()

	switch lexer.character {
	case '=':
		nextCharacter := lexer.Peek(1)
		if nextCharacter == '=' {
			return newTwoCharacterToken(lexer, token.EQ)
		} else {
			nextToken = token.NewToken(token.ASSIGN, lexer.character)
		}
	case ';':
		nextToken = token.NewToken(token.SEMICOLON, lexer.character)
	case '(':
		nextToken = token.NewToken(token.LPAREN, lexer.character)
	case ')':
		nextToken = token.NewToken(token.RPAREN, lexer.character)
	case ',':
		nextToken = token.NewToken(token.COMMA, lexer.character)
	case '+':
		nextToken = token.NewToken(token.PLUS, lexer.character)
	case '{':
		nextToken = token.NewToken(token.LBRACE, lexer.character)
	case '}':
		nextToken = token.NewToken(token.RBRACE, lexer.character)
	case '[':
		nextToken = token.NewToken(token.LBRACKET, lexer.character)
	case ']':
		nextToken = token.NewToken(token.RBRACKET, lexer.character)
	case '!':
		if lexer.Peek(1) == '=' {
			return newTwoCharacterToken(lexer, token.NOT_EQ)
		} else {
			nextToken = token.NewToken(token.BANG, lexer.character)
		}
	case '-':
		nextToken = token.NewToken(token.MINUS, lexer.character)
	case '/':
		nextToken = token.NewToken(token.SLASH, lexer.character)
	case '*':
		nextToken = token.NewToken(token.ASTERISK, lexer.character)
	case '<':
		nextToken = token.NewToken(token.LT, lexer.character)
	case '>':
		nextToken = token.NewToken(token.GT, lexer.character)
	case '"':
		nextToken = lexer.newStringLiteral()
	case 0:
		nextToken = token.NewToken(token.EOF, byte(0))
	default:
		if isLetter(lexer.character) {
			return newLanguageElement(lexer)
		}
		if unicode.IsDigit(rune(lexer.character)) {
			return newNumber(lexer)
		}
		nextToken = token.NewToken(token.ILLEGAL, lexer.character)
	}

	lexer.readCharacter()

	return nextToken
}

func (lexer *Lexer) newStringLiteral() *token.Token {
	lexer.readCharacter()
	beforeReadPosition := lexer.currentPosition
	for lexer.character != '"' {
		lexer.readCharacter()
	}
	return token.NewMultiByteToken(token.STRING, lexer.input[beforeReadPosition:lexer.currentPosition])
}

func newTwoCharacterToken(lexer *Lexer, tokenType token.Type) *token.Token {
	beforeReadPosition := lexer.currentPosition
	lexer.readCharacter()
	lexer.readCharacter()
	return token.NewMultiByteToken(tokenType, lexer.input[beforeReadPosition:lexer.currentPosition])
}

func newNumber(lexer *Lexer) *token.Token {
	beforeReadPosition := lexer.currentPosition
	for unicode.IsDigit(rune(lexer.character)) {
		lexer.readCharacter()
	}

	return token.NewMultiByteToken(token.INT, lexer.input[beforeReadPosition:lexer.currentPosition])
}

func newLanguageElement(lexer *Lexer) *token.Token {
	identifier := lexer.readIdentifier()
	tokenType := token.LookUpTokenType(identifier)
	return token.NewMultiByteToken(tokenType, identifier)
}

func (lexer *Lexer) readIdentifier() string {
	var positionBefore = lexer.currentPosition
	for isLetter(lexer.character) {
		lexer.readCharacter()
	}
	return lexer.input[positionBefore:lexer.currentPosition]
}

func (lexer *Lexer) skipWhitespace() {
	for lexer.character == ' ' || lexer.character == '\t' || lexer.character == '\n' || lexer.character == '\r' {
		lexer.readCharacter()
	}
}

func New(input string) *Lexer {
	lexer := &Lexer{input: input, readPosition: 0}
	lexer.readCharacter()
	return lexer
}

func isLetter(character byte) bool {
	return unicode.IsLetter(rune(character)) || character == '_'
}
