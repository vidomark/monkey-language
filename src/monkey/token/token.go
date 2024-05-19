package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENTIFIER = "IDENTIFIER" // add, foobar, x, y, ...
	INT        = "INT"        // 1343456

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	EQ     = "=="
	NOT_EQ = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

type Type string

type Token struct {
	Type    Type
	Literal string
}

func NewToken(tokenType Type, literal byte) Token {
	return Token{Type: tokenType, Literal: string(literal)}
}

func NewMultiByteToken(tokenType Type, literal string) Token {
	return Token{Type: tokenType, Literal: literal}
}

func LookUpTokenType(identifier string) Type {
	if tokenType, hasValue := keywords[identifier]; hasValue {
		return tokenType
	}

	return IDENTIFIER
}

var keywords = map[string]Type{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}
