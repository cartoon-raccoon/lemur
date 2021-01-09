package eval

import (
	"github.com/cartoon-raccoon/monkey-jit/ast"
	"github.com/cartoon-raccoon/monkey-jit/lexer"
	"github.com/cartoon-raccoon/monkey-jit/object"
)

var (
	// NULL - the single null object
	NULL = &object.Null{}
	// TRUE - An invariant
	TRUE = &object.Boolean{Value: true}
	// FALSE - An invariant
	FALSE = &object.Boolean{Value: false}
)

// Evaluator epresents the program that walks the tree
type Evaluator struct {
	Data map[string]object.Object
}

// Evaluate runs the evaluator, walking the tree and executing code
func (e *Evaluator) Evaluate(node ast.Node) object.Object {
	switch node.(type) {
	case *ast.Program:
		ret := &object.ProgramResult{}
		ret.Results = []object.Object{}
		for _, stmt := range node.(*ast.Program).Statements {
			ret.Results = append(ret.Results, e.Evaluate(stmt))
		}
		return ret
	case ast.Statement:
		switch node.(ast.Statement).(type) {
		case *ast.LetStatement:
			letstmt := node.(ast.Statement).(*ast.LetStatement)
			e.Data[letstmt.Name.String()] = e.Evaluate(letstmt.Value)
		case *ast.ExprStatement:
			expr := node.(ast.Statement).(*ast.ExprStatement)
			return e.Evaluate(expr.Expression)
		case *ast.ReturnStatement:

		case *ast.BlockStatement:
			blkstmt := node.(ast.Statement).(*ast.BlockStatement)
			for _, stmt := range blkstmt.Statements {
				e.Evaluate(stmt)
			}
		default:
			return nil
		}
	case ast.Expression:
		switch node.(ast.Expression).(type) {
		case *ast.Identifier:
			ident := node.(ast.Expression).(*ast.Identifier)
			data, ok := e.Data[ident.Value]
			if !ok {
				//throw error
			}
			return data
		case *ast.PrefixExpr:
			pexpr := node.(ast.Expression).(*ast.PrefixExpr)
			return e.evalPrefixExpr(pexpr)
		case *ast.InfixExpr:
			iexpr := node.(ast.Expression).(*ast.InfixExpr)
			return e.evalInfixExpr(iexpr)
		case *ast.IfExpression:
			ifexpr := node.(ast.Expression).(*ast.IfExpression)
			condition := e.Evaluate(ifexpr.Condition)
			if evaluateTruthiness(condition) {
				return e.Evaluate(ifexpr.Result)
			}
			switch ifexpr.Alternative.(type) {
			case *ast.BlockStatement:
				return e.Evaluate(ifexpr.Alternative.(*ast.BlockStatement))
			case *ast.IfExpression:
				return e.Evaluate(ifexpr.Alternative.(*ast.IfExpression))
			default:
				//todo: throw error
				return nil
			}
		case *ast.FnLiteral:
		case *ast.FunctionCall:
		case *ast.Int:
			intexpr := node.(ast.Expression).(*ast.Int)
			return &object.Integer{Value: intexpr.Inner}
		case *ast.Flt:
			fltexpr := node.(ast.Expression).(*ast.Flt)
			return &object.Float{Value: fltexpr.Inner}
		case *ast.Str:
			strexpr := node.(ast.Expression).(*ast.Str)
			return &object.String{Value: strexpr.Inner}
		case *ast.Bool:
			boolexpr := node.(ast.Expression).(*ast.Bool)
			return nativeBooltoObj(boolexpr.Inner)
		default:
			return NULL
		}
	default:
		return nil
	}
	return nil
}

func (e *Evaluator) evalPrefixExpr(expr *ast.PrefixExpr) object.Object {
	switch expr.Operator {
	case "!":
		return e.evalBangPExpr(expr.Right)
	case "-":
		return e.evalMinusPExpr(expr.Right)
	default:
		return NULL
	}
}

