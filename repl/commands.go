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
	fmt.Println("Sorry, help is not implemented at this time")
}
