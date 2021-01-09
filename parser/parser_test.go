package parser

import (
	"testing"

	"github.com/cartoon-raccoon/monkey-jit/ast"
	"github.com/cartoon-raccoon/monkey-jit/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `let x = 10;
	let y = 5_000;
	let nice = 69;`

	l := lexer.New(input)
	p, err := New(l)
	if err != nil {
		t.Fatalf("Error instantiating parser: malformed input")
	}
	program := p.Parse()
	if program == nil {
		t.Errorf("Error while parsing program")
	}
	if p.checkErrors() != nil {
		for _, err := range p.checkErrors() {
			t.Logf(err.Error())
		}
		t.Fatalf("Aborting test")
	}

	tests := []struct {
		expectedIdent string
	}{
		{"x"},
		{"y"},
		{"nice"},
	}

	for i, ident := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, ident.expectedIdent) {
			return
		}
	}
}

func testLetStatement(t *testing.T, stmt ast.Statement, ident string) bool {
	if stmt.TokenLiteral() != lexer.LET {
		t.Errorf("error: wrong token, expected \"let\" got %q", stmt.TokenLiteral())
	}
	letstmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Fatalf("statement not let statement: got %T", stmt)
	}
	if letstmt.Name.Value != ident {
		t.Errorf("wrong identifier: got %q", letstmt.Name.Value)
	}
	_, ok = letstmt.Value.(*ast.Int)
	if !ok {
		t.Errorf("value is not int")
	}

	return true
}

func TestReturnStatement(t *testing.T) {
	input := `return x;
	return 0.78;
	return "what";`

	l := lexer.New(input)
	p, err := New(l)
	if err != nil {
		t.Fatalf("Error instantiating parser: malformed input")
	}
	program := p.Parse()
	if program == nil {
		t.Errorf("Error while parsing program")
	}
	if p.checkErrors() != nil {
		for _, err := range p.checkErrors() {
			t.Logf(err.Error())
		}
		t.Fatalf("Aborting test")
	}

	if len := len(program.Statements); len != 3 {
		t.Errorf("Wrong number of statements: expected 3, got %d", len)
		for _, stmt := range program.Statements {
			t.Error(stmt.String())
		}
	}

	for _, stmt := range program.Statements {
		if stmt.TokenLiteral() != lexer.RETURN {
			t.Fatalf("Error: expected RETURN, got %q", stmt.TokenLiteral())
		}
		_, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("statement is not a return statement")
		}
	}
}

func TestIdentExpr(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p, err := New(l)
	if err != nil {
		t.Fatalf("Error: lexer encountered an error - %s", err)
	}

	prog := p.Parse()
	if prog == nil {
		t.Errorf("Error while parsing program, checking for errors")
	}
	if p.checkErrors() != nil {
		for _, err := range p.checkErrors() {
			t.Logf(err.Error())
		}
		t.Fatalf("Aborting test")
	} else {
		t.Logf("No errors found")
	}

	if len := len(prog.Statements); len != 1 {
		t.Fatalf("Not enough statements, got: %d", len)
	}

	if stmt, ok := prog.Statements[0].(*ast.ExprStatement); !ok {
		t.Fatalf("Statement is not expression statement, got %q", stmt.String())
		if ident, ok := stmt.Expression.(*ast.Identifier); ok {
			if ident.Token.Type != lexer.IDENT || ident.Value != "foobar" {
				t.Fatalf("Incorrect identifier, got %s", ident.Value)
			}
		} else {
			t.Fatalf("Statement is not identifier")
		}
	}
}

func TestLiteralExpr(t *testing.T) {
	tests := []struct {
		Input string
	}{
		{"5"},
		{"420.69"},
		{`"Hello there"`},
		{"true"},
	}

	for i, tt := range tests {
		l := lexer.New(tt.Input)
		p, err := New(l)
		if err != nil {
			t.Errorf("Error on statement %d", i)
			t.Errorf("Error: lexer encountered an error - %s", err)
			continue
		}

		prog := p.Parse()
		if prog == nil {
			t.Errorf("Error while parsing program %d, checking for errors", i)
		}
		if p.checkErrors() != nil {
			for _, err := range p.checkErrors() {
				t.Logf(err.Error())
			}
			t.Fatalf("Aborting test")
		} else {
			t.Logf("No errors found")
		}

		if len := len(prog.Statements); len != 1 {
			t.Errorf("Error on statement %d", i)
			t.Errorf("Wrong statement count: expected 1, got: %d", len)
			continue
		}

		if stmt, ok := prog.Statements[0].(*ast.ExprStatement); ok {
			switch i {
			case 0:
				intexpr, ok := stmt.Expression.(*ast.Int)
				if !ok {
					t.Errorf("Statement 0 is not an int, got %T", stmt.Expression)
					continue
				}
				if intexpr.Inner != 5 {
					t.Errorf("Int value is not 5: got %d", intexpr.Inner)
				}
			case 1:
				fltexpr, ok := stmt.Expression.(*ast.Flt)
				if !ok {
					t.Errorf("Statement 1 is not a flt, got %T", stmt.Expression)
					continue
				}
				if fltexpr.Inner != 420.69 {
					t.Errorf("Flt value is not 420.69: got %f", fltexpr.Inner)
				}
			case 2:
				strexpr, ok := stmt.Expression.(*ast.Str)
				if !ok {
					t.Errorf("Statement 2 is not a Str, got %T", stmt.Expression)
					continue
				}
				if strexpr.Inner != "Hello there" {
					t.Errorf("Str value does not match: got %s", strexpr.Inner)
				}
			case 3:
				boolexpr, ok := stmt.Expression.(*ast.Bool)
				if !ok {
					t.Errorf("Statement 3 is not a Bool, got %T", stmt.Expression)
					continue
				}
				if boolexpr.Inner != true {
					t.Errorf("Bool value does not match")
				}
			default:
				t.Fatal("How the hell did you even get here?")
			}
		} else {
			t.Fatalf("statement is not an expression statement")
		}
	}

}

func TestPrefixExpr(t *testing.T) {
	input := "-5; !hello; -420.69;"

	l := lexer.New(input)
	p, err := New(l)
	if err != nil {
		t.Fatalf("Could not generate AST; %s", err)
	}
	prog := p.Parse()
	if prog == nil {
		t.Errorf("Error while parsing program, checking for errors")
	}
	if p.checkErrors() != nil {
		for _, err := range p.checkErrors() {
			t.Logf(err.Error())
		}
		t.Fatalf("Aborting test")
	} else {
		t.Logf("No errors found")
	}

	expected := []struct {
		Token    string
		Operator string
		Right    ast.Expression
		String   string
	}{
		{lexer.SUB, "-", &ast.Int{Inner: 5}, "(-5)"},
		{lexer.BANG, "!", &ast.Identifier{Value: "hello"}, "(!hello)"},
		{lexer.SUB, "-", &ast.Flt{Inner: 420.69}, "(-420.69)"},
	}

	for i, stmt := range prog.Statements {
		if !testPrefixExpr(t, stmt, expected[i].Token, expected[i].String, i) {
			return
		}
	}
}

func testPrefixExpr(t *testing.T, stmt ast.Statement, tt string, right string, i int) bool {
	exstmt, ok := stmt.(*ast.ExprStatement)
	if !ok {
		t.Fatalf("statement %d is not expression statement, got %s", i, stmt.String())
	}
	pexpr, ok := exstmt.Expression.(*ast.PrefixExpr)
	if !ok {
		t.Fatalf("statement %d is not prefix expression, got %s", i, exstmt.String())
	}
	if pexpr.Token.Type != tt || pexpr.String() != right {
		t.Errorf("got wrong type or expr: %s", pexpr.String())
	}
	return true
}

func TestInfixExpr(t *testing.T) {
	infixTests := []struct {
		Input    string
		Left     int64
		Operator string
		Right    int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 <= 5;", 5, "<=", 5},
		{"5 >= 5;", 5, ">=", 5},
		{"5 != 5;", 5, "!=", 5},
		{"5 == 5;", 5, "==", 5},
	}

	for i, tt := range infixTests {
		l := lexer.New(tt.Input)
		p, err := New(l)
		if err != nil {
			t.Fatalf("Error parsing input")
		}
		prog := p.Parse()
		if prog == nil {
			t.Errorf("Error while parsing program, checking for errors")
		}
		if p.checkErrors() != nil {
			for _, err := range p.checkErrors() {
				t.Logf(err.Error())
			}
			t.Fatalf("Aborting test")
		} else {
			t.Logf("No errors found")
		}

		if len := len(prog.Statements); len != 1 {
			t.Errorf("wrong statement count for program %d: expected 1, got %d", i, len)
			continue
		}
		exprstmt, ok := prog.Statements[0].(*ast.ExprStatement)
		if !ok {
			t.Errorf("statement %d is not an expression statement", i)
			continue
		}
		expr, ok := exprstmt.Expression.(*ast.InfixExpr)
		if !ok {
			t.Errorf("expression of stmt %d is not infix expression", i)
			continue
		}
		left, lok := expr.Left.(*ast.Int)
		right, rok := expr.Right.(*ast.Int)
		if !lok || !rok {
			t.Errorf("one of the sides of statement %d is not an Int", i)
			t.Errorf("left side: %s", left.String())
			t.Errorf("right side: %s", right.String())
			continue
		}
		if left.Inner != tt.Left || right.Inner != tt.Right {
			t.Errorf("incorrect literal for statement %d", i)
			t.Errorf("left side: %d", left.Inner)
			t.Errorf("right side: %d", right.Inner)
		}
		if expr.Operator != tt.Operator {
			t.Errorf("operator for statement %d does not match", i)
			t.Errorf("expected %s, got %s", tt.Operator, expr.Operator)
		}

		t.Logf("program %d passed!", i)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"1 + 1",
			"(1 + 1)",
		},
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
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 || 3 * 1 + 4 * 5",
			"((3 + ((4 * (5 || 3)) * 1)) + (4 * 5))",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
	}
	for i, tt := range tests {
		l := lexer.New(tt.input)
		p, err := New(l)
		if err != nil {
			t.Fatalf("Test %d: Could not build parser", i)
		}
		prog := p.Parse()
		if prog == nil {
			t.Errorf("Error while parsing program, checking for errors")
		}
		if p.checkErrors() != nil {
			for _, err := range p.checkErrors() {
				t.Logf(err.Error())
			}
			t.Fatalf("Aborting test")
		} else {
			t.Logf("No errors found")
		}

		expr, ok := prog.Statements[0].(*ast.ExprStatement)
		if !ok {
			t.Errorf("Expected exprstatement, got %T", expr)
		}
		actual := prog.String()

		if actual != tt.expected {
			t.Errorf("Test %d: expected: %q, got: %q", i, tt.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y || !false) { 
		return (x + y) * 2; 
	} else if (x * y == 5) { 
		return x * y; 
	}
	
	if (!hello) { let night = true; }`

	// lexing
	l := lexer.New(input)
	p, err := New(l)
	if err != nil {
		t.Fatalf("Could not tokenize expression")
	}
	if p.current.Type != lexer.IF {
		panic("wrong token")
	}
	if p.next.Type != lexer.LPAREN {
		panic("wrong next token")
	}

	// parsing
	prog := p.Parse()
	if prog == nil {
		t.Errorf("Error while parsing program, checking for errors")
	}
	if p.checkErrors() != nil {
		for _, err := range p.checkErrors() {
			t.Logf(err.Error())
		}
		t.Fatalf("Aborting test")
	} else {
		t.Logf("No errors found")
	}

	// testing top level statements
	if len := len(prog.Statements); len != 2 {
		t.Errorf("Expected 2 statements, got %d", len)
		for _, stmt := range prog.Statements {
			if stmt == nil {
				t.Log("<nil>")
			} else {
				t.Logf(stmt.TokenLiteral())
				t.Logf("%T", stmt)
				t.Logf(stmt.String())
				t.Logf("-------")
			}
		}
		t.Fatalf("Aborting test")
	}

	// Asserting types, going down the tree
	stmt, ok := prog.Statements[0].(*ast.ExprStatement)
	if !ok {
		t.Fatalf("Statement is not an expression: got %T", stmt)
	}
	ifexpr, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("Expr is not an if expr: got %T", ifexpr)
	}
	cond, ok := ifexpr.Condition.(*ast.InfixExpr)
	if !ok {
		t.Fatalf("condition is not an infix expr")
	}
	if cond.Operator != lexer.LT {
		t.Errorf("Wrong conditional operator")
	}
	if len := len(ifexpr.Result.Statements); len != 1 {
		t.Errorf("Incorrect statements in blockstmt, got %d", len)
		t.Logf(ifexpr.Result.TokenLiteral())
		for _, stmt := range ifexpr.Result.Statements {
			t.Logf(stmt.String())
		}
	}

	if ifexpr.Alternative == nil {
		t.Fatalf("Alternative should not be nil")
	}

	alt, ok := ifexpr.Alternative.(*ast.IfExpression)
	if !ok {
		t.Fatalf("Alternative is not ifexpr")
	}

	cond, ok = alt.Condition.(*ast.InfixExpr)
	if !ok {
		t.Fatalf("condition is not an infix expr")
	}

	if len := len(alt.Result.Statements); len != 1 {
		t.Errorf("Incorrect statements in blockstmt, got %d", len)
		t.Logf(ifexpr.Result.TokenLiteral())
		for _, stmt := range ifexpr.Result.Statements {
			t.Logf(stmt.String())
		}
	}
}

