package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/spossner/gator/internal/database"
	"github.com/spossner/gator/internal/rss"
)

type command struct {
	name string
	args []string
}

type handler func(*state, command) error
type authenticatedHandler func(*state, command, database.User) error

type commands struct {
	list map[string]handler
}

func (c *commands) register(name string, fn func(*state, command) error) {
	c.list[name] = fn
}

func (c *commands) run(s *state, cmd command) error {
	fn, ok := c.list[cmd.name]
	if !ok {
		return fmt.Errorf("unknwon command %v", cmd.name)
	}
	return fn(s, cmd)
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("missing username")
	}
	name := cmd.args[0]
	user, err := s.db.GetUserByName(context.Background(), name)
	if err != nil {
		return fmt.Errorf("error fetching user %v: %w", name, err)
	}

	if err := s.cfg.SetUser(user.Name); err != nil {
		return fmt.Errorf("error logging in: %w", err)
	}
	fmt.Printf("User %v logged in successfully\n", s.cfg.CurrentUserName)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("missing username")
	}
	user, err := s.db.CreateUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	if err := s.cfg.SetUser(user.Name); err != nil {
		return fmt.Errorf("error logging in %v after registration: %w", user.Name, err)
	}

	fmt.Printf("User %s generated and successfully logged in\n", user.Name)
	fmt.Println(user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.Reset(context.Background()); err != nil {
		return fmt.Errorf("error wiping users table: %w", err)
	}
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error retrieving users: %w", err)
	}
	for _, user := range users {
		fmt.Printf("* %s", user.Name)
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf(" (current)")
		}
		fmt.Println()
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("error reading rss feed: %w", err)
	}
	fmt.Println(feed)
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("missing name + url")
	}
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		Name:   cmd.args[0],
		Url:    cmd.args[1],
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed entry: %w", err)
	}

	fmt.Printf("created feed %s\n", feed.Name)

	follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error following the new feed %s: %w", feed.Name, err)
	}

	fmt.Printf("%s is now following %s\n", follow.UserName, follow.FeedName)

	return nil
}

func handlerFeeds(s *state, _ command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching feeds: %w", err)
	}
	for _, feed := range feeds {
		fmt.Printf("* %s (%s), %s\n", feed.Name, feed.Url, feed.UserName)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return errors.New("missing feed url")
	}

	url := cmd.args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error fetching feed %s: %w", url, err)
	}

	follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating entry to follow the feed: %w", err)
	}

	fmt.Printf("%s is now following %s\n", follow.UserName, follow.FeedName)

	return nil
}

func handlerFollowing(s *state, _ command, user database.User) error {
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error fetching follows for user %s: %w", user.Name, err)
	}

	for _, follow := range follows {
		fmt.Printf("* %s\n", follow.FeedName)
	}

	return nil
}
