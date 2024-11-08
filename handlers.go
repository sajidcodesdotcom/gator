package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sajidcodess/gator/internal/database"
)

func loginHandler(state *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage %s <name>", cmd.Name)
	}
	username := cmd.Args[0]
	dbUser, err := state.db.GetUser(context.Background(), username)
	if err != nil && dbUser.Name != username {
		return fmt.Errorf("User not found, please register before loging in: %s", err)
	}

	if err := state.cfg.SetUser(username); err != nil {
		return fmt.Errorf("the user %s is not found: %w", username, err)
	}
	fmt.Println("The user has been set")
	return nil
}

func registerHandler(state *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage %s <name>", cmd.Name)
	}
	username := cmd.Args[0]
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	dbUser, err := state.db.GetUser(context.Background(), username)
	if err != nil && dbUser.Name == username {
		return fmt.Errorf("User already exist: %s", err)
	}

	newUser, err := state.db.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Println("the user is created")

	if err := state.cfg.SetUser(newUser.Name); err != nil {
		return fmt.Errorf("the user %s is not found: %w", username, err)
	}

	return nil
}

func resetHandler(state *state, cmd command) error {
	err := state.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Deleting/reseting users table failed: %w", err)
	}
	fmt.Println("The users table has been reset successfully.")
	return nil
}

func getUsersHandler(state *state, cmd command) error {
	users, err := state.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error while getting all the users list: %w", err)
	}
	for _, val := range users {
		if state.cfg.CurrentUserName == val.Name {
		fmt.Println("* " + val.Name + " (current)")
		} else {
		fmt.Println("* " + val.Name)
    }
	}
	return nil
}
