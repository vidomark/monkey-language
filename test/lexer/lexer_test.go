package lexer

import (
	"testing"
	"writing-in-interpreter-in-go/src/monkey/lexer"
	"writing-in-interpreter-in-go/src/monkey/token"
)

func TestNextToken(test *testing.T) {
	tests := []*token.Token{
		token.NewMultiByteToken(token.LET, "let"),
		token.NewMultiByteToken(token.IDENTIFIER, "five"),
		token.NewToken(token.ASSIGN, '='),
		token.NewToken(token.INT, '5'),
		token.NewToken(token.SEMICOLON, ';'),
		token.NewMultiByteToken(token.LET, "let"),
		token.NewMultiByteToken(token.IDENTIFIER, "ten"),
		token.NewToken(token.ASSIGN, '='),
		token.NewMultiByteToken(token.INT, "10"),
		token.NewToken(token.SEMICOLON, ';'),
		token.NewMultiByteToken(token.LET, "let"),
		token.NewMultiByteToken(token.IDENTIFIER, "add"),
		token.NewToken(token.ASSIGN, '='),
		token.NewMultiByteToken(token.FUNCTION, "fn"),
		token.NewToken(token.LPAREN, '('),
		token.NewToken(token.IDENTIFIER, 'x'),
		token.NewMultiByteToken(token.COMMA, ","),
		token.NewToken(token.IDENTIFIER, 'y'),
		token.NewToken(token.RPAREN, ')'),
		token.NewToken(token.LBRACE, '{'),
		token.NewToken(token.IDENTIFIER, 'x'),
		token.NewToken(token.PLUS, '+'),
		token.NewToken(token.IDENTIFIER, 'y'),
		token.NewToken(token.SEMICOLON, ';'),
		token.NewToken(token.RBRACE, '}'),
		token.NewToken(token.SEMICOLON, ';'),
		token.NewMultiByteToken(token.LET, "let"),
		token.NewMultiByteToken(token.IDENTIFIER, "result"),
		token.NewToken(token.ASSIGN, '='),
		token.NewMultiByteToken(token.IDENTIFIER, "add"),
		token.NewToken(token.LPAREN, '('),
		token.NewMultiByteToken(token.IDENTIFIER, "five"),
		token.NewMultiByteToken(token.COMMA, ","),
		token.NewMultiByteToken(token.IDENTIFIER, "ten"),
		token.NewToken(token.RPAREN, ')'),
		token.NewToken(token.SEMICOLON, ';'),
		token.NewToken(token.BANG, '!'),
		token.NewToken(token.MINUS, '-'),
		token.NewToken(token.SLASH, '/'),
		token.NewToken(token.ASTERISK, '*'),
		token.NewToken(token.INT, '5'),
		token.NewToken(token.SEMICOLON, ';'),
		token.NewToken(token.INT, '5'),
		token.NewToken(token.LT, '<'),
		token.NewMultiByteToken(token.INT, "10"),
		token.NewToken(token.GT, '>'),
		token.NewToken(token.INT, '5'),
		token.NewToken(token.SEMICOLON, ';'),
		token.NewMultiByteToken(token.IF, "if"),
		token.NewToken(token.LPAREN, '('),
		token.NewToken(token.INT, '5'),
		token.NewToken(token.LT, '<'),
		token.NewMultiByteToken(token.INT, "10"),
		token.NewToken(token.RPAREN, ')'),
		token.NewToken(token.LBRACE, '{'),
		token.NewMultiByteToken(token.RETURN, "return"),
		token.NewMultiByteToken(token.TRUE, "true"),
		token.NewToken(token.SEMICOLON, ';'),
		token.NewToken(token.RBRACE, '}'),
		token.NewMultiByteToken(token.ELSE, "else"),
		token.NewToken(token.LBRACE, '{'),
		token.NewMultiByteToken(token.RETURN, "return"),
		token.NewMultiByteToken(token.FALSE, "false"),
		token.NewToken(token.SEMICOLON, ';'),
		token.NewToken(token.RBRACE, '}'),
		token.NewMultiByteToken(token.INT, "10"),
		token.NewMultiByteToken(token.EQ, "=="),
		token.NewMultiByteToken(token.INT, "10"),
		token.NewToken(token.SEMICOLON, ';'),
		token.NewMultiByteToken(token.INT, "10"),
		token.NewMultiByteToken(token.NOT_EQ, "!="),
		token.NewToken(token.INT, '9'),
		token.NewToken(token.SEMICOLON, ';'),
		token.NewMultiByteToken(token.STRING, "foobar"),
		token.NewMultiByteToken(token.STRING, "foo bar"),
		token.NewToken(token.LBRACKET, '['),
		token.NewToken(token.INT, '1'),
		token.NewToken(token.COMMA, ','),
		token.NewToken(token.INT, '2'),
		token.NewToken(token.RBRACKET, ']'),
		token.NewToken(token.SEMICOLON, ';'),
		token.NewToken(token.EOF, byte(0)),
	}

	tokens := lexer.New(inputString())

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
			10 != 9;
			"foobar"
			"foo bar"
			[1, 2];
`
}
