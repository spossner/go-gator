package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/spossner/gator/internal/database"
	"github.com/spossner/gator/internal/rss"
	"strconv"
	"strings"
	"time"
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

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error identifying next feed to fetch: %w", err)
	}

	fmt.Printf("scraping %s...\n", feed.Name)

	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("error marking feed %s as fetched: %w", feed.ID, err)
	}
	rssFeed, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("error fetching feed %s: %w", feed.ID, err)
	}
	for _, item := range rssFeed.Channel.Item {
		pubDate, err := time.Parse("Mon, 02 Jan 2006 15:04:05 +0000", item.PubDate)
		if err != nil {
			return err
		}
		post, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			FeedID:      feed.ID,
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: sql.NullTime{Time: pubDate, Valid: true},
		})
		if err != nil {
			if !strings.HasPrefix(err.Error(), "pq: duplicate key") {
				fmt.Printf("error persisting post %s: %v\n", item.Title, err)
			}
		} else {
			fmt.Printf("* %s\n", post.Title)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	timeBetweenReqs := 30 * time.Second
	if len(cmd.args) > 0 {
		d, err := time.ParseDuration(cmd.args[0])
		if err != nil {
			return fmt.Errorf("invalid duration %s: %w", cmd.args[0], err)
		}
		timeBetweenReqs = d
	}

	if timeBetweenReqs < 5*time.Second {
		return fmt.Errorf("interval %s is too fast - min 5 seconds", timeBetweenReqs)
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)
	defer ticker.Stop()
	for range ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			fmt.Println(err)
		}
	}
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
		return fmt.Errorf("unknown feed %s: %w", url, err)
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
		fmt.Printf("* %s: %s\n", follow.FeedName, follow.FeedUrl)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return errors.New("missing feed url")
	}

	url := cmd.args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error fetching feed %s: %w", url, err)
	}

	follow, err := s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error removing the following of feed %s: %w", url, err)
	}

	fmt.Printf("%s is not following %s anymore\n", follow.UserName, follow.FeedName)

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int32 = 2
	if len(cmd.args) > 0 {
		if i, err := strconv.Atoi(cmd.args[0]); err == nil {
			limit = int32(i)
		}
	}
	posts, err := s.db.GetPostsByUser(context.Background(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("error fetching posts for user %s: %w", user.Name, err)
	}
	for _, post := range posts {
		fmt.Printf("** %s **\n%s\n%v\n---\n%s\n\n", strings.ToUpper(post.Name), post.Title, post.PublishedAt.Time.Format(time.DateTime), strings.TrimSpace(post.Description.String))
	}
	return nil
}
