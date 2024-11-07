package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/sajidcodess/gator/internal/config"
	"github.com/sajidcodess/gator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading the config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
    log.Fatalf("Error while opening DB connection: ", err)
	}
  dbQueries := database.New(db)

	programState := &state{
		cfg: &cfg,
    db: dbQueries,
	}

	cmds := commands{
		registerCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", loginHandler)
  cmds.register("register", registerHandler)

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
