package eval

import (
	"github.com/cartoon-raccoon/monkey-jit/ast"
	"github.com/cartoon-raccoon/monkey-jit/lexer"
	"github.com/cartoon-raccoon/monkey-jit/object"
)

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

	if isComparisonOp(expr.Operator) {
		return e.evaluateComp(left, right, expr.Operator)
	}

	if ok, res := e.evaluateSides(left, right, expr.Operator); ok {
		return res
	}

	return nil
}

func (e *Evaluator) evaluateComp(left, right object.Object, op string) object.Object {
	switch left.(type) {
	case *object.Integer:
		if right, ok := right.(*object.Integer); ok {
			left := left.(*object.Integer)
			return nativeBooltoObj(executeCompInt(left.Value, right.Value, op))
		}
		//todo: error
		return nil
	case *object.Float:
		if right, ok := right.(*object.Float); ok {
			left := left.(*object.Float)
			return nativeBooltoObj(executeCompFlt(left.Value, right.Value, op))
		}
		return nil
	case *object.String:
		if right, ok := right.(*object.String); ok {
			left := left.(*object.String)
			return nativeBooltoObj(executeCompStr(left.Value, right.Value, op))
		}
		return nil
	case *object.Boolean:
		if !isValidBoolOp(op) {
			//todo: throw error
			return nil
		}
		if right, ok := right.(*object.Boolean); ok {
			left := left.(*object.Boolean)
			return nativeBooltoObj(executeCompBool(left.Value, right.Value, op))
		}
		//todo: error
		return nil
	case *object.Null:
		panic("eval.evaluateSides (case Null): unimplemented")
	default:
		return nil
	}
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

func executeCompInt(left, right int64, op string) bool {
	switch op {
	case lexer.EQ:
		return left == right
	case lexer.NE:
		return left != right
	case lexer.GT:
		return left > right
	case lexer.LT:
		return left < right
	case lexer.GE:
		return left >= right
	case lexer.LE:
		return left <= right
	default:
		panic("eval.executeCompInt: reached unreachable code")
	}
}

func executeCompFlt(left, right float64, op string) bool {
	switch op {
	case lexer.EQ:
		return left == right
	case lexer.NE:
		return left != right
	case lexer.GT:
		return left > right
	case lexer.LT:
		return left < right
	case lexer.GE:
		return left >= right
	case lexer.LE:
		return left <= right
	default:
		panic("eval.executeCompFlt: reached unreachable code")
	}
}

func executeCompStr(left, right string, op string) bool {
	switch op {
	case lexer.EQ:
		return left == right
	case lexer.NE:
		return left != right
	case lexer.GT:
		return left > right
	case lexer.LT:
		return left < right
	case lexer.GE:
		return left >= right
	case lexer.LE:
		return left <= right
	default:
		panic("eval.executeCompStr: reached unreachable code")
	}
}

func executeCompBool(left, right bool, op string) bool {
	switch op {
	case lexer.EQ:
		return left == right
	case lexer.NE:
		return left != right
	default:
		panic("eval.executeCompStr: reached unreachable code")
	}
}

func isValidBoolOp(op string) bool {
	switch op {
	case lexer.EQ:
		return true
	case lexer.NE:
		return true
	default:
		return false
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

func isComparisonOp(input string) bool {
	switch input {
	case lexer.EQ:
		return true
	case lexer.NE:
		return true
	case lexer.LT:
		return true
	case lexer.GT:
		return true
	case lexer.LE:
		return true
	case lexer.GE:
		return true
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

func nativeBooltoObj(input bool) object.Object {
	if input {
		return TRUE
	}
	return FALSE
}
