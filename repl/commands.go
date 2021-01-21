package repl

import (
	"fmt"
	"os"
)

// Commands - the commands that the user can enter in the repl
var Commands = map[string]func(){
	"quit": quit,
	"exit": quit,
	"help": help,
}

func quit() {
	os.Exit(0)
}

func help() {
	fmt.Println(
		`List of Commands:
	:quit    - Leave the REPL
	:exit    - Alias for quit
	:prompt  - Set the prompt
	:engine  - Set the engine used for evaluation
	`)
}

func tokenizeCommand(cmd string) []string {
	tokens := []string{}

	pos := 0

	for i, char := range cmd {
		switch char {
		case ' ':
			tokens = append(tokens, cmd[pos:i])
			pos = i + 1

		}
	}

	return tokens
}
