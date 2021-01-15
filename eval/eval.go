package eval

import (
	"fmt"
	"os"

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

var builtins = map[string]*object.Builtin{
	// Gets the number of items in a collection
	// Can be called on strings, arrays and maps.
	"len": {
		Fn: func(ctxt lexer.Context, args ...object.Object) object.Object {
			if len := len(args); len != 1 {
				return &object.Exception{
					Msg: fmt.Sprintf("Expected 1 argument, got %d", len),
					Con: ctxt,
				}
			}
			arg := args[0]
			switch args[0].(type) {
			case *object.String:
				str := arg.(*object.String)
				return &object.Integer{Value: int64(len(str.Value))}
			case *object.Array:
				arr := arg.(*object.Array)
				return &object.Integer{Value: int64(len(arr.Elements))}
			default:
				return &object.Exception{
					Msg: fmt.Sprintf("Cannot use type %T as argument for len()", arg),
					Con: ctxt,
				}
			}
		},
	},
	// Variadic, pushes all its arguments onto the first argument which must be a collection
	// Can be called on arrays and maps
	// When called on a map, the arguments must be array of length 2
	//todo: change this to a method instead of a function once dot expressions are implemented
	"push": {
		Fn: func(ctxt lexer.Context, args ...object.Object) object.Object {
			arr, ok := args[0].(*object.Array)
			if !ok {
				return &object.Exception{
					Msg: fmt.Sprintf("Cannot push to type %T", args[0]),
					Con: ctxt,
				}
			}
			for _, elem := range args[1:] {
				arr.Elements = append(arr.Elements, elem)
			}
			return arr
		},
	},
	"first": {
		Fn: func(ctxt lexer.Context, args ...object.Object) object.Object {
			if len := len(args); len != 1 {
				return &object.Exception{
					Msg: fmt.Sprintf("Expected 1 argument for call to first(), got %d", len),
					Con: ctxt,
				}
			}
			arr, ok := args[0].(*object.Array)
			if !ok {
				return &object.Exception{
					Msg: fmt.Sprintf("Cannot call first() on type %T", arr),
					Con: ctxt,
				}
			}

			if len(arr.Elements) < 1 {
				return NULL
			}

			return arr.Elements[0]
		},
	},
	"print": {
		Fn: func(ctxt lexer.Context, args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect())
				fmt.Print(" ")
			}
			fmt.Print("\n")
			return NULL
		},
	},
	"quit": { // for exiting normally
		Fn: func(ctxt lexer.Context, args ...object.Object) object.Object {
			if len := len(args); len != 0 {
				return &object.Exception{
					Msg: fmt.Sprintf(`Expected 0 arguments for quit(), got %d
					Use exit() to exit with a status code`, len),
					Con: ctxt,
				}
			}
			os.Exit(0)
			return NULL
		},
	},
	"exit": { // for exiting with a code
		Fn: func(ctxt lexer.Context, args ...object.Object) object.Object {
			if len := len(args); len != 1 {
				return &object.Exception{
					Msg: fmt.Sprintf("Expected 1 argument for exit(), got %d", len),
					Con: ctxt,
				}
			}
			arg := args[0]
			code, ok := arg.(*object.Integer)
			if !ok {
				return &object.Exception{
					Msg: fmt.Sprintf("Cannot use %T as argument in exit()", code),
					Con: ctxt,
				}
			}
			os.Exit(int(code.Value))
			return NULL
		},
	},
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

		// adding statements
		for _, stmt := range node.(*ast.Program).Statements {
			if ret, ok := stmt.(*ast.ReturnStatement); ok {
				return e.Evaluate(ret, env)
			}
			result := e.Evaluate(stmt, env)
			res.Results = append(res.Results, result)
		}

		// adding functions
		//todo: this should function differently than closures
		for _, fn := range node.(*ast.Program).Functions {
			body := fn.Body
			params := fn.Params
			env.Data[fn.Name.Value] = &object.Function{
				Params: params,
				Body:   body,
				Env:    env,
			}
		}

		//todo: adding classes

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

		case *ast.WhileStatement:
			whilestmt := stmt.(*ast.WhileStatement)

			var result object.Object

			for {
				val := e.Evaluate(whilestmt.Condition, env)
				if !evaluateTruthiness(val) {
					break
				}
				result = e.evalBlockStmt(whilestmt.Body, env)
				if object.IsErr(result) {
					return result
				}
			}

			return result

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
			if data, ok := env.Data[ident.Value]; ok {
				return data
			}
			if bltn, ok := builtins[ident.Value]; ok {
				return bltn
			}
			return &object.Exception{
				Msg: fmt.Sprintf("Could not find symbol %s", ident.Value),
				Con: ident.Context(),
			}

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
		case *ast.Array:
			array := node.(ast.Expression).(*ast.Array)
			arr := &object.Array{}

			// preallocating so we don't have to waste cycles
			// reallocating every time we append
			elements := make([]object.Object, 0, len(array.Elements))

			for _, elem := range array.Elements {
				elements = append(elements, e.Evaluate(elem, env))
			}
			arr.Elements = elements

			return arr

		case *ast.Map:
			hash := node.(ast.Expression).(*ast.Map)
			newmap := &object.Map{}
			newmap.Elements = make(map[object.HashKey]object.Object)

			for key, val := range hash.Elements {
				nkey, nval := e.Evaluate(key, env), e.Evaluate(val, env)

				if object.IsErr(nkey) {
					return nkey
				}
				if object.IsErr(nval) {
					return nval
				}

				hashable, ok := nkey.(object.Hashable)

				if !ok {
					return &object.Exception{
						Msg: fmt.Sprintf("Cannot use type %T as key for Map", nkey),
						Con: hash.Context(),
					}
				}

				newmap.Elements[hashable.HashKey()] = nval
			}

			return newmap

		case *ast.IndexExpr:
			idx := node.(ast.Expression).(*ast.IndexExpr)
			return e.evalIndexExpr(idx, env)

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
		if builtin, ok := fn.(*object.Builtin); ok {
			return builtin.Fn(e.Ctxt, args...)
		}
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

