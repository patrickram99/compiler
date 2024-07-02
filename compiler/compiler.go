package compiler

import (
	"fmt"
	"main/ast"
	"strings"
)

type SymbolTable struct {
	symbols     map[string]*Symbol
	parent      *SymbolTable
	nextAddress int
}

type Symbol struct {
	Name    string
	Type    string
	Address int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		symbols:     make(map[string]*Symbol),
		nextAddress: 0,
	}
}

func (st *SymbolTable) Define(name, typ string) *Symbol {
	symbol := &Symbol{Name: name, Type: typ, Address: st.nextAddress}
	st.symbols[name] = symbol
	st.nextAddress += 4 // Assume all variables/parameters are 4 bytes
	return symbol
}

func (st *SymbolTable) Resolve(name string) (*Symbol, bool) {
	symbol, ok := st.symbols[name]
	if !ok && st.parent != nil {
		return st.parent.Resolve(name)
	}
	return symbol, ok
}

func (st *SymbolTable) NewEnclosedSymbolTable() *SymbolTable {
	enclosed := NewSymbolTable()
	enclosed.parent = st
	return enclosed
}

var (
	symbolTable        *SymbolTable
	currentRegister    int
	labelCount         int
	functionName       string
	intRegisterCount   int
	floatRegisterCount int
)

func GenerateMIPS(node ast.Node) string {
	var output strings.Builder
	symbolTable = NewSymbolTable()
	currentRegister = 0
	labelCount = 0
	functionName = ""

	writeLines(&output, []string{
		".data",
		"newline: .asciiz \"\\n\"",
		"true_str: .asciiz \"true\"",
		"false_str: .asciiz \"false\"",
		".text",
		".globl main",
	})

	if program, ok := node.(*ast.Program); ok {
		generateProgram(&output, program)
	} else {
		generateNode(&output, node)
	}

	return output.String()
}

func generateProgram(output *strings.Builder, program *ast.Program) {
	// Generate main function
	writeLines(output, []string{
		"main:",
		"move $fp, $sp",
		"sw $ra, 0($sp)",
		"addi $sp, $sp, -4",
	})

	for _, statement := range program.Statements {
		generateNode(output, statement)
	}

	// Main function epilogue
	writeLines(output, []string{
		"lw $ra, 4($sp)",
		"addi $sp, $sp, 4",
		"jr $ra",
	})
}

func generateNode(output *strings.Builder, node ast.Node) (int, string) {
	switch n := node.(type) {
	case *ast.Program:
		generateProgram(output, n)
		return 0, ""
	case *ast.ExpressionStatement:
		return generateNode(output, n.Expression)
	case *ast.InfixExpression:
		return generateInfixExpression(output, n)
	case *ast.IntegerLiteral:
		reg := getNextRegister()
		writeLine(output, fmt.Sprintf("li $t%d, %d", reg, n.Value))
		return reg, "int"
	case *ast.Boolean:
		reg := getNextRegister()
		if n.Value {
			writeLine(output, fmt.Sprintf("li $t%d, 1", reg))
		} else {
			writeLine(output, fmt.Sprintf("li $t%d, 0", reg))
		}
		return reg, "bool"
	case *ast.IfExpression:
		return generateIfExpression(output, n)

	case *ast.BlockStatement:
		generateBlockStatement(output, n)
		return 0, ""

	case *ast.CallExpression:
		if ident, ok := n.Function.(*ast.Variable); ok && ident.Value == "SpeakNow" {
			return generateSpeakNow(output, n)
		}
	case *ast.LetStatement:
		return generateLetStatement(output, n)
	case *ast.Variable:
		return generateVariable(output, n)
	}
	return 0, ""
}

func generateLetStatement(output *strings.Builder, node *ast.LetStatement) (int, string) {
	valueReg, valueType := generateNode(output, node.Value)
	symbol := symbolTable.Define(node.Name.Value, valueType)
	if functionName == "" {
		// Global variable
		writeLine(output, fmt.Sprintf(".data"))
		writeLine(output, fmt.Sprintf("%s: .word 0", node.Name.Value))
		writeLine(output, fmt.Sprintf(".text"))
		writeLine(output, fmt.Sprintf("sw $t%d, %s", valueReg, node.Name.Value))
	} else {
		// Local variable
		writeLine(output, fmt.Sprintf("sw $t%d, -%d($fp)", valueReg, symbol.Address))
	}
	return valueReg, valueType
}

func generateVariable(output *strings.Builder, node *ast.Variable) (int, string) {
	symbol, ok := symbolTable.Resolve(node.Value)
	if !ok {
		fmt.Printf("Error: Undefined variable %s\n", node.Value)
		return 0, ""
	}
	resultReg := getNextRegister()
	if functionName == "" {
		// Global variable
		writeLine(output, fmt.Sprintf("lw $t%d, %s", resultReg, node.Value))
	} else if symbol.Address > 0 {
		// Parameter
		writeLine(output, fmt.Sprintf("lw $t%d, %d($fp)", resultReg, symbol.Address))
	} else {
		// Local variable
		writeLine(output, fmt.Sprintf("lw $t%d, -%d($fp)", resultReg, -symbol.Address))
	}
	return resultReg, symbol.Type
}

