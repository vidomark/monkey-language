package lexer

import (
	"unicode"
	. "writing-in-interpreter-in-go/src/monkey/token"
)

type Lexer struct {
	input           string
	currentPosition int
	readPosition    int
	character       byte
}

func (lexer *Lexer) readCharacter() byte {
	if lexer.readPosition+1 > len(lexer.input) {
		lexer.character = 0
		return lexer.character
	}

	lexer.character = lexer.input[lexer.readPosition]
	lexer.currentPosition = lexer.readPosition
	lexer.readPosition += 1
	return lexer.character
}

func (lexer *Lexer) Peek(count int) byte {
	if lexer.readPosition+count > len(lexer.input) {
		return 0
	}
	return lexer.input[lexer.currentPosition+count]
}

func (lexer *Lexer) NextToken() Token {
	var token Token
	lexer.skipWhitespace()

	switch lexer.character {
	case '=':
		nextCharacter := lexer.Peek(1)
		if nextCharacter == '=' {
			return newTwoCharacterToken(lexer, EQ)
		} else {
			token = NewToken(ASSIGN, lexer.character)
		}
	case ';':
		token = NewToken(SEMICOLON, lexer.character)
	case '(':
		token = NewToken(LPAREN, lexer.character)
	case ')':
		token = NewToken(RPAREN, lexer.character)
	case ',':
		token = NewToken(COMMA, lexer.character)
	case '+':
		token = NewToken(PLUS, lexer.character)
	case '{':
		token = NewToken(LBRACE, lexer.character)
	case '}':
		token = NewToken(RBRACE, lexer.character)
	case '!':
		if lexer.Peek(1) == '=' {
			return newTwoCharacterToken(lexer, NOT_EQ)
		} else {
			token = NewToken(BANG, lexer.character)
		}
	case '-':
		token = NewToken(MINUS, lexer.character)
	case '/':
		token = NewToken(SLASH, lexer.character)
	case '*':
		token = NewToken(ASTERISK, lexer.character)
	case '<':
		token = NewToken(LT, lexer.character)
	case '>':
		token = NewToken(GT, lexer.character)
	case 0:
		token = NewToken(EOF, byte(0))
	default:
		if isLetter(lexer.character) {
			return newLanguageElement(lexer)
		}
		if unicode.IsDigit(rune(lexer.character)) {
			return newNumber(lexer)
		}
		token = NewToken(ILLEGAL, lexer.character)
	}

	lexer.readCharacter()
	return token
}

func newTwoCharacterToken(lexer *Lexer, tokenType Type) Token {
	beforeReadPosition := lexer.currentPosition
	lexer.readCharacter()
	lexer.readCharacter()
	return NewMultiByteToken(tokenType, lexer.input[beforeReadPosition:lexer.currentPosition])
}

func newNumber(lexer *Lexer) Token {
	beforeReadPosition := lexer.currentPosition
	for unicode.IsDigit(rune(lexer.character)) {
		lexer.readCharacter()
	}

	return NewMultiByteToken(INT, lexer.input[beforeReadPosition:lexer.currentPosition])
}

func newLanguageElement(lexer *Lexer) Token {
	identifier := lexer.readIdentifier()
	tokenType := LookUpTokenType(identifier)
	return NewMultiByteToken(tokenType, identifier)
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

func NewLexer(input string) *Lexer {
	lexer := &Lexer{input: input, readPosition: 0}
	lexer.readCharacter()
	return lexer
}

func isLetter(character byte) bool {
	return unicode.IsLetter(rune(character)) || character == '_'
}
