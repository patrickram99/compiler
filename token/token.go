package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILEGAL"
	EOF     = "EOF"

	// Identificadores y primitivos
	ID    = "ID"
	INT   = "INT"
	FLOAT = "FLOAT"

	// Operadores
	ASSIGN  = "="
	PLUS    = "+"
	MINUS   = "-"
	EXCL    = "!"
	TIMES   = "*"
	DIVIDES = "/"
	COLON   = ":"

	LT = "<"
	GT = ">"

	EQ     = "=="
	NOT_EQ = "!="

	// Delimitadores
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Palabras clave
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"

	STRING = "STRING"
)

var palabras_reservadas = map[string]TokenType{
	"isme":      FUNCTION,
	"enchanted": LET,
	"SparksFly": TRUE,
	"BadBlood":  FALSE,
	"LoverEra":  IF,
	"RepEra":    ELSE,
	"hi":        RETURN,
}

func CheckIdentificador(identificador string) TokenType {
	if tok, ok := palabras_reservadas[identificador]; ok {
		return tok
	}
	return ID
}
