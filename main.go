package main

import (
	"github.com/spossner/gator/internal/config"
	"log"
	"os"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal("error reading config file", err)
	}
	s := &state{cfg: cfg}
	cmds := commands{
		list: make(map[string]handler),
	}
	cmds.register("login", handlerLogin)

	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal("Usage...")
	}
	if err := cmds.run(s, command{
		name: args[0],
		args: args[1:],
	}); err != nil {
		log.Fatal("error logging in: ", err)
	}
}