func (e *Evaluator) evalBangPExpr(expr ast.Expression) object.Object {
	pexpr := e.Evaluate(expr)
	truth := evaluateTruthiness(pexpr)
	return nativeBooltoObj(!truth)
}

func (e *Evaluator) evalMinusPExpr(expr ast.Expression) object.Object {
	pexpr := e.Evaluate(expr)
	switch pexpr.(type) {
	case *object.Integer:
		num := pexpr.(*object.Integer).Value
		return &object.Integer{Value: -num}
	case *object.Float:
		num := pexpr.(*object.Float).Value
		return &object.Float{Value: -num}
	default:
		// todo: add error handling
		//! boolean and strings cannot be operated on with a -
		return NULL
	}
}

func (e *Evaluator) evalInfixExpr(expr *ast.InfixExpr) object.Object {
	left := e.Evaluate(expr.Left)
	right := e.Evaluate(expr.Right)

	if ok, res := e.evaluateSides(left, right, expr.Operator); ok {
		return res
	}

	return nil
}

func (e *Evaluator) evaluateSides(left, right object.Object, op string) (bool, object.Object) {
	switch left.(type) {
	case *object.Integer:
		if right, ok := right.(*object.Integer); ok {
			left := left.(*object.Integer)
			return ok, &object.Integer{Value: executeOpInt(left.Value, right.Value, op)}
		}
		return false, nil
	case *object.Float:
		if right, ok := right.(*object.Float); ok {
			if !isValidFltOp(op) {
				return false, nil
			}
			left := left.(*object.Float)
			return ok, &object.Float{Value: executeOpFlt(left.Value, right.Value, op)}
		}
		return false, nil
	case *object.String:
		panic("eval.evaluateSides (case String): unimplemented")
	case *object.Boolean:
		panic("eval.evaluateSides (case boolean): unimplemented")
	case *object.Null:
		panic("eval.evaluateSides (case Null): unimplemented")
	default:
		return false, nil
	}
}

func executeOpInt(left, right int64, op string) int64 {
	switch op {
	case lexer.ADD:
		return left + right
	case lexer.SUB:
		return left - right
	case lexer.MUL:
		return left * right
	case lexer.DIV:
		return left / right
	case lexer.BWAND:
		return left & right
	case lexer.BWOR:
		return left | right
	case lexer.BSL:
		return left << right
	case lexer.BSR:
		return left >> right
	default:
		panic("eval.executeOpInt: reached unreachable code")
	}
}

func executeOpFlt(left, right float64, op string) float64 {
	switch op {
	case lexer.ADD:
		return left + right
	case lexer.SUB:
		return left - right
	case lexer.MUL:
		return left * right
	case lexer.DIV:
		return left / right
	default:
		panic("eval.executeOpFlt: reached unreachable code")
	}
}

func nativeBooltoObj(input bool) object.Object {
	if input {
		return TRUE
	}
	return FALSE
}

func isValidFltOp(input string) bool {
	switch input {
	case lexer.ADD:
		return true
	case lexer.SUB:
		return true
	case lexer.MUL:
		return true
	case lexer.DIV:
		return true
	case lexer.BWAND:
		return false
	case lexer.BWOR:
		return false
	case lexer.BSL:
		return false
	case lexer.BSR:
		return false
	default:
		return false
	}
}

func evaluateTruthiness(in object.Object) bool {
	switch in.(type) {
	case *object.Integer:
		if in.(*object.Integer).Value == 0 {
			return false
		}
		return true
	case *object.Float:
		if in.(*object.Float).Value == 0 {
			return false
		}
		return true
	case *object.String:
		if len(in.(*object.String).Value) == 0 {
			return false
		}
		return true
	case *object.Boolean:
		if in.(*object.Boolean).Value {
			return true
		}
		return false
	default:
		return false
	}
}
