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
	Data map[string]object.Object
}

// New - returns a new evaluator
func New() *Evaluator {
	eval := &Evaluator{
		Data: map[string]object.Object{},
	}
	return eval
}

// Evaluate runs the evaluator, walking the tree and executing code
func (e *Evaluator) Evaluate(node ast.Node) object.Object {
	switch node.(type) {
	case *ast.Program:
		res := &object.StmtResults{}
		res.Results = []object.Object{}
		for _, stmt := range node.(*ast.Program).Statements {
			if ret, ok := stmt.(*ast.ReturnStatement); ok {
				return e.Evaluate(ret)
			}
			result := e.Evaluate(stmt)
			res.Results = append(res.Results, result)
		}
		return res
	case ast.Statement:
		stmt := node.(ast.Statement)
		switch node.(ast.Statement).(type) {
		case *ast.LetStatement:
			letstmt := stmt.(*ast.LetStatement)
			val := e.Evaluate(letstmt.Value)
			e.Data[letstmt.Name.String()] = val
			return NULL
		case *ast.ExprStatement:
			expr := stmt.(*ast.ExprStatement)
			return e.Evaluate(expr.Expression)
		case *ast.ReturnStatement:
			retstmt := stmt.(*ast.ReturnStatement)
			res := e.Evaluate(retstmt.Value)
			return &object.Return{Inner: res}
		case *ast.BlockStatement:
			blkstmt := stmt.(*ast.BlockStatement)
			return e.evalBlockStmt(blkstmt)
		default:
			return NULL
		}
	case ast.Expression:
		switch node.(ast.Expression).(type) {
		case *ast.Identifier:
			ident := node.(ast.Expression).(*ast.Identifier)
			data, ok := e.Data[ident.Value]
			if !ok {
				return &object.Exception{
					Msg: "Variable not yet declared",
					Con: ident.Context(),
				}
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
			if condition == nil {
				return &object.Exception{
					Msg: "If condition returned nil",
					Con: ifexpr.Context(),
				}
			}
			if evaluateTruthiness(condition) {
				return e.Evaluate(ifexpr.Result)
			}
			if ifexpr.Alternative != nil {
				switch ifexpr.Alternative.(type) {
				case *ast.BlockStatement:
					return e.Evaluate(ifexpr.Alternative.(*ast.BlockStatement))
				case *ast.IfExpression:
					return e.Evaluate(ifexpr.Alternative.(*ast.IfExpression))
				default:
					return &object.Exception{
						Msg: "Invalid else branch",
						Con: ifexpr.Alternative.Context(),
					}
				}
			}
		case *ast.FnLiteral:
			//todo
			return &object.Exception{
				Msg: "FnLiteral: unimplemented",
				Con: node.Context(),
			}
		case *ast.FunctionCall:
			//todo
			return &object.Exception{
				Msg: "FuncCall: unimplemented",
				Con: node.Context(),
			}
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

func (e *Evaluator) evalProgram(prog *ast.Program) (object.Object, error) {
	var result object.Object

	for _, stmt := range prog.Statements {
		result := e.Evaluate(stmt)

		if returnVal, ok := result.(*object.Return); ok {
			return returnVal.Inner, nil
		}
	}

	return result, nil
}

func (e *Evaluator) evalBlockStmt(stmt *ast.BlockStatement) object.Object {
	var result object.Object

	for _, stmt := range stmt.Statements {
		result := e.Evaluate(stmt)
		if !object.IsNull(result) && result.Type() == object.RETURN {
			return result
		}
	}
	return result
}

// Err - Error returned by the evaluator
type Err struct {
	Msg string
	Con lexer.Context
}

func (e Err) Error() string {
	return fmt.Sprintf("%s - Line %d, Col %d", e.Msg, e.Con.Line, e.Con.Col)
}
