package repl

import (
	"bufio"
	"fmt"
	"main/ast"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func PrintAST(node ast.Node, indent string) {
	switch n := node.(type) {
	case *ast.Program:
		fmt.Println(indent + "Program:")
		for _, stmt := range n.Statements {
			PrintAST(stmt, indent+"  ")
		}
	case *ast.IfExpression:
		fmt.Println(indent + "IfExpression:")
		fmt.Println(indent + "  Condition:")
		PrintAST(n.Condition, indent+"    ")
		fmt.Println(indent + "  Consequence:")
		PrintAST(n.Consequence, indent+"    ")
		if n.Alternative != nil {
			fmt.Println(indent + "  Alternative:")
			PrintAST(n.Alternative, indent+"    ")
		}
	case *ast.BlockStatement:
		fmt.Println(indent + "BlockStatement:")
		for _, stmt := range n.Statements {
			PrintAST(stmt, indent+"  ")
		}
	case *ast.FunctionLiteral:
		fmt.Println(indent + "FunctionLiteral:")
		fmt.Println(indent + "  Parameters:")
		for _, param := range n.Parameters {
			PrintAST(param, indent+"    ")
		}
		fmt.Println(indent + "  Body:")
		PrintAST(n.Body, indent+"    ")
	case *ast.CallExpression:
		fmt.Println(indent + "CallExpression:")
		fmt.Println(indent + "  Function:")
		PrintAST(n.Function, indent+"    ")
		fmt.Println(indent + "  Arguments:")
		for _, arg := range n.Arguments {
			PrintAST(arg, indent+"    ")
		}
	case *ast.Boolean:
		fmt.Printf(indent+"Boolean: %v (%v)\n", n.Value, n.TokenLiteral())
	case *ast.IntegerLiteral:
		fmt.Printf(indent+"IntegerLiteral: %v (%v)\n", n.Value, n.TokenLiteral())
	case *ast.FloatLiteral:
		fmt.Printf(indent+"FloatLiteral: %v (%v)\n", n.Value, n.TokenLiteral())
	case *ast.PrefixExpression:
		fmt.Println(indent + "PrefixExpression:")
		fmt.Printf(indent+"  Operator: %v\n", n.Operator)
		fmt.Println(indent + "  Right:")
		PrintAST(n.Right, indent+"    ")
	case *ast.InfixExpression:
		fmt.Println(indent + "InfixExpression:")
		fmt.Println(indent + "  Left:")
		PrintAST(n.Left, indent+"    ")
		fmt.Printf(indent+"  Operator: %v\n", n.Operator)
		fmt.Println(indent + "  Right:")
		PrintAST(n.Right, indent+"    ")
	case *ast.LetStatement:
		fmt.Println(indent + "LetStatement:")
		fmt.Println(indent + "  Name:")
		PrintAST(n.Name, indent+"    ")
		fmt.Println(indent + "  Value:")
		PrintAST(n.Value, indent+"    ")
	case *ast.Variable:
		fmt.Printf(indent+"Variable: %v (%v)\n", n.Value, n.TokenLiteral())
	case *ast.ReturnStatement:
		fmt.Println(indent + "ReturnStatement:")
		if n.ReturnValue != nil {
			fmt.Println(indent + "  ReturnValue:")
			PrintAST(n.ReturnValue, indent+"    ")
		}
	case *ast.ExpressionStatement:
		fmt.Println(indent + "ExpressionStatement:")
		if n.Expression != nil {
			PrintAST(n.Expression, indent+"  ")
		}
	default:
		fmt.Printf(indent+"Unknown node type: %T\n", node)
	}
}

type dotNode struct {
	ID   string
	Name string
}

var nodeCounter int

func nextNodeID() string {
	nodeCounter++
	return fmt.Sprintf("Node%d", nodeCounter)
}

func writeDotNode(n dotNode, f *os.File) {
	fmt.Fprintf(f, "%s [label=\"%s\"];\n", n.ID, n.Name)
}

func writeDotEdge(fromID, toID string, f *os.File) {
	fmt.Fprintf(f, "%s -> %s;\n", fromID, toID)
}

func generateDot(node ast.Node, parentID string, f *os.File) string {
	nodeID := nextNodeID()
	switch n := node.(type) {
	case *ast.Program:
		writeDotNode(dotNode{nodeID, "Program"}, f)
		for _, stmt := range n.Statements {
			childID := generateDot(stmt, nodeID, f)
			writeDotEdge(nodeID, childID, f)
		}
	case *ast.IfExpression:
		writeDotNode(dotNode{nodeID, "IfExpression"}, f)
		conditionID := generateDot(n.Condition, nodeID, f)
		writeDotEdge(nodeID, conditionID, f)
		consequenceID := generateDot(n.Consequence, nodeID, f)
		writeDotEdge(nodeID, consequenceID, f)
		if n.Alternative != nil {
			alternativeID := generateDot(n.Alternative, nodeID, f)
			writeDotEdge(nodeID, alternativeID, f)
		}
	case *ast.BlockStatement:
		writeDotNode(dotNode{nodeID, "BlockStatement"}, f)
		for _, stmt := range n.Statements {
			childID := generateDot(stmt, nodeID, f)
			writeDotEdge(nodeID, childID, f)
		}
	case *ast.FunctionLiteral:
		writeDotNode(dotNode{nodeID, "FunctionLiteral"}, f)
		for _, param := range n.Parameters {
			paramID := generateDot(param, nodeID, f)
			writeDotEdge(nodeID, paramID, f)
		}
		bodyID := generateDot(n.Body, nodeID, f)
		writeDotEdge(nodeID, bodyID, f)
	case *ast.CallExpression:
		writeDotNode(dotNode{nodeID, "CallExpression"}, f)
		functionID := generateDot(n.Function, nodeID, f)
		writeDotEdge(nodeID, functionID, f)
		for _, arg := range n.Arguments {
			argID := generateDot(arg, nodeID, f)
			writeDotEdge(nodeID, argID, f)
		}
	case *ast.Boolean:
		writeDotNode(dotNode{nodeID, fmt.Sprintf("Boolean: %v", n.Value)}, f)
	case *ast.StringLiteral:
		writeDotNode(dotNode{nodeID, fmt.Sprintf("String: %v", n.Value)}, f)
	case *ast.IntegerLiteral:
		writeDotNode(dotNode{nodeID, fmt.Sprintf("IntegerLiteral: %v", n.Value)}, f)
	case *ast.FloatLiteral:
		writeDotNode(dotNode{nodeID, fmt.Sprintf("FloatLiteral: %v", n.Value)}, f)
	case *ast.PrefixExpression:
		writeDotNode(dotNode{nodeID, fmt.Sprintf("PrefixExpression: %v", n.Operator)}, f)
		rightID := generateDot(n.Right, nodeID, f)
		writeDotEdge(nodeID, rightID, f)
	case *ast.InfixExpression:
		writeDotNode(dotNode{nodeID, fmt.Sprintf("InfixExpression: %v", n.Operator)}, f)
		leftID := generateDot(n.Left, nodeID, f)
		writeDotEdge(nodeID, leftID, f)
		rightID := generateDot(n.Right, nodeID, f)
		writeDotEdge(nodeID, rightID, f)
	case *ast.LetStatement:
		writeDotNode(dotNode{nodeID, "LetStatement"}, f)
		nameID := generateDot(n.Name, nodeID, f)
		writeDotEdge(nodeID, nameID, f)
		if n.Value != nil {
			valueID := generateDot(n.Value, nodeID, f)
			writeDotEdge(nodeID, valueID, f)
		}
	case *ast.Variable:
		writeDotNode(dotNode{nodeID, fmt.Sprintf("Variable: %v", n.Value)}, f)
	case *ast.ReturnStatement:
		writeDotNode(dotNode{nodeID, "ReturnStatement"}, f)
		if n.ReturnValue != nil {
			returnValueID := generateDot(n.ReturnValue, nodeID, f)
			writeDotEdge(nodeID, returnValueID, f)
		}
	case *ast.ExpressionStatement:
		writeDotNode(dotNode{nodeID, "ExpressionStatement"}, f)
		if n.Expression != nil {
			expressionID := generateDot(n.Expression, nodeID, f)
			writeDotEdge(nodeID, expressionID, f)
		}
	default:
		writeDotNode(dotNode{nodeID, fmt.Sprintf("Desconocido: %T", node)}, f)
	}

	if parentID != "" {
		writeDotEdge(parentID, nodeID, f)
	}

	return nodeID
}

func CreateGraphvizImage(node ast.Node, filename string) error {
	nodeCounter = 0

	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	dotFilename := strings.TrimSuffix(filename, ".png") + ".dot"
	f, err := os.Create(dotFilename)
	if err != nil {
		return fmt.Errorf("failed to create .dot file: %v", err)
	}
	defer f.Close()

	fmt.Fprintln(f, "digraph AST {")
	fmt.Fprintln(f, "  node [shape=box];")

	generateDot(node, "", f)

	fmt.Fprintln(f, "}")

	// Remove duplicate lines from the .dot file
	if err := RemoveDuplicateLines(dotFilename); err != nil {
		return fmt.Errorf("failed to remove duplicate lines: %v", err)
	}

	cmd := exec.Command("dot", "-Tpng", dotFilename, "-o", filename)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to generate image using graphviz: %v\nOutput: %s", err, string(output))
	}

	return nil
}

func RemoveDuplicateLines(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open .dot file: %v", err)
	}
	defer file.Close()

	lines := make(map[string]struct{})
	var uniqueLines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if _, exists := lines[line]; !exists {
			lines[line] = struct{}{}
			uniqueLines = append(uniqueLines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading .dot file: %v", err)
	}

	// Write the unique lines back to the file
	outputFile, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to write back to .dot file: %v", err)
	}
	defer outputFile.Close()

	for _, line := range uniqueLines {
		fmt.Fprintln(outputFile, line)
	}

	return nil
}
