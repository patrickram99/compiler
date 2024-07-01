package evaluator

import (
	"main/lexer"
	"main/object"
	"main/parser"
	"math"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input  string
		output int64
	}{
		{"5", 5},
		{"12", 12},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.output)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("obj no es un entero. Sino: %T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("obj contiene un valor erroneo. Obtuvo: %d, en vez de: %d", result.Value, expected)
		return false
	}

	return true
}

func TestEvalFloatExpression(t *testing.T) {
	tests := []struct {
		input  string
		output float64
	}{
		{"5.5", 5.5},
		{"12.1234", 12.1234},
		{"-5.4", -5.4},
		{"-10.1110", -10.1110},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2.5 * 2.78 * 25", 695},
		{"-50 + 100.50 + -50", 0.5},
		{"5 * 2.5 + 10", 22.5},
		{"5 + 2 * 10.5", 26},
		{"20 + 2.5 * -10", -5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.output)
	}
}

func testFloatObject(t *testing.T, obj object.Object, expected float64) bool {
	const tolerance = 1e-9 // Adjust this value based on the precision needed
	var resultValue float64

	switch result := obj.(type) {
	case *object.Float:
		resultValue = result.Value
	case *object.Integer:
		resultValue = float64(result.Value)
	default:
		t.Errorf("obj no es un flotante ni un entero. Sino: %T (%+v)", obj, obj)
		return false
	}

	if math.Abs(resultValue-expected) > tolerance {
		t.Errorf("obj contiene un valor erroneo. Obtuvo: %f, en vez de: %f", resultValue, expected)
		return false
	}

	return true
}

func TestEvalBoolExpression(t *testing.T) {
	tests := []struct {
		input  string
		output bool
	}{
		{"SparksFly", true},
		{"BadBlood", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"SparksFly != BadBlood", true},
		{"BadBlood != SparksFly", true},
		{"(1 < 2) == SparksFly", true},
		{"(1 < 2) == BadBlood", false},
		{"(1 > 2) == SparksFly", false},
		{"(1 > 2) == BadBlood", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBoolObject(t, evaluated, tt.output)
	}
}

func testBoolObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Bool)
	if !ok {
		t.Errorf("obj no es un booleano. Sino: %T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("obj contiene un valor erroneo. Obtuvo: %t, en vez de: %t", result.Value, expected)
		return false
	}

	return true
}

func TestExclOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!SparksFly", false},
		{"!BadBlood", true},
		{"!5", false},
		{"!!SparksFly", true},
		{"!!BadBlood", false},
		{"!!5", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBoolObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"LoverEra (SparksFly) { 10 }", 10},
		{"LoverEra (BadBlood) { 10 }", nil},
		{"LoverEra (1) { 10 }", 10},
		{"LoverEra (1 < 2) { 10 }", 10},
		{"LoverEra (1 > 2) { 10 }", nil},
		{"LoverEra (1 > 2) { 10 } RepEra { 20 }", 20},
		{"LoverEra (1 < 2) { 10 } RepEra { 20 }", 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}
