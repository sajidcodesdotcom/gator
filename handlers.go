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

func aggHandler(state *state, cmd command) error {
	// if len(cmd.Args) != 1 {
	// 	return fmt.Errorf("usage %s <URL>", cmd.Name)
	// }
	// url := cmd.Args[0]
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
}

func addFeedHandler(state *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage %s <name> <feedURL>", cmd.Name)
	}
	// get args [name] and URL and save it in the feeds table.
	currentUser, err := state.db.GetUser(context.Background(), state.cfg.CurrentUserName)
	if err != nil {
		return err
	}
	name := cmd.Args[0]
	url := cmd.Args[1]
	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    currentUser.ID,
	}
	feed, err := state.db.CreateFeed(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Error, while creating a feed: %w", err)
	}
	fmt.Println("The feed has been successfully created")
	feedFollow, err := state.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(feed, currentUser)
	fmt.Println()
	fmt.Println("Feed followed successfully:")
	printFeedFollow(feedFollow.UserName, feedFollow.FeedName)
	fmt.Println("=====================================")
	return nil

}

func listFeeds(state *state, cmd command) error {
	feeds, err := state.db.ListFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error while listing the feeds: %w", err)
	}
	if len(feeds) == 0 {
		fmt.Println("0 feeds found")
		return nil
	}
	fmt.Println("List of feeds stored in the DB")
	fmt.Println("=================================")

	for i, feed := range feeds {
		fmt.Println(i+1, "___________________")
		fmt.Printf("* Feed Name:       %s\n", feed.Feedname)
		fmt.Printf("* Feed URL:       %s\n", feed.Url)
		fmt.Printf("* Created By:       %s\n", feed.Username)

	}

	return nil
}

func followHandler(state *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage %s <name>", cmd.Name)
	}
	url := cmd.Args[0]
	user, err := state.db.GetUser(context.Background(), state.cfg.CurrentUserName)
	if err != nil {
		return err
	}
	feed, err := state.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	feedFollow, err := state.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Couldn't create feed follow: %w", err)
	}

	fmt.Println("Feed is created successfully")
	fmt.Println("Feed is created successfully")
	fmt.Printf("UserName: %v, FeedName: %v", feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func followingHandler(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get feed follows: %w", err)
	}

	if len(feedFollows) == 0 {
		fmt.Println("No feed follows found for this user.")
		return nil
	}

	fmt.Printf("Feed follows for user %s:\n", user.Name)
	for _, ff := range feedFollows {
		fmt.Printf("* %s\n", ff.FeedName)
	}

	return nil
}

func printFeedFollow(username, feedname string) {
	fmt.Printf("* User:          %s\n", username)
	fmt.Printf("* Feed:          %s\n", feedname)
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* User:          %s\n", user.Name)
}
