package eval

import (
	"fmt"

	"github.com/cartoon-raccoon/monkey-jit/ast"
	"github.com/cartoon-raccoon/monkey-jit/lexer"
	"github.com/cartoon-raccoon/monkey-jit/object"
)

func (e *Evaluator) evalPrefixExpr(expr *ast.PrefixExpr) (object.Object, error) {
	switch expr.Operator {
	case "!":
		return e.evalBangPExpr(expr.Right)
	case "-":
		return e.evalMinusPExpr(expr.Right)
	default:
		return NULL, Err{"evalPrefixExpr: unreachable", expr.Token.Pos}
	}
}

func (e *Evaluator) evalBangPExpr(expr ast.Expression) (object.Object, error) {
	pexpr, err := e.Evaluate(expr)
	if err != nil {
		return NULL, err
	}
	truth := evaluateTruthiness(pexpr)
	return nativeBooltoObj(!truth), nil
}

func (e *Evaluator) evalMinusPExpr(expr ast.Expression) (object.Object, error) {
	pexpr, err := e.Evaluate(expr)
	if err != nil {
		return NULL, err
	}
	switch pexpr.(type) {
	case *object.Integer:
		num := pexpr.(*object.Integer).Value
		return &object.Integer{Value: -num}, nil
	case *object.Float:
		num := pexpr.(*object.Float).Value
		return &object.Float{Value: -num}, nil
	default:
		return NULL, Err{"MINUS usage on boolean or string", expr.Context()}
	}
}

func (e *Evaluator) evalInfixExpr(expr *ast.InfixExpr) (object.Object, error) {
	left, err1 := e.Evaluate(expr.Left)
	right, err2 := e.Evaluate(expr.Right)

	if err1 != nil {
		return NULL, Err{"Could not evaluate LHS", expr.Context()}
	}
	if err2 != nil {
		return NULL, Err{"Could not evaluate RHS", expr.Context()}
	}

	if isComparisonOp(expr.Operator) {
		return e.evaluateComp(left, right, expr.Operator, expr.Context())
	}
	return e.evaluateSides(left, right, expr.Operator, expr.Context())
}

func (e *Evaluator) evaluateComp(left, right object.Object, op string, con lexer.Context) (object.Object, error) {
	switch left.(type) {
	case *object.Integer:
		if right, ok := right.(*object.Integer); ok {
			left := left.(*object.Integer)
			return nativeBooltoObj(executeCompInt(left.Value, right.Value, op)), nil
		}
		return NULL, Err{fmt.Sprintf("Cannot compare INT and %T", right), con}
	case *object.Float:
		if right, ok := right.(*object.Float); ok {
			left := left.(*object.Float)
			return nativeBooltoObj(executeCompFlt(left.Value, right.Value, op)), nil
		}
		return NULL, Err{fmt.Sprintf("Cannot compare FLT and %T", right), con}
	case *object.String:
		if right, ok := right.(*object.String); ok {
			left := left.(*object.String)
			return nativeBooltoObj(executeCompStr(left.Value, right.Value, op)), nil
		}
		return NULL, Err{fmt.Sprintf("Cannot compare STR and %T", right), con}
	case *object.Boolean:
		if !isValidBoolOp(op) {
			return NULL, Err{fmt.Sprintf("Cannot use operator `%s` with BOOL", op), con}
		}
		if right, ok := right.(*object.Boolean); ok {
			left := left.(*object.Boolean)
			return nativeBooltoObj(executeCompBool(left.Value, right.Value, op)), nil
		}
		return NULL, Err{fmt.Sprintf("Cannot compare BOOL and %T", right), con}
	case *object.Null:
		return NULL, Err{"LHS is null", con}
	default:
		return NULL, Err{"e.evaluateComp: Unreachable", con}
	}
}

func (e *Evaluator) evaluateSides(left, right object.Object, op string, con lexer.Context) (object.Object, error) {
	switch left.(type) {
	case *object.Integer:
		if right, ok := right.(*object.Integer); ok {
			left := left.(*object.Integer)
			return &object.Integer{Value: executeOpInt(left.Value, right.Value, op)}, nil
		}
		return NULL, Err{fmt.Sprintf("Cannot operate on INT and %T", right), con}
	case *object.Float:
		if right, ok := right.(*object.Float); ok {
			if !isValidFltOp(op) {
				return NULL, Err{fmt.Sprintf("Cannot use operator `%s` on FLT", op), con}
			}
			left := left.(*object.Float)
			return &object.Float{Value: executeOpFlt(left.Value, right.Value, op)}, nil
		}
		return NULL, Err{fmt.Sprintf("Cannot operate on FLT and %T", right), con}
	case *object.String:
		panic("eval.evaluateSides (case String): unimplemented")
	case *object.Boolean:
		return NULL, Err{fmt.Sprintf("Cannot operate `%s` on BOOL", op), con}
	case *object.Null:
		return NULL, Err{"LHS is null", con}
	default:
		return NULL, Err{"e.evaluateSides: Unreacheable", con}
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
