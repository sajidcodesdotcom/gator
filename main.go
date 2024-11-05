package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sajidcodess/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading the config: %v", err)
	}

	programState := &state{
		cfg: &cfg,
	}

	cmds := commands{
		registerCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", loginHandler)

	if len(os.Args) < 2 {
    os.Exit(1)
		fmt.Println("Usage: cli <command> [args...]")
		return
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
