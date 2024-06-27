package lexer

import (
	"testing"

	"main/token"
)

func TestNextToken(t *testing.T) {
	input := `enchanted five = 5.5;
	enchanted ten = 10;
	enchanted add = isme(x, y) {
	x + y;
	};
	enchanted result = add(five, ten);
	!-/*5;
	5 < 10 > 5;
	LoverEra (5 < 10) {
		hi SparksFly;
	} RepEra {
		hi BadBlood;
	}
	10 == 10;
	10 != 9;
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "enchanted"},
		{token.ID, "five"},
		{token.ASSIGN, "="},
		{token.FLOAT, "5.5"},
		{token.SEMICOLON, ";"},
		{token.LET, "enchanted"},
		{token.ID, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "enchanted"},
		{token.ID, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "isme"},
		{token.LPAREN, "("},
		{token.ID, "x"},
		{token.COMMA, ","},
		{token.ID, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.ID, "x"},
		{token.PLUS, "+"},
		{token.ID, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "enchanted"},
		{token.ID, "result"},
		{token.ASSIGN, "="},
		{token.ID, "add"},
		{token.LPAREN, "("},
		{token.ID, "five"},
		{token.COMMA, ","},
		{token.ID, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EXCL, "!"},
		{token.MINUS, "-"},
		{token.DIVIDES, "/"},
		{token.TIMES, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "LoverEra"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "hi"},
		{token.TRUE, "SparksFly"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "RepEra"},
		{token.LBRACE, "{"},
		{token.RETURN, "hi"},
		{token.FALSE, "BadBlood"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - Tipo erroneo de token. Esperaba %q, obtuvo %q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Tipo erroneo de literal. Esperaba %q, obtuvo %q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
