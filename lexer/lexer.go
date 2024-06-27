package lexer

import "main/token"

type Lexer struct {
	input        string
	position     int  // Index char actual
	readPosition int  // Actual, luego de leer car
	ch           byte // Char actual
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.EXCL, l.ch)
		}
	case '/':
		tok = newToken(token.DIVIDES, l.ch)
	case '*':
		tok = newToken(token.TIMES, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if esLetra(l.ch) {
			tok.Literal = l.readIdentificador()
			tok.Type = token.CheckIdentificador(tok.Literal)
			return tok
		} else if esDigito(l.ch) {
			tok.Literal = l.readNumero()
			if l.ch == '.' {
				l.readChar()
				if esDigito(l.ch) {
					tok.Literal += "." + l.readNumero()
					tok.Type = token.FLOAT
				} else {
					tok = newToken(token.ILLEGAL, l.ch)
				}
			} else {
				tok.Type = token.INT
			}
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}

}

func (l *Lexer) readNumero() string {
	position := l.position

	for esDigito(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func esDigito(char byte) bool {
	return '0' <= char && char <= '9'
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentificador() string {
	position := l.position
	for esLetra(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func esLetra(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func newToken(tipoToken token.TokenType, caracter byte) token.Token {
	return token.Token{Type: tipoToken, Literal: string(caracter)}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII de nulo o fin de archivo
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}
