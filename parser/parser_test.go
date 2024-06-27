package parser

import (
	"fmt"
	"main/ast"
	"main/lexer"
	"math"
	"strings"
	"testing"
)

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}
	t.Errorf("parser tiene %d errores", len(errors))
	for _, msg := range errors {
		t.Errorf("Error en el parser: %q", msg)
	}
	t.FailNow()
}

func TestDeclaracionVariables(t *testing.T) {
	input := `
	enchanted x = 5;
	enchanted y = 10.12;
	enchanted kekw = 123456;
	`
	l := lexer.New(input)
	p := New(l)

	codigo := p.ParseProgram()
	checkParserErrors(t, p)

	if codigo == nil {
		t.Fatalf("No se reconocio ningun token")
	}
	if len(codigo.Statements) != 3 {
		t.Fatalf("Se reconoció una cantidad distinta a los 3 definidos: %d", len(codigo.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"kekw"},
	}

	for i, tt := range tests {
		stmt := codigo.Statements[i]
		if !testDeclaracion(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testDeclaracion(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "enchanted" {
		t.Errorf("TokenLiteral no es enchanted, es: %T", s)
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s no es *ast.LetStatement. es: %T", s)
		return false
	}
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value no es '%s'. es: %s", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name no es '%s'. es: %s", name, letStmt.Name)
		return false
	}
	return true
}

func TestReturns(t *testing.T) {
	input := `
	hi 5;
	hi 10;
	hi 993322;
	`
	l := lexer.New(input)
	p := New(l)
	fmt.Print(p)

	codigo := p.ParseProgram()
	checkParserErrors(t, p)

	if len(codigo.Statements) != 3 {
		t.Fatalf("Se reconoció una cantidad distinta a los 3 definidos: %d", len(codigo.Statements))
	}

	for _, stmt := range codigo.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("Return no coincide *ast.returnStatement. es: %T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "hi" {
			t.Errorf("returnStmt.TokenLiteral no es 'hi', es: %q",
				returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "omggg"

	l := lexer.New(input)
	p := New(l)
	codigo := p.ParseProgram()
	checkParserErrors(t, p)

	if len(codigo.Statements) != 1 {
		t.Fatalf("El programa no reconocio el statement. Es: %d", len(codigo.Statements))
	}
	stmt, ok := codigo.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] no es ast.ExpressionStatement. Es: %T",
			codigo.Statements[0])
	}

	variable, ok := stmt.Expression.(*ast.Variable)

	if !ok {
		t.Fatalf("Expresion no es *ast.Identifier. Es: %T", stmt.Expression)
	}

	if variable.Value != "omggg" {
		t.Errorf("Variable no es %s. sino: %s", "omggg", variable.Value)
	}
	if variable.TokenLiteral() != "omggg" {
		t.Errorf("variable.TokenLiteral no es %s. sino: %s", "omggg",
			variable.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5",
			literal.TokenLiteral())
	}
}

func TestFloatExpressions(t *testing.T) {
	input := "5.5484;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("El programa no registro suficientes declaraciones. SOn: %d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] no es ast.ExpressionStatement. Sino: %T",
			program.Statements[0])
	}
	literal, ok := stmt.Expression.(*ast.FloatLiteral)
	if !ok {
		t.Fatalf("exp not *ast.FloatLiteral. Es: %T", stmt.Expression)
	}
	if literal.Value != 5.5484 {
		t.Errorf("literal.Value no es %d. Es: %f", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5.5484" {
		t.Errorf("literal.TokenLiteral no es %s. Es: %s", "5.5",
			literal.TokenLiteral())
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integ.TokenLiteral())
		return false
	}
	return true
}

