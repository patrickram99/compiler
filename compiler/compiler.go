package compiler

import (
	"fmt"
	"main/ast"
	"strings"
)

var (
	labelCount         int
	intRegisterCount   int
	floatRegisterCount int
)

func GenerateMIPS(node ast.Node) string {
	var output strings.Builder
	intRegisterCount = 0
	floatRegisterCount = 0

	writeLines(&output, []string{
		".data",
		"newline: .asciiz \"\\n\"",
		".text",
		".globl main",
		"main:",
	})

	generateNode(&output, node)

	writeLines(&output, []string{
		"li $v0, 10",
		"syscall",
	})

	return output.String()
}

func writeLines(output *strings.Builder, lines []string) {
	for _, line := range lines {
		output.WriteString(line + "\n")
	}
}

func writeLine(output *strings.Builder, line string) {
	output.WriteString(line + "\n")
}

func generateNode(output *strings.Builder, node ast.Node) (int, bool) {
	switch n := node.(type) {
	case *ast.Program:
		var lastReg int
		var isFloat bool
		for _, stmt := range n.Statements {
			lastReg, isFloat = generateNode(output, stmt)
		}
		return lastReg, isFloat
	case *ast.ExpressionStatement:
		return generateNode(output, n.Expression)
	case *ast.InfixExpression:
		return generateInfixExpression(output, n)
	case *ast.IntegerLiteral:
		reg := getNextIntRegister()
		writeLine(output, fmt.Sprintf("li $t%d, %d", reg, n.Value))
		return reg, false
	case *ast.FloatLiteral:
		reg := getNextFloatRegister()
		writeLine(output, fmt.Sprintf("li.s $f%d, %f", reg, n.Value))
		return reg, true
	case *ast.CallExpression:
		if ident, ok := n.Function.(*ast.Variable); ok && ident.Value == "SpeakNow" {
			return generateSpeakNow(output, n)
		}
	}
	return 0, false
}

func generateInfixExpression(output *strings.Builder, node *ast.InfixExpression) (int, bool) {
	leftReg, leftIsFloat := generateNode(output, node.Left)
	rightReg, rightIsFloat := generateNode(output, node.Right)

	isFloat := leftIsFloat || rightIsFloat

	if isFloat {
		if !leftIsFloat {
			// Convert left integer to float
			floatReg := getNextFloatRegister()
			writeLine(output, fmt.Sprintf("mtc1 $t%d, $f%d", leftReg, floatReg))
			writeLine(output, fmt.Sprintf("cvt.s.w $f%d, $f%d", floatReg, floatReg))
			leftReg = floatReg
		}
		if !rightIsFloat {
			// Convert right integer to float
			floatReg := getNextFloatRegister()
			writeLine(output, fmt.Sprintf("mtc1 $t%d, $f%d", rightReg, floatReg))
			writeLine(output, fmt.Sprintf("cvt.s.w $f%d, $f%d", floatReg, floatReg))
			rightReg = floatReg
		}

		resultReg := getNextFloatRegister()
		switch node.Operator {
		case "+":
			writeLine(output, fmt.Sprintf("add.s $f%d, $f%d, $f%d", resultReg, leftReg, rightReg))
		case "-":
			writeLine(output, fmt.Sprintf("sub.s $f%d, $f%d, $f%d", resultReg, leftReg, rightReg))
		case "*":
			writeLine(output, fmt.Sprintf("mul.s $f%d, $f%d, $f%d", resultReg, leftReg, rightReg))
		case "/":
			writeLine(output, fmt.Sprintf("div.s $f%d, $f%d, $f%d", resultReg, leftReg, rightReg))
		}
		return resultReg, true
	} else {
		resultReg := getNextIntRegister()
		switch node.Operator {
		case "+":
			writeLine(output, fmt.Sprintf("add $t%d, $t%d, $t%d", resultReg, leftReg, rightReg))
		case "-":
			writeLine(output, fmt.Sprintf("sub $t%d, $t%d, $t%d", resultReg, leftReg, rightReg))
		case "*":
			writeLine(output, fmt.Sprintf("mul $t%d, $t%d, $t%d", resultReg, leftReg, rightReg))
		case "/":
			writeLine(output, fmt.Sprintf("div $t%d, $t%d, $t%d", resultReg, leftReg, rightReg))
		}
		return resultReg, false
	}
}

func generateSpeakNow(output *strings.Builder, node *ast.CallExpression) (int, bool) {
	if len(node.Arguments) != 1 {
		fmt.Println("Error: SpeakNow expects exactly one argument")
		return 0, false
	}

	reg, isFloat := generateNode(output, node.Arguments[0])

	if isFloat {
		writeLines(output, []string{
			fmt.Sprintf("mov.s $f12, $f%d", reg),
			"li $v0, 2", // System call for print float
			"syscall",
		})
	} else {
		writeLines(output, []string{
			fmt.Sprintf("move $a0, $t%d", reg),
			"li $v0, 1", // System call for print integer
			"syscall",
		})
	}

	writeLines(output, []string{
		"li $v0, 4",
		"la $a0, newline",
		"syscall",
	})

	return reg, isFloat
}

func getNextIntRegister() int {
	reg := intRegisterCount
	intRegisterCount++
	return reg
}

func getNextFloatRegister() int {
	reg := floatRegisterCount
	floatRegisterCount++
	return reg
}

func getNextLabel() string {
	labelCount++
	return fmt.Sprintf("label_%d", labelCount)
}
