package ast

import (
	"main/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "enchanted"},
				Name: &Variable{
					Token: token.Token{Type: token.ID, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Variable{
					Token: token.Token{Type: token.ID, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "enchanted myVar = anotherVar;" {
		t.Errorf("program.String() genero un error en: %q", program.String())
	}
}
