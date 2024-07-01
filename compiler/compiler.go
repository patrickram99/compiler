package compiler

import (
	"main/ast"
)

func GenerateMIPS(node ast.Node) string {
	generator := NewMIPSGenerator()
	return generator.Generate(node)
}

// Remove or modify other functions that were used for interpretation
