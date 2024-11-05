package main

import "fmt"

func loginHandler (state *state, cmd command) error {
  if len(cmd.Args) != 1 {
    return fmt.Errorf("usage %s <name>", cmd.Name)
  }
  username := cmd.Args[0]
  if err := state.cfg.SetUser(username); err != nil {
    return fmt.Errorf("the user %s is not found: %w", username, err)
  }
  fmt.Println("The user has been set")
  return nil
}