func resetRegisterAllocation() {
	currentRegister = 0
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

	argReg, argType := generateNode(output, node.Arguments[0])

	switch argType {
	case "int":
		writeLines(output, []string{
			fmt.Sprintf("move $a0, $t%d", argReg),
			"li $v0, 1", // System call for print integer
			"syscall",
		})
	case "bool":
		labelFalse := getNextLabel()
		labelEnd := getNextLabel()
		writeLines(output, []string{
			fmt.Sprintf("beq $t%d, $zero, %s", argReg, labelFalse),
			"la $a0, true_str",
			fmt.Sprintf("j %s", labelEnd),
			fmt.Sprintf("%s:", labelFalse),
			"la $a0, false_str",
			fmt.Sprintf("%s:", labelEnd),
			"li $v0, 4", // System call for print string
			"syscall",
		})
	default:
		fmt.Printf("Unsupported type for SpeakNow: %s\n", argType)
		return 0, ""
	}

	// Print newline
	writeLines(output, []string{
		"li $v0, 4",
		"la $a0, newline",
		"syscall",
	})

	return argReg, argType
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
	condReg, _ := generateNode(output, ifExpr.Condition)

	labelElse := getNextLabel()
	labelEnd := getNextLabel()

	// If condition is false, jump to else
	writeLine(output, fmt.Sprintf("beq $t%d, $zero, %s", condReg, labelElse))

	// Generate code for consequence
	if ifExpr.Consequence != nil {
		generateBlockStatement(output, ifExpr.Consequence)
	}

	// Jump to end after consequence
	writeLine(output, fmt.Sprintf("j %s", labelEnd))

	// Else label
	writeLine(output, fmt.Sprintf("%s:", labelElse))

	// Generate code for alternative (if it exists)
	if ifExpr.Alternative != nil {
		generateBlockStatement(output, ifExpr.Alternative)
	}

	// End label
	writeLine(output, fmt.Sprintf("%s:", labelEnd))

	return 0, "" // If-expressions don't return a value in this implementation
}

func generateBlockStatement(output *strings.Builder, block *ast.BlockStatement) {
	for _, statement := range block.Statements {
		generateNode(output, statement)
	}
}

func getNextRegister() int {
	reg := currentRegister
	currentRegister = (currentRegister + 1) % 8 // Use only $t0 to $t7
	return reg
}

func writeLines(output *strings.Builder, lines []string) {
	for _, line := range lines {
		output.WriteString(line + "\n")
	}
}

func writeLine(output *strings.Builder, line string) {
	output.WriteString(line + "\n")
}

func generateFunction(output *strings.Builder, node *ast.FunctionLiteral) (int, string) {
	prevFunctionName := functionName
	functionName = fmt.Sprintf("func_%d", labelCount)
	labelCount++

	writeLine(output, fmt.Sprintf("%s:", functionName))

	// Function prologue
	writeLines(output, []string{
		"sw $ra, 0($sp)",
		"sw $fp, -4($sp)",
		"move $fp, $sp",
		fmt.Sprintf("addi $sp, $sp, -%d", 8+len(node.Parameters)*4),
	})

	// Save used registers
	for i := 0; i < 8; i++ {
		writeLine(output, fmt.Sprintf("sw $t%d, -%d($fp)", i, 12+i*4))
	}
	writeLine(output, fmt.Sprintf("addi $sp, $sp, -%d", 8*4))

	// Create a new symbol table for the function's scope
	enclosedSymbolTable := symbolTable.NewEnclosedSymbolTable()
	prevSymbolTable := symbolTable
	symbolTable = enclosedSymbolTable

	// Add parameters to the symbol table
	for i, param := range node.Parameters {
		symbol := symbolTable.Define(param.Value, "int") // Assume all parameters are ints for now
		symbol.Address = 4 * (i + 2)                     // Parameters are above the frame pointer
	}

	// Generate code for the function body
	generateNode(output, node.Body)

	// Function epilogue
	writeLines(output, []string{
		// Restore used registers
		fmt.Sprintf("addi $sp, $sp, %d", 8*4),
	})
	for i := 0; i < 8; i++ {
		writeLine(output, fmt.Sprintf("lw $t%d, -%d($fp)", i, 12+i*4))
	}
	writeLines(output, []string{
		"move $sp, $fp",
		"lw $ra, 0($sp)",
		"lw $fp, -4($sp)",
		"jr $ra",
	})

	// Restore the previous symbol table and function name
	symbolTable = prevSymbolTable
	functionName = prevFunctionName

	return 0, "function"
}

func generateFunctionCall(output *strings.Builder, node *ast.CallExpression) (int, string) {
	if ident, ok := node.Function.(*ast.Variable); ok && ident.Value == "SpeakNow" {
		return generateSpeakNow(output, node)
	}

	// Evaluate and push arguments
	for i := len(node.Arguments) - 1; i >= 0; i-- {
		argReg, _ := generateNode(output, node.Arguments[i])
		writeLine(output, fmt.Sprintf("sw $t%d, 0($sp)", argReg))
		writeLine(output, "addi $sp, $sp, -4")
	}

	// Call the function
	funcName := node.Function.(*ast.Variable).Value
	writeLine(output, fmt.Sprintf("jal %s", funcName))

	// Clean up the stack
	writeLine(output, fmt.Sprintf("addi $sp, $sp, %d", len(node.Arguments)*4))

	// The result is in $v0, move it to a temporary register
	resultReg := getNextRegister()
	writeLine(output, fmt.Sprintf("move $t%d, $v0", resultReg))

	return resultReg, "int" // Assume all functions return int for now
}
