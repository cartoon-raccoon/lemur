package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/cartoon-raccoon/monkey-jit/eval"
	"github.com/cartoon-raccoon/monkey-jit/lexer"
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

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
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
			fmt.Println("Errors encountered 1")
			fmt.Fprintf(os.Stdout, err.Error())
			continue
		}
		prog := p.Parse()
		if p.CheckErrors() != nil {
			fmt.Println("Errors encountered 2")
			for _, err := range p.CheckErrors() {
				fmt.Fprintf(os.Stdout, err.Error())
			}
			continue
		}

		evaluator := eval.Evaluator{}

		obj := evaluator.Evaluate(prog)

		if obj == nil {
			continue
		}
		obj.Display()

	}

}
