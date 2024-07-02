package repl

import (
	"io"
	"log"
	"main/compiler"
	"main/evaluator"
	"main/lexer"
	"main/object"
	"main/parser"
	"os"
)

func writeToFile(filePath string, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func Start(filePath string, out io.Writer) {
	env := object.NewEnvironment()

	// Read the entire file content
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// Process the file content as a single input string
	line := string(fileContent)
	l := lexer.New(line)
	p := parser.New(l)
	program := p.ParseProgram()

	// PrintAST(program, "  ")

	// Create Graphviz image
	err = CreateGraphvizImage(program, "ast.jpeg")
	if err != nil {
		log.Fatalf("Error generating AST: %v", err)
	} else {
		log.Println("AST saved in: ast.jpeg")
	}

	mipsCode := compiler.GenerateMIPS(program)
	// fmt.Println(mipsCode) // Print the generated code

	writeToFile("out.s", mipsCode)

	if len(p.Errors()) != 0 {
		printParserErrors(out, p.Errors())
		return
	}

	evaluated := evaluator.Eval(program, env)

	if evaluated != nil {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
