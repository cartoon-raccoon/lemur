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
func (e *Evaluator) Evaluate(node ast.Node) (object.Object, error) {
	switch node.(type) {
	case *ast.Program:
		res := &object.StmtResults{}
		res.Results = []object.Object{}
		for _, stmt := range node.(*ast.Program).Statements {
			if ret, ok := stmt.(*ast.ReturnStatement); ok {
				return e.Evaluate(ret)
			}
			result, err := e.Evaluate(stmt)
			if err != nil {
				return NULL, err
			}
			res.Results = append(res.Results, result)
		}
		return res, nil
	case ast.Statement:
		stmt := node.(ast.Statement)
		switch node.(ast.Statement).(type) {
		case *ast.LetStatement:
			letstmt := stmt.(*ast.LetStatement)
			val, err := e.Evaluate(letstmt.Value)
			if err != nil {
				return NULL, err
			}
			e.Data[letstmt.Name.String()] = val
			return NULL, nil
		case *ast.ExprStatement:
			expr := stmt.(*ast.ExprStatement)
			return e.Evaluate(expr.Expression)
		case *ast.ReturnStatement:
			retstmt := stmt.(*ast.ReturnStatement)
			res, err := e.Evaluate(retstmt.Value)
			if err != nil {
				return NULL, err
			}
			return &object.Return{Inner: res}, nil
		case *ast.BlockStatement:
			blkstmt := stmt.(*ast.BlockStatement)
			return e.evalBlockStmt(blkstmt)
		default:
			return NULL, nil
		}
	case ast.Expression:
		switch node.(ast.Expression).(type) {
		case *ast.Identifier:
			ident := node.(ast.Expression).(*ast.Identifier)
			data, ok := e.Data[ident.Value]
			if !ok {
				return NULL, Err{
					Msg: "Variable not yet declared",
				}
			}
			return data, nil
		case *ast.PrefixExpr:
			pexpr := node.(ast.Expression).(*ast.PrefixExpr)
			return e.evalPrefixExpr(pexpr)
		case *ast.InfixExpr:
			iexpr := node.(ast.Expression).(*ast.InfixExpr)
			return e.evalInfixExpr(iexpr)
		case *ast.IfExpression:
			ifexpr := node.(ast.Expression).(*ast.IfExpression)
			condition, err := e.Evaluate(ifexpr.Condition)
			if err != nil {
				return NULL, err
			}
			if condition == nil {
				return NULL, Err{"If condition returned nil", ifexpr.Context()}
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
					return NULL, Err{"Invalid else branch", ifexpr.Alternative.Context()}
				}
			}
		case *ast.FnLiteral:
			//todo
		case *ast.FunctionCall:
			//todo
		case *ast.DotExpression:
			//todo
		case *ast.Int:
			intexpr := node.(ast.Expression).(*ast.Int)
			return &object.Integer{Value: intexpr.Inner}, nil
		case *ast.Flt:
			fltexpr := node.(ast.Expression).(*ast.Flt)
			return &object.Float{Value: fltexpr.Inner}, nil
		case *ast.Str:
			strexpr := node.(ast.Expression).(*ast.Str)
			return &object.String{Value: strexpr.Inner}, nil
		case *ast.Bool:
			boolexpr := node.(ast.Expression).(*ast.Bool)
			return nativeBooltoObj(boolexpr.Inner), nil
		default:
			return NULL, nil
		}
	default:
		return NULL, Err{"Unimplemented type", node.Context()}
	}
	return NULL, Err{"Evaluate: unreachable code", node.Context()}
}

func (e *Evaluator) evalProgram(prog *ast.Program) (object.Object, error) {
	var result object.Object

	for _, stmt := range prog.Statements {
		result, err := e.Evaluate(stmt)
		if err != nil {
			return NULL, err
		}

		if returnVal, ok := result.(*object.Return); ok {
			return returnVal.Inner, nil
		}
	}

	return result, nil
}

func (e *Evaluator) evalBlockStmt(stmt *ast.BlockStatement) (object.Object, error) {
	var result object.Object

	for _, stmt := range stmt.Statements {
		result, err := e.Evaluate(stmt)
		if err != nil {
			return NULL, err
		}

		if !object.IsNull(result) && result.Type() == object.RETURN {
			return result, nil
		}
	}
	return result, nil
}

// Err - Error returned by the evaluator
type Err struct {
	Msg string
	Con lexer.Context
}

func (e Err) Error() string {
	return fmt.Sprintf("%s - Line %d, Col %d", e.Msg, e.Con.Line, e.Con.Col)
}
