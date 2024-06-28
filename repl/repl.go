package repl

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"main/evaluator"
	"main/lexer"
	"main/parser"
)

const PROMPT = "Speak Noooowww >> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		// PrintAST(program, "    ")

		// Create Graphviz image
		err := CreateGraphvizImage(program, "ast.jpeg")
		if err != nil {
			log.Fatalf("Error generado AST: %v", err)
		} else {
			log.Println("AST guardado en: ast.jpeg")
		}

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}
func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
