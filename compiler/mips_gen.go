package compiler

import (
	"fmt"
	"main/ast"
	"strings"
)

type MIPSGenerator struct {
	code        strings.Builder
	labelCount  int
	stackOffset int
}

func NewMIPSGenerator() *MIPSGenerator {
	return &MIPSGenerator{
		labelCount:  0,
		stackOffset: 0,
	}
}

func (g *MIPSGenerator) Generate(node ast.Node) string {
	g.visitNode(node)
	return g.code.String()
}

func (g *MIPSGenerator) visitNode(node ast.Node) {
	switch n := node.(type) {
	case *ast.Program:
		g.emitHeader()
		for _, stmt := range n.Statements {
			g.visitNode(stmt)
		}
		g.emitFooter()
	case *ast.ExpressionStatement:
		g.visitNode(n.Expression)
		g.emitPrintInt() // Print the result of the expression
	case *ast.IntegerLiteral:
		g.emitLoadImmediate("$t0", n.Value)
	case *ast.InfixExpression:
		g.visitInfixExpression(n)
		// Add more cases for other AST nodes
	}
}

func (g *MIPSGenerator) visitInfixExpression(node *ast.InfixExpression) {
	if g.shouldParenthesize(node) {
		g.visitNode(node.Left)
		g.emit("move $s0, $t0") // Save left result to $s0

		g.visitNode(node.Right)
		g.emit("move $s1, $t0") // Save right result to $s1

		// Perform the operation
		switch node.Operator {
		case "+":
			g.emit("add $t0, $s0, $s1")
		case "-":
			g.emit("sub $t0, $s0, $s1")
		case "*":
			g.emit("mul $t0, $s0, $s1")
		case "/":
			g.emit("div $s0, $s1")
			g.emit("mflo $t0")
		}
	} else {
		g.visitNode(node.Left)
		g.emit("move $t1, $t0") // Save left result
		g.visitNode(node.Right)
		g.emit("move $t2, $t0") // Save right result

		// Perform the operation
		switch node.Operator {
		case "+":
			g.emit("add $t0, $t1, $t2")
		case "-":
			g.emit("sub $t0, $t1, $t2")
		case "*":
			g.emit("mul $t0, $t1, $t2")
		case "/":
			g.emit("div $t1, $t2")
			g.emit("mflo $t0")
		}
	}
}

func (g *MIPSGenerator) shouldParenthesize(node *ast.InfixExpression) bool {
	parentPrecedence := g.getOperatorPrecedence(node.Operator)

	if left, ok := node.Left.(*ast.InfixExpression); ok {
		if g.getOperatorPrecedence(left.Operator) < parentPrecedence {
			return true
		}
	}

	if right, ok := node.Right.(*ast.InfixExpression); ok {
		if g.getOperatorPrecedence(right.Operator) <= parentPrecedence {
			return true
		}
	}

	return false
}

func (g *MIPSGenerator) getOperatorPrecedence(operator string) int {
	switch operator {
	case "*", "/":
		return 3
	case "+", "-":
		return 2
	default:
		return 1
	}
}

func (g *MIPSGenerator) visitCallExpression(node *ast.CallExpression) {
	if function, ok := node.Function.(*ast.Variable); ok {
		if function.Value == "SpeakNow" {
			for _, arg := range node.Arguments {
				g.visitNode(arg)
				g.emitPrintInt()
			}
		}
		// Add support for other builtin functions here
	}
}

func (g *MIPSGenerator) emitHeader() {
	g.emit(".data")
	g.emit("newline: .asciiz \"\\n\"")
	g.emit(".text")
	g.emit(".globl main")
	g.emit("main:")
}

func (g *MIPSGenerator) emitFooter() {
	g.emit("li $v0, 10")
	g.emit("syscall")
}

func (g *MIPSGenerator) emit(instruction string) {
	g.code.WriteString(instruction + "\n")
}

func (g *MIPSGenerator) emitLoadImmediate(register string, value int64) {
	g.emit(fmt.Sprintf("li %s, %d", register, value))
}

func (g *MIPSGenerator) emitPrintInt() {
	g.emit("move $a0, $t0")
	g.emit("li $v0, 1")
	g.emit("syscall")
	g.emit("la $a0, newline")
	g.emit("li $v0, 4")
	g.emit("syscall")
}

// Add more helper methods for generating MIPS instructions
