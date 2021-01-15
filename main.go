package main

import (
	"os"
	"os/user"

	"github.com/cartoon-raccoon/lemur/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	r := repl.New()
	r.Run(user.Username, os.Stdin, os.Stdout)
}
