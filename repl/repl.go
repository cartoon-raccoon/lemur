package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/cartoon-raccoon/lemur/compiler"
	"github.com/cartoon-raccoon/lemur/lexer"
	"github.com/cartoon-raccoon/lemur/parser"
	"github.com/cartoon-raccoon/lemur/vm"
)

// PROMPT - the prompt that the user sees
const PROMPT = "lemur >> "

// CONT - When a construct is incomplete
const CONT = "> "

var nestingLevel = 0

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

	fmt.Printf("Lemur Interactive Shell v0.1 ")
	fmt.Printf("(%s %s)\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Type :help for a list of commands\n")
	fmt.Printf("Welcome, %s\n", username)

	// env := object.NewEnv()
	// e := eval.New()
	c := compiler.New()
	vm := vm.New()

	for {
		prompt := PROMPT
		var line string
		fmt.Printf(prompt)

		for {
			scanned := scanner.Scan()
			if !scanned {
				return
			}

			line += scanner.Text() + "\n"

			if isComplete(line) {
				nestingLevel = 0
				break
			} else {
				nestingLevel = 0
				prompt = CONT
				fmt.Printf(prompt)
			}
		}

		if strings.HasPrefix(line, ":") {
			if fn, ok := Commands[strings.TrimSpace(line[1:])]; ok {
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

		// res := e.Evaluate(prog, env)

		// res.Display()

		err = c.Compile(prog)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			continue
		}

		err = vm.Run(c.Bytecode())
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			continue
		}

		vm.LastPopped().Display()

	}
}

// todo: add a stack to track whether the matching bracket is correct
func isComplete(input string) bool {
	for _, char := range input {
		switch char {
		case '[':
			nestingLevel++
		case '{':
			nestingLevel++
		case ']':
			nestingLevel--
		case '}':
			nestingLevel--
		}
	}

	if nestingLevel != 0 {
		return false
	}

	return true
}
