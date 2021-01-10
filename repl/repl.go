package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/cartoon-raccoon/monkey-jit/eval"
	"github.com/cartoon-raccoon/monkey-jit/lexer"
	"github.com/cartoon-raccoon/monkey-jit/object"
	"github.com/cartoon-raccoon/monkey-jit/parser"
)

// PROMPT - the prompt that the user sees
const PROMPT = ">>> "

// Repl - the interative Monkey shell
type Repl struct {
	lexer lexer.Lexer
}

// New - returns an instance of a Repl
func New() *Repl {
	r := &Repl{}
	return r
}

// Run - Runs the Repl
func (r *Repl) Run(username string, in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	fmt.Printf("Monkey Interactive Shell v0.1\n")
	fmt.Printf("running on (%s %s)\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Welcome, %s\n", username)

	env := object.NewEnv()
	evaluator := eval.New()

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if strings.HasPrefix(line, ":") {
			if fn, ok := Commands[line[1:]]; ok {
				fn()
			} else {
				fmt.Printf("No command `%s` found\n", line[1:])
			}
			continue
		}
		l := lexer.New(line)
		// for tok, err := l.NextToken(); tok.Type != lexer.EOF; tok, err = l.NextToken() {
		// 	if err != nil {
		// 		fmt.Println(err)
		// 		break
		// 	}
		// 	fmt.Printf("%+v\n", tok)
		// }
		p, err := parser.New(l)
		if err != nil {
			fmt.Fprintf(os.Stdout, "%s\n", err.Error())
			continue
		}
		prog := p.Parse()
		if p.CheckErrors() != nil {
			for _, err := range p.CheckErrors() {
				fmt.Fprintf(os.Stdout, "%s\n", err.Error())
			}
			continue
		}

		obj := evaluator.Evaluate(prog, env)

		if obj == nil {
			continue
		}
		obj.Display()

	}

}
