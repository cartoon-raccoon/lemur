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
	program, err := p.Parse()
	if program == nil || err != nil {
		t.Fatalf(err.Error())
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
	program, err := p.Parse()
	if program == nil {
		t.Fatalf(err.Error())
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

	prog, err := p.Parse()
	if err != nil {
		t.Fatalf("Error: parser encountered an error - %s", err)
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

		prog, err := p.Parse()
		if err != nil {
			t.Errorf("Error on statement %d", i)
			t.Errorf("Error: parser encountered an error - %s", err)
			continue
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
	prog, err := p.Parse()

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

	if prog == nil {
		t.Fatalf("Invalid program")
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
		prog, err := p.Parse()
		if err != nil {
			t.Fatalf("Error parsing program: %s", err)
		}
		if prog == nil {
			t.Fatalf("AST could not be generated")
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
		prog, err := p.Parse()
		if err != nil {
			t.Fatalf("Test %d: Could not parse AST", i)
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
	input := "if (x < y || !false) { return (x + y) * 2; } else { return x * y; }"

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
	prog, err := p.Parse()

	if err != nil {
		t.Fatalf("Error while parsing: %s", err)
	}
	if prog == nil {
		t.Fatalf("Could not parse AST")
	}

	// testing top level statements
	if len := len(prog.Statements); len != 1 {
		t.Errorf("Expected 1 statement, got %d", len)
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

	if len := len(ifexpr.Alternative.Statements); len != 1 {
		t.Errorf("Incorrect statements in blockstmt, got %d", len)
		t.Logf(ifexpr.Result.TokenLiteral())
		for _, stmt := range ifexpr.Result.Statements {
			t.Logf(stmt.String())
		}
	}
}