func (e *Evaluator) evalIndexExpr(idx *ast.IndexExpr, env *object.Environment) object.Object {
	left := e.Evaluate(idx.Left, env)
	index := e.Evaluate(idx.Index, env)

	switch left.(type) {
	case *object.Array:
		left := left.(*object.Array)
		pos, ok := index.(*object.Integer)
		if !ok {
			return &object.Exception{
				Msg: fmt.Sprintf("Cannot index into array with index of type %T", pos),
				Con: idx.Context(),
			}
		}
		if len := len(left.Elements); int(pos.Value) > len-1 {
			return &object.Exception{
				Msg: fmt.Sprintf("Cannot get index %d of array of length %d", pos.Value, len),
				Con: idx.Context(),
			}
		}
		return left.Elements[int(pos.Value)]

	case *object.String:
		left := left.(*object.String)
		pos, ok := index.(*object.Integer)
		if !ok {
			return &object.Exception{
				Msg: fmt.Sprintf("Cannot index into string with index of type %T", pos),
				Con: idx.Context(),
			}
		}
		if len := len(left.Value); int(pos.Value) > len-1 {
			return &object.Exception{
				Msg: fmt.Sprintf("Cannot get index %d of array of length %d", pos.Value, len),
				Con: idx.Context(),
			}
		}
		return &object.String{Value: string(left.Value[int(pos.Value)])}

	case *object.Map:
		left := left.(*object.Map)

		hashable, ok := index.(object.Hashable)
		if !ok {
			return &object.Exception{
				Msg: fmt.Sprintf("Cannot use type %T as key for Map", index),
				Con: idx.Index.Context(),
			}
		}
		ret, ok := left.Elements[hashable.HashKey()]
		if !ok {
			return NULL
		}
		return ret

	default:
		return &object.Exception{
			Msg: fmt.Sprintf("Cannot use type %T as index", left),
			Con: idx.Context(),
		}
	}
}
