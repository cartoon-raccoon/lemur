package repl

import (
	"bufio"
	"fmt"
	"io"
	"runtime"

	"github.com/cartoon-raccoon/monkey-jit/lexer"
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
	fmt.Printf("running on %s %s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Welcome to Monkey, %s\n", username)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		for tok, err := l.NextToken(); tok.Type != lexer.EOF; tok, err = l.NextToken() {
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Printf("%+v\n", tok)
		}
	}

}
