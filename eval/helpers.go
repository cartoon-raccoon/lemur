package eval

import (
	"fmt"

	"github.com/cartoon-raccoon/monkey-jit/ast"
	"github.com/cartoon-raccoon/monkey-jit/lexer"
	"github.com/cartoon-raccoon/monkey-jit/object"
)

func (e *Evaluator) evalPrefixExpr(expr *ast.PrefixExpr, env *object.Environment) object.Object {
	switch expr.Operator {
	case "!":
		return e.evalBangPExpr(expr.Right, env)
	case "-":
		return e.evalMinusPExpr(expr.Right, env)
	default:
		return &object.Exception{Msg: "evalPrefixExpr: unreachable", Con: expr.Token.Pos}
	}
}

func (e *Evaluator) evalBangPExpr(expr ast.Expression, env *object.Environment) object.Object {
	pexpr := e.Evaluate(expr, env)

	truth := evaluateTruthiness(pexpr)
	return nativeBooltoObj(!truth)
}

func (e *Evaluator) evalMinusPExpr(expr ast.Expression, env *object.Environment) object.Object {
	pexpr := e.Evaluate(expr, env)

	switch pexpr.(type) {
	case *object.Integer:
		num := pexpr.(*object.Integer).Value
		return &object.Integer{Value: -num}
	case *object.Float:
		num := pexpr.(*object.Float).Value
		return &object.Float{Value: -num}
	default:
		return &object.Exception{Msg: "MINUS usage on boolean or string", Con: expr.Context()}
	}
}

func (e *Evaluator) evalInfixExpr(expr *ast.InfixExpr, env *object.Environment) object.Object {
	left := e.Evaluate(expr.Left, env)
	right := e.Evaluate(expr.Right, env)

	if left == nil {
		return &object.Exception{Msg: "Could not evaluate LHS", Con: expr.Context()}
	}
	if right == nil {
		return &object.Exception{Msg: "Could not evaluate RHS", Con: expr.Context()}
	}

	if isComparisonOp(expr.Operator) {
		return e.evaluateComp(left, right, expr.Operator, expr.Context())
	}
	return e.evaluateSides(left, right, expr.Operator, expr.Context())
}

func (e *Evaluator) evaluateComp(left, right object.Object, op string, con lexer.Context) object.Object {
	switch left.(type) {
	case *object.Integer:
		if right, ok := right.(*object.Integer); ok {
			left := left.(*object.Integer)
			return nativeBooltoObj(executeCompInt(left.Value, right.Value, op))
		}
		return &object.Exception{
			Msg: fmt.Sprintf("Cannot compare INT and %T", right),
			Con: con,
		}
	case *object.Float:
		if right, ok := right.(*object.Float); ok {
			left := left.(*object.Float)
			return nativeBooltoObj(executeCompFlt(left.Value, right.Value, op))
		}
		return &object.Exception{
			Msg: fmt.Sprintf("Cannot compare FLT and %T", right),
			Con: con,
		}
	case *object.String:
		if right, ok := right.(*object.String); ok {
			left := left.(*object.String)
			return nativeBooltoObj(executeCompStr(left.Value, right.Value, op))
		}
		return &object.Exception{
			Msg: fmt.Sprintf("Cannot compare STR and %T", right),
			Con: con,
		}
	case *object.Boolean:
		if !isValidBoolOp(op) {
			return &object.Exception{
				Msg: fmt.Sprintf("Cannot use operator `%s` with BOOL", op),
				Con: con,
			}
		}
		if right, ok := right.(*object.Boolean); ok {
			left := left.(*object.Boolean)
			return nativeBooltoObj(executeCompBool(left.Value, right.Value, op))
		}
		return &object.Exception{
			Msg: fmt.Sprintf("Cannot compare BOOL and %T", right),
			Con: con,
		}
	case *object.Null:
		return &object.Exception{
			Msg: "LHS is null",
			Con: con,
		}
	default:
		return &object.Exception{
			Msg: "e.evaluateComp: Unreachable",
			Con: con,
		}
	}
}

func (e *Evaluator) evaluateSides(left, right object.Object, op string, con lexer.Context) object.Object {
	switch left.(type) {
	case *object.Integer:
		if right, ok := right.(*object.Integer); ok {
			left := left.(*object.Integer)
			return &object.Integer{Value: executeOpInt(left.Value, right.Value, op)}
		}
		return &object.Exception{
			Msg: fmt.Sprintf("Cannot operate on INT and %T", right),
			Con: con,
		}
	case *object.Float:
		if right, ok := right.(*object.Float); ok {
			if !isValidFltOp(op) {
				return &object.Exception{
					Msg: fmt.Sprintf("Cannot use operator `%s` on FLT", op),
					Con: con,
				}
			}
			left := left.(*object.Float)
			return &object.Float{Value: executeOpFlt(left.Value, right.Value, op)}
		}
		return &object.Exception{
			Msg: fmt.Sprintf("Cannot operate on FLT and %T", right),
			Con: con,
		}
	case *object.String:
		panic("eval.evaluateSides (case String): unimplemented")
	case *object.Boolean:
		return &object.Exception{
			Msg: fmt.Sprintf("Cannot operate `%s` on BOOL", op),
			Con: con,
		}
	case *object.Null:
		return &object.Exception{
			Msg: "LHS is null",
			Con: con,
		}
	default:
		return &object.Exception{
			Msg: "e.evaluateSides: Unreacheable",
			Con: con,
		}
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
