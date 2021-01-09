package eval

import (
	"github.com/cartoon-raccoon/monkey-jit/ast"
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
		ret := &object.StmtResults{}
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
			return NULL
		case *ast.ExprStatement:
			expr := node.(ast.Statement).(*ast.ExprStatement)
			return e.Evaluate(expr.Expression)
		case *ast.ReturnStatement:

		case *ast.BlockStatement:
			ret := &object.StmtResults{}
			ret.Results = []object.Object{}
			blkstmt := node.(ast.Statement).(*ast.BlockStatement)
			for _, stmt := range blkstmt.Statements {
				ret.Results = append(ret.Results, e.Evaluate(stmt))
			}
			return ret
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
			if condition == nil {
				//todo: return error
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
					//todo: throw error
					return nil
				}
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