func testFloatLiteral(t *testing.T, il ast.Expression, value float64) bool {
	floatNum, ok := il.(*ast.FloatLiteral)
	if !ok {
		t.Errorf("il not *ast.FloatLiteral. got=%T", il)
		return false
	}

	if !closeEnough(floatNum.Value, value, floatTolerance) {
		t.Errorf("floatNum.Value not %f. got=%f", value, floatNum.Value)
		return false
	}

	expectedSubstring := fmt.Sprintf("%.1f", value)
	if !strings.Contains(floatNum.TokenLiteral(), expectedSubstring) {
		t.Errorf("floatNum.TokenLiteral does not contain %s. got=%s",
			expectedSubstring, floatNum.TokenLiteral())
		return false
	}

	return true
}

const floatTolerance = 0.000001

func closeEnough(a, b, tolerance float64) bool {
	return math.Abs(a-b) < tolerance
}

func TestPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}
	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	type value struct {
		intValue   int64
		floatValue float64
		isFloat    bool
	}

	infixTests := []struct {
		input      string
		leftValue  value
		operator   string
		rightValue value
	}{
		{"5 + 5;", value{intValue: 5}, "+", value{intValue: 5}},
		{"5 - 5;", value{intValue: 5}, "-", value{intValue: 5}},
		{"5 * 5;", value{intValue: 5}, "*", value{intValue: 5}},
		{"5 / 5;", value{intValue: 5}, "/", value{intValue: 5}},
		{"5 > 5;", value{intValue: 5}, ">", value{intValue: 5}},
		{"5 < 5;", value{intValue: 5}, "<", value{intValue: 5}},
		{"5 == 5;", value{intValue: 5}, "==", value{intValue: 5}},
		{"5 != 5;", value{intValue: 5}, "!=", value{intValue: 5}},
		{"5.1 + 2.2;", value{floatValue: 5.1, isFloat: true}, "+", value{floatValue: 2.2, isFloat: true}},
		{"5.1 - 2.2;", value{floatValue: 5.1, isFloat: true}, "-", value{floatValue: 2.2, isFloat: true}},
		{"5.1 * 2.2;", value{floatValue: 5.1, isFloat: true}, "*", value{floatValue: 2.2, isFloat: true}},
		{"5.1 / 2.2;", value{floatValue: 5.1, isFloat: true}, "/", value{floatValue: 2.2, isFloat: true}},
		{"5.1 > 2.2;", value{floatValue: 5.1, isFloat: true}, ">", value{floatValue: 2.2, isFloat: true}},
		{"5.1 < 2.2;", value{floatValue: 5.1, isFloat: true}, "<", value{floatValue: 2.2, isFloat: true}},
		{"5.1 == 2.2;", value{floatValue: 5.1, isFloat: true}, "==", value{floatValue: 2.2, isFloat: true}},
		{"5.1 != 2.2;", value{floatValue: 5.1, isFloat: true}, "!=", value{floatValue: 2.2, isFloat: true}},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
		}

		if tt.leftValue.isFloat {
			if !testFloatLiteral(t, exp.Left, tt.leftValue.floatValue) {
				return
			}
		} else {
			if !testIntegerLiteral(t, exp.Left, tt.leftValue.intValue) {
				return
			}
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}

		if tt.rightValue.isFloat {
			if !testFloatLiteral(t, exp.Right, tt.rightValue.floatValue) {
				return
			}
		} else {
			if !testIntegerLiteral(t, exp.Right, tt.rightValue.intValue) {
				return
			}
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3.5 + 4.1; -5.2 * 5",
			"(3.5 + 4.1)((-5.2) * 5)",
		},
		{
			"5.5 > 4.2 == 3.3 < 4.8",
			"((5.5 > 4.2) == (3.3 < 4.8))",
		},
		{
			"5.1 < 4.6 != 3.7 > 4.2",
			"((5.1 < 4.6) != (3.7 > 4.2))",
		},
		{
			"3.8 + 4 * 5 == 3 * 1.0 + 4 * 5",
			"((3.8 + (4 * 5)) == ((3 * 1.0) + (4 * 5)))",
		},
		{
			"3.2 + 4.2 * 5.5 == 3.1 * 1.2 + 4.0 * 5.5",
			"((3.2 + (4.2 * 5.5)) == ((3.1 * 1.2) + (4.0 * 5.5)))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}