func TestFnLiteralParsing(t *testing.T) {
	input := `fn(a, b) { let sum = x + y; }`
	l := lexer.New(input)
	p, err := New(l)
	if err != nil {
		t.Fatalf("Could not parse")
	}
	if p.current.Type != lexer.FUNCTION {
		t.Fatalf("Current token is not fn, got %s", p.current.Type)
	}
	if p.next.Type != lexer.LPAREN {
		t.Fatalf("Next token is not lparen, got %s", p.next.Type)
	}
	prog := p.Parse()
	if prog == nil {
		t.Errorf("Error while parsing program, checking for errors")
	}
	if p.checkErrors() != nil {
		for _, err := range p.checkErrors() {
			t.Logf(err.Error())
		}
		t.Fatalf("Aborting test")
	} else {
		t.Logf("No errors found")
	}

	if len := len(prog.Statements); len != 1 {
		t.Errorf("Expected 1 statement, got %d", len)
		for _, stmt := range prog.Statements {
			if stmt == nil {
				t.Logf("<nil>")
			} else {
				t.Logf(stmt.TokenLiteral())
			}
			t.Logf("-------")
		}
		t.Fatalf("Aborting test")
	}

	expr, ok := prog.Statements[0].(*ast.ExprStatement)
	if !ok {
		t.Fatalf("statement is not exprstmt, got %s", expr)
	}

	fnexpr, ok := expr.Expression.(*ast.FnLiteral)
	if !ok {
		t.Fatalf("statement is not function literal, got %s", expr)
	}

	// testing the literal itself
	if len := len(fnexpr.Params); len != 2 {
		t.Errorf("expected 2 params, got %d", len)
	}

	if len := len(fnexpr.Body.Statements); len != 1 {
		t.Errorf("expected 1 statement in body, got %d", len)
	}

	stmt, ok := fnexpr.Body.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("wrong statement, expected let, got %T", stmt)
	}
}

