package lexer

import (
	. "testing"
	"writing-in-interpreter-in-go/src/monkey/lexer"
	. "writing-in-interpreter-in-go/src/monkey/token"
)

func TestNextToken(test *T) {
	tests := []Token{
		NewMultiByteToken(LET, "let"),
		NewMultiByteToken(IDENTIFIER, "five"),
		NewToken(ASSIGN, '='),
		NewToken(INT, '5'),
		NewToken(SEMICOLON, ';'),
		NewMultiByteToken(LET, "let"),
		NewMultiByteToken(IDENTIFIER, "ten"),
		NewToken(ASSIGN, '='),
		NewMultiByteToken(INT, "10"),
		NewToken(SEMICOLON, ';'),
		NewMultiByteToken(LET, "let"),
		NewMultiByteToken(IDENTIFIER, "add"),
		NewToken(ASSIGN, '='),
		NewMultiByteToken(FUNCTION, "fn"),
		NewToken(LPAREN, '('),
		NewToken(IDENTIFIER, 'x'),
		NewMultiByteToken(COMMA, ","),
		NewToken(IDENTIFIER, 'y'),
		NewToken(RPAREN, ')'),
		NewToken(LBRACE, '{'),
		NewToken(IDENTIFIER, 'x'),
		NewToken(PLUS, '+'),
		NewToken(IDENTIFIER, 'y'),
		NewToken(SEMICOLON, ';'),
		NewToken(RBRACE, '}'),
		NewToken(SEMICOLON, ';'),
		NewMultiByteToken(LET, "let"),
		NewMultiByteToken(IDENTIFIER, "result"),
		NewToken(ASSIGN, '='),
		NewMultiByteToken(IDENTIFIER, "add"),
		NewToken(LPAREN, '('),
		NewMultiByteToken(IDENTIFIER, "five"),
		NewMultiByteToken(COMMA, ","),
		NewMultiByteToken(IDENTIFIER, "ten"),
		NewToken(RPAREN, ')'),
		NewToken(SEMICOLON, ';'),
		NewToken(BANG, '!'),
		NewToken(MINUS, '-'),
		NewToken(SLASH, '/'),
		NewToken(ASTERISK, '*'),
		NewToken(INT, '5'),
		NewToken(SEMICOLON, ';'),
		NewToken(INT, '5'),
		NewToken(LT, '<'),
		NewMultiByteToken(INT, "10"),
		NewToken(GT, '>'),
		NewToken(INT, '5'),
		NewToken(SEMICOLON, ';'),
		NewMultiByteToken(IF, "if"),
		NewToken(LPAREN, '('),
		NewToken(INT, '5'),
		NewToken(LT, '<'),
		NewMultiByteToken(INT, "10"),
		NewToken(RPAREN, ')'),
		NewToken(LBRACE, '{'),
		NewMultiByteToken(RETURN, "return"),
		NewMultiByteToken(TRUE, "true"),
		NewToken(SEMICOLON, ';'),
		NewToken(RBRACE, '}'),
		NewMultiByteToken(ELSE, "else"),
		NewToken(LBRACE, '{'),
		NewMultiByteToken(RETURN, "return"),
		NewMultiByteToken(FALSE, "false"),
		NewToken(SEMICOLON, ';'),
		NewToken(RBRACE, '}'),
		NewMultiByteToken(INT, "10"),
		NewMultiByteToken(EQ, "=="),
		NewMultiByteToken(INT, "10"),
		NewToken(SEMICOLON, ';'),
		NewMultiByteToken(INT, "10"),
		NewMultiByteToken(NOT_EQ, "!="),
		NewToken(INT, '9'),
		NewToken(SEMICOLON, ';'),
		NewToken(EOF, byte(0)),
	}

	tokens := lexer.NewLexer(inputString())

	for index, expectedToken := range tests {
		currentToken := tokens.NextToken()
		if currentToken.Type != expectedToken.Type {
			test.Fatalf("tests[%d] - tokentype wrong. expected=%qgot=%q",
				index, expectedToken.Type, currentToken.Type)
		}

		if currentToken.Literal != expectedToken.Literal {
			test.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				index, expectedToken.Literal, currentToken.Literal)
		}
	}
}

func inputString() string {
	return `let five = 5;
			let ten = 10;

			let add = fn(x, y) {
			  x + y;
			};

			let result = add(five, ten);
			!-/*5;
			5 < 10 > 5;

			if (5 < 10) {
				return true;
			} else {
				return false;
			}

			10 == 10;
			10 != 9;`
}
