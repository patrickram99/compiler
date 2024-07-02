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

var currentRegister int = 0

func getNextRegister() int {
	reg := currentRegister
	currentRegister = (currentRegister + 1) % 8 // Use only $t0 to $t7
	return reg
}

func resetRegisterAllocation() {
	currentRegister = 0
}

func GenerateMIPS(node ast.Node) string {
	var output strings.Builder
	intRegisterCount = 0
	floatRegisterCount = 0

	writeLines(&output, []string{
		".data",
		"newline: .asciiz \"\\n\"",
		"true_str: .asciiz \"true\"",
		"false_str: .asciiz \"false\"",
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
func generateNode(output *strings.Builder, node ast.Node) (int, string) {
	switch n := node.(type) {
	case *ast.Program:
		var lastReg int
		var lastType string
		for _, stmt := range n.Statements {
			lastReg, lastType = generateNode(output, stmt)
		}
		return lastReg, lastType
	case *ast.ExpressionStatement:
		return generateNode(output, n.Expression)
	case *ast.InfixExpression:
		return generateInfixExpression(output, n)
	case *ast.IntegerLiteral:
		reg := getNextIntRegister()
		writeLine(output, fmt.Sprintf("li $t%d, %d", reg, n.Value))
		return reg, "int"
	case *ast.FloatLiteral:
		reg := getNextFloatRegister()
		writeLine(output, fmt.Sprintf("li.s $f%d, %f", reg, n.Value))
		return reg, "float"
	case *ast.Boolean:
		reg := getNextIntRegister()
		if n.Value {
			writeLine(output, fmt.Sprintf("li $t%d, 1", reg))
		} else {
			writeLine(output, fmt.Sprintf("li $t%d, 0", reg))
		}
		return reg, "bool"
	case *ast.CallExpression:
		if ident, ok := n.Function.(*ast.Variable); ok && ident.Value == "SpeakNow" {
			return generateSpeakNow(output, n)
		}
	case *ast.IfExpression:
		return generateIfExpression(output, n)
	case *ast.BlockStatement:
		return generateBlockStatement(output, n)
	}
	return 0, ""
}

func generateInfixExpression(output *strings.Builder, node *ast.InfixExpression) (int, string) {
	leftReg, leftType := generateNode(output, node.Left)
	rightReg, rightType := generateNode(output, node.Right)

	if leftType == "float" || rightType == "float" {
		return generateFloatInfixExpression(output, node.Operator, leftReg, rightReg, leftType, rightType)
	} else if leftType == "bool" || rightType == "bool" {
		return generateBoolInfixExpression(output, node.Operator, leftReg, rightReg)
	} else {
		return generateIntInfixExpression(output, node.Operator, leftReg, rightReg)
	}
}

func generateFloatInfixExpression(output *strings.Builder, operator string, leftReg, rightReg int, leftType, rightType string) (int, string) {
	if leftType != "float" {
		// Convert left integer to float
		floatReg := getNextFloatRegister()
		writeLine(output, fmt.Sprintf("mtc1 $t%d, $f%d", leftReg, floatReg))
		writeLine(output, fmt.Sprintf("cvt.s.w $f%d, $f%d", floatReg, floatReg))
		leftReg = floatReg
	}
	if rightType != "float" {
		// Convert right integer to float
		floatReg := getNextFloatRegister()
		writeLine(output, fmt.Sprintf("mtc1 $t%d, $f%d", rightReg, floatReg))
		writeLine(output, fmt.Sprintf("cvt.s.w $f%d, $f%d", floatReg, floatReg))
		rightReg = floatReg
	}

	resultReg := getNextFloatRegister()
	switch operator {
	case "+":
		writeLine(output, fmt.Sprintf("add.s $f%d, $f%d, $f%d", resultReg, leftReg, rightReg))
	case "-":
		writeLine(output, fmt.Sprintf("sub.s $f%d, $f%d, $f%d", resultReg, leftReg, rightReg))
	case "*":
		writeLine(output, fmt.Sprintf("mul.s $f%d, $f%d, $f%d", resultReg, leftReg, rightReg))
	case "/":
		writeLine(output, fmt.Sprintf("div.s $f%d, $f%d, $f%d", resultReg, leftReg, rightReg))
	case "<", ">", "<=", ">=", "==", "!=":
		intResultReg := getNextIntRegister()
		writeLine(output, fmt.Sprintf("c.%s.s $f%d, $f%d", floatComparisonOp(operator), leftReg, rightReg))
		writeLine(output, fmt.Sprintf("li $t%d, 1", intResultReg))
		writeLine(output, fmt.Sprintf("bc1t float_true_%d", labelCount))
		writeLine(output, fmt.Sprintf("li $t%d, 0", intResultReg))
		writeLine(output, fmt.Sprintf("float_true_%d:", labelCount))
		labelCount++
		return intResultReg, "bool"
	}
	return resultReg, "float"
}

func generateIntInfixExpression(output *strings.Builder, operator string, leftReg, rightReg int) (int, string) {
	resultReg := getNextRegister()
	switch operator {
	case "+":
		writeLine(output, fmt.Sprintf("add $t%d, $t%d, $t%d", resultReg, leftReg, rightReg))
	case "-":
		writeLine(output, fmt.Sprintf("sub $t%d, $t%d, $t%d", resultReg, leftReg, rightReg))
	case "*":
		writeLine(output, fmt.Sprintf("mul $t%d, $t%d, $t%d", resultReg, leftReg, rightReg))
	case "/":
		writeLine(output, fmt.Sprintf("div $t%d, $t%d", leftReg, rightReg))
		writeLine(output, fmt.Sprintf("mflo $t%d", resultReg))
	case "==":
		writeLine(output, fmt.Sprintf("seq $t%d, $t%d, $t%d", resultReg, leftReg, rightReg))
	case "!=":
		writeLine(output, fmt.Sprintf("sne $t%d, $t%d, $t%d", resultReg, leftReg, rightReg))
	case "<":
		writeLine(output, fmt.Sprintf("slt $t%d, $t%d, $t%d", resultReg, leftReg, rightReg))
	case ">":
		writeLine(output, fmt.Sprintf("sgt $t%d, $t%d, $t%d", resultReg, leftReg, rightReg))
	case "<=":
		writeLine(output, fmt.Sprintf("sle $t%d, $t%d, $t%d", resultReg, leftReg, rightReg))
	case ">=":
		writeLine(output, fmt.Sprintf("sge $t%d, $t%d, $t%d", resultReg, leftReg, rightReg))
	default:
		fmt.Printf("Unsupported integer operation: %s\n", operator)
		return 0, "int"
	}
	return resultReg, "int"
}

func generateBoolInfixExpression(output *strings.Builder, operator string, leftReg, rightReg int) (int, string) {
	resultReg := getNextIntRegister()
	switch operator {
	case "&&":
		writeLine(output, fmt.Sprintf("and $t%d, $t%d, $t%d", resultReg, leftReg, rightReg))
	case "||":
		writeLine(output, fmt.Sprintf("or $t%d, $t%d, $t%d", resultReg, leftReg, rightReg))
	case "==":
		writeLine(output, fmt.Sprintf("xor $t%d, $t%d, $t%d", resultReg, leftReg, rightReg))
		writeLine(output, fmt.Sprintf("sltiu $t%d, $t%d, 1", resultReg, resultReg))
	case "!=":
		writeLine(output, fmt.Sprintf("xor $t%d, $t%d, $t%d", resultReg, leftReg, rightReg))
		writeLine(output, fmt.Sprintf("sltu $t%d, $zero, $t%d", resultReg, resultReg))
	default:
		fmt.Printf("Unsupported boolean operation: %s\n", operator)
		return 0, "bool"
	}
	return resultReg, "bool"
}

func generateSpeakNow(output *strings.Builder, node *ast.CallExpression) (int, string) {
	if len(node.Arguments) != 1 {
		fmt.Println("Error: SpeakNow expects exactly one argument")
		return 0, ""
	}

	reg, valType := generateNode(output, node.Arguments[0])

	switch valType {
	case "int":
		writeLines(output, []string{
			fmt.Sprintf("move $a0, $t%d", reg),
			"li $v0, 1", // System call for print integer
			"syscall",
		})
	case "bool":
		labelFalse := getNextLabel()
		labelEnd := getNextLabel()
		writeLines(output, []string{
			fmt.Sprintf("beq $t%d, $zero, %s", reg, labelFalse),
			"la $a0, true_str",
			fmt.Sprintf("j %s", labelEnd),
			fmt.Sprintf("%s:", labelFalse),
			"la $a0, false_str",
			fmt.Sprintf("%s:", labelEnd),
			"li $v0, 4", // System call for print string
			"syscall",
		})
	default:
		fmt.Printf("Unsupported type for SpeakNow: %s\n", valType)
		return 0, ""
	}

	writeLines(output, []string{
		"li $v0, 4",
		"la $a0, newline",
		"syscall",
	})

	return reg, valType
}

func floatComparisonOp(operator string) string {
	switch operator {
	case "<":
		return "lt"
	case ">":
		return "gt"
	case "<=":
		return "le"
	case ">=":
		return "ge"
	case "==":
		return "eq"
	case "!=":
		return "ne"
	default:
		return ""
	}
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
func generateIfExpression(output *strings.Builder, ifExpr *ast.IfExpression) (int, string) {
	resetRegisterAllocation()

	// Generate code for the condition
	condReg, _ := generateNode(output, ifExpr.Condition)

	labelElse := getNextLabel()
	labelEnd := getNextLabel()

	// If condition is false, jump to else
	writeLine(output, fmt.Sprintf("beq $t%d, $zero, %s", condReg, labelElse))

	// Generate code for consequence
	consequenceReg, consequenceType := generateNode(output, ifExpr.Consequence)

	// Jump to end after consequence
	writeLine(output, fmt.Sprintf("j %s", labelEnd))

	// Else label
	writeLine(output, fmt.Sprintf("%s:", labelElse))

	// Generate code for alternative (if it exists)
	if ifExpr.Alternative != nil {
		_, alternativeType := generateNode(output, ifExpr.Alternative)
		if consequenceType != alternativeType {
			fmt.Printf("Warning: Mismatched types in if-else: %s and %s\n", consequenceType, alternativeType)
		}
	}

	// End label
	writeLine(output, fmt.Sprintf("%s:", labelEnd))

	return consequenceReg, consequenceType
}
func generateBlockStatement(output *strings.Builder, block *ast.BlockStatement) (int, string) {
	var lastReg int
	var lastType string
	for _, statement := range block.Statements {
		lastReg, lastType = generateNode(output, statement)
	}
	return lastReg, lastType
}
