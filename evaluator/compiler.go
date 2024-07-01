package evaluator

import (
	"fmt"
	"main/ast"
	"os"
	"strings"
)

func GenerateMIPS(node ast.Node) error {
	var code strings.Builder
	code.WriteString(".text\n")
	code.WriteString(".globl main\n")
	code.WriteString("main:\n")

	generateMIPSCode(node, &code)

	code.WriteString("    li $v0, 10\n")
	code.WriteString("    syscall\n")

	// Guardar el código generado en out.s
	err := os.WriteFile("out.s", []byte(code.String()), 0644)
	if err != nil {
		return fmt.Errorf("error al escribir el archivo out.s: %v", err)
	}

	return nil
}

func generateMIPSCode(node ast.Node, code *strings.Builder) {
	switch node := node.(type) {
	case *ast.Program:
		for _, statement := range node.Statements {
			generateMIPSCode(statement, code)
		}

	case *ast.ExpressionStatement:
		generateMIPSCode(node.Expression, code)

	case *ast.IntegerLiteral:
		code.WriteString(fmt.Sprintf("    li $t0, %d\n", node.Value))

	case *ast.InfixExpression:
		generateMIPSCode(node.Left, code)
		code.WriteString("    move $t1, $t0\n")
		generateMIPSCode(node.Right, code)
		code.WriteString("    move $t2, $t0\n")

		switch node.Operator {
		case "+":
			code.WriteString("    add $t0, $t1, $t2\n")
		case "-":
			code.WriteString("    sub $t0, $t1, $t2\n")
		case "*":
			code.WriteString("    mul $t0, $t1, $t2\n")
		case "/":
			code.WriteString("    div $t0, $t1, $t2\n")
		}

	case *ast.PrefixExpression:
		if node.Operator == "-" {
			generateMIPSCode(node.Right, code)
			code.WriteString("    neg $t0, $t0\n")
		}

	default:
		// Handle other node types or unsupported expressions
		code.WriteString(fmt.Sprintf("    # Unsupported node type: %T\n", node))
	}
}