func TestCallExpression(t *testing.T) {
	input := `add(
		1, 
		2 * 3, 
		fn(a,b) { a + b; },
		hello("there", 6)
	);`

	l := lexer.New(input)
	p, err := New(l)
	if err != nil {
		t.Fatalf("Error in parsing: %s", err)
	}
	if p == nil {
		t.Fatalf("Could not parse AST")
	}
	if p.current.Type != lexer.IDENT {
		t.Fatalf("Incorrect current token: got %s", p.current.Type)
	}
	if p.next.Type != lexer.LPAREN {
		t.Fatalf("Incorrect next token: got %s", p.next.Type)
	}

	prog := p.Parse()
	if prog == nil {
		t.Errorf("Error while parsing program, checking for errors")
	}
	if p.checkErrors() != nil {
		for _, err := range p.checkErrors() {
			t.Logf(err.Error())
		}
		t.Fatalf("Aborting test")
	} else {
		t.Logf("No errors found")
	}

	if len := len(prog.Statements); len != 1 {
		t.Errorf("Expected 1 statement, got %d", len)
		for _, stmt := range prog.Statements {
			t.Logf(stmt.TokenLiteral())
		}
		// /t.Fatalf("Aborting test")
	}
	expr, ok := prog.Statements[0].(*ast.ExprStatement)
	if !ok {
		t.Fatalf("Expected expr stmt, got %T", expr)
	}
	call, ok := expr.Expression.(*ast.FunctionCall)
	if !ok {
		t.Fatalf("Expected function call, got %T", call)
	}
	ident, ok := call.Ident.(*ast.Identifier)
	if !ok {
		t.Errorf("Expected identifier as call name, got %T", ident)
	}
	if len := len(call.Params); len != 4 {
		t.Errorf("Expected 4 parameters, got %d", len)
		t.Logf("%d", p.current.Pos.Line)
		t.Logf("%+v\n", call)
	}
}
