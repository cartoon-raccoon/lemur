package eval

import (
	"fmt"

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
	Ctxt lexer.Context
}

// New - returns a new evaluator
func New() *Evaluator {
	eval := &Evaluator{
		Ctxt: lexer.Context{Line: 1, Col: 1, Ctxt: ""},
	}
	return eval
}

// Evaluate runs the evaluator, walking the tree and executing code
func (e *Evaluator) Evaluate(node ast.Node, env *object.Environment) object.Object {
	e.Ctxt = node.Context()
	switch node.(type) {
	case *ast.Program:
		res := &object.StmtResults{}
		res.Results = []object.Object{}
		for _, stmt := range node.(*ast.Program).Statements {
			if ret, ok := stmt.(*ast.ReturnStatement); ok {
				return e.Evaluate(ret, env)
			}
			result := e.Evaluate(stmt, env)
			res.Results = append(res.Results, result)
		}
		return res
	case ast.Statement:
		stmt := node.(ast.Statement)
		switch node.(ast.Statement).(type) {
		case *ast.LetStatement:
			letstmt := stmt.(*ast.LetStatement)
			val := e.Evaluate(letstmt.Value, env)
			env.Data[letstmt.Name.String()] = val
			return NULL
		case *ast.ExprStatement:
			expr := stmt.(*ast.ExprStatement)
			return e.Evaluate(expr.Expression, env)
		case *ast.ReturnStatement:
			retstmt := stmt.(*ast.ReturnStatement)
			res := e.Evaluate(retstmt.Value, env)
			return &object.Return{Inner: res}
		case *ast.BlockStatement:
			blkstmt := stmt.(*ast.BlockStatement)
			return e.evalBlockStmt(blkstmt, env)
		default:
			return NULL
		}
	case ast.Expression:
		expr := node.(ast.Expression)
		switch node.(ast.Expression).(type) {
		case *ast.Identifier:
			ident := expr.(*ast.Identifier)
			data, ok := env.Data[ident.Value]
			if !ok {
				return &object.Exception{
					Msg: "Variable not yet declared",
					Con: ident.Context(),
				}
			}
			return data
		case *ast.PrefixExpr:
			pexpr := expr.(*ast.PrefixExpr)
			return e.evalPrefixExpr(pexpr, env)
		case *ast.InfixExpr:
			iexpr := expr.(*ast.InfixExpr)
			return e.evalInfixExpr(iexpr, env)
		case *ast.IfExpression:
			ifexpr := expr.(*ast.IfExpression)
			condition := e.Evaluate(ifexpr.Condition, env)
			if condition == nil {
				return &object.Exception{
					Msg: "If condition returned nil",
					Con: ifexpr.Context(),
				}
			}
			if evaluateTruthiness(condition) {
				return e.Evaluate(ifexpr.Result, env)
			}
			if ifexpr.Alternative != nil {
				switch ifexpr.Alternative.(type) {
				case *ast.BlockStatement:
					return e.Evaluate(ifexpr.Alternative.(*ast.BlockStatement), env)
				case *ast.IfExpression:
					return e.Evaluate(ifexpr.Alternative.(*ast.IfExpression), env)
				default:
					return &object.Exception{
						Msg: "Invalid else branch",
						Con: ifexpr.Alternative.Context(),
					}
				}
			}
		case *ast.FnLiteral:
			fnlit := expr.(*ast.FnLiteral)
			params := fnlit.Params
			body := fnlit.Body
			return &object.Function{Params: params, Env: env, Body: body}

		case *ast.FunctionCall:
			// asserting type
			fncall := expr.(*ast.FunctionCall)

			// resolving to object
			function := e.Evaluate(fncall.Ident, env)
			if object.IsErr(function) {
				return function
			}

			args := e.evalExpressions(fncall.Params, env)
			if len(args) == 1 && object.IsErr(args[0]) {
				return args[0]
			}

			return e.applyFunction(function, args)
		case *ast.DotExpression:
			//todo
			return &object.Exception{
				Msg: "DotExpr: unimplemented",
				Con: node.Context(),
			}
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
		return &object.Exception{
			Msg: "Unimplemented type",
			Con: node.Context(),
		}
	}
	return &object.Exception{
		Msg: "Evaluate: unreachable code",
		Con: node.Context(),
	}
}

func (e *Evaluator) evalProgram(prog *ast.Program, env *object.Environment) (object.Object, error) {
	var result object.Object

	for _, stmt := range prog.Statements {
		result := e.Evaluate(stmt, env)

		if returnVal, ok := result.(*object.Return); ok {
			return returnVal.Inner, nil
		}
	}

	return result, nil
}

func (e *Evaluator) evalBlockStmt(stmt *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range stmt.Statements {
		result = e.Evaluate(stmt, env)
		if !object.IsNull(result) && result.Type() == object.RETURN {
			return result
		}
	}
	return result
}

func (e *Evaluator) evalExpressions(
	exprs []ast.Expression,
	env *object.Environment,
) []object.Object {
	objects := []object.Object{}

	for _, expr := range exprs {
		res := e.Evaluate(expr, env)
		if object.IsErr(res) {
			return []object.Object{res}
		}
		objects = append(objects, res)
	}

	return objects
}

func (e *Evaluator) applyFunction(
	fn object.Object,
	args []object.Object,
) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		return &object.Exception{
			Msg: "Not a function",
			Con: e.Ctxt,
		}
	}

	len1 := len(args)
	len2 := len(function.Params)
	if len1 != len2 {
		return &object.Exception{
			Msg: fmt.Sprintf(
				"Param mismatch: expected %d, got %d", len2, len1,
			),
			Con: e.Ctxt,
		}
	}

	extendedEnv := extendFunctionEnv(function, args)
	evaluated := e.Evaluate(function.Body, extendedEnv)

	return unwrapReturnValue(evaluated)
}

func extendFunctionEnv(
	fn *object.Function,
	args []object.Object,
) *object.Environment {
	env := object.NewEnclosedEnv(fn.Env)

	for i, param := range fn.Params {
		env.Set(param.Value, args[i])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if ret, ok := obj.(*object.Return); ok {
		return ret.Inner
	}

	return obj
}
