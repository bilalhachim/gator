package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/bilalhachim/gator/internal/database"
	"github.com/google/uuid"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func handlerLogin(s *state, cmd command) error {

	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]
	_, ok := s.db.GetUser(context.Background(), name)
	if ok != nil {
		return fmt.Errorf("you can't login to an account that doesn't exist!: %w", ok)
	}
	err := s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User switched successfully!")
	return nil
}
func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	name := cmd.Args[0]

	dbUser, err := s.db.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: name})
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	s.cfg.SetUser(dbUser.Name)
	return nil
}
func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	return err

}
func handlerusers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for i := 0; i < len(users); i++ {
		if users[i] == s.cfg.CurrentUserName {
			fmt.Println("* " + users[i] + " (current)")
			continue
		}
		fmt.Println("* " + users[i])

	}
	return err
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}
	var c http.Client
	res, err := c.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	resByte, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, err
	}
	var feed RSSFeed
	err = xml.Unmarshal(resByte, &feed)
	if err != nil {
		return &RSSFeed{}, err
	}
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	return &feed, nil
}
func handleragg(s *state, cmd command) error {
	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	fmt.Print("Collecting feeds every")
	fmt.Println(timeBetweenRequests)
	if err != nil {
		return err
	}
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func handlerFeed(s *state, cmd command, user database.User) error {

	if len(cmd.Args) <= 1 {

		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	dbFeed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{FeedID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: name, Url: url, ReferenceID: user.ID})
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	FeedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: user.ID, FeedID: dbFeed.FeedID})
	if err != nil {
		return err
	}
	fmt.Println(FeedFollow)
	return nil
}
func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for i := 0; i < len(feeds); i++ {
		dbUser, _ := s.db.GetUserById(context.Background(), feeds[i].ReferenceID)
		fmt.Println("* name" + feeds[i].Name + " ,url: " + feeds[i].Url + " ,users who created the feed: " + dbUser.Name)

	}
	return err
}
func handlerFollow(s *state, cmd command, user database.User) error {
	url := cmd.Args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}

	FeedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: user.ID, FeedID: feed.FeedID})
	if err != nil {
		return err
	}
	fmt.Print(FeedFollow)
	return nil
}
func handlerFollowing(s *state, cmd command, user database.User) error {

	feedsName, err := s.db.GetFeedFollowsForUser(context.Background(), s.cfg.CurrentUserName)
	for i := 0; i < len(feedsName); i++ {
		fmt.Println(feedsName[i])

	}

	if err != nil {
		return err
	}
	return nil

}
func handlerUnfollow(s *state, cmd command, user database.User) error {
	url := cmd.Args[0]
	err := s.db.DeleteFeedFollowsForUser(context.Background(), database.DeleteFeedFollowsForUserParams{Url: url, UserID: user.ID})

	if err != nil {
		return err
	}
	return nil

}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	result := func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}
		err = handler(s, cmd, user)
		if err != nil {
			return err
		}
		return nil
	}
	return result
}
func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	err = s.db.MarkFeedFetched(context.Background(), feed.FeedID)
	if err != nil {
		return nil
	}
	posts, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return nil
	}
	for i := 0; i < len(posts.Channel.Item); i++ {
		published_at, err := time.Parse("2006-01-02 03:04:05", posts.Channel.Item[i].PubDate)
		if err != nil {
			fmt.Println(err)
			return err
		}
		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Title: posts.Channel.Item[i].Title, Url: posts.Channel.Item[i].Link, Description: posts.Channel.Item[i].Description, PublishedAt: published_at, FeedID: feed.FeedID})
		if err != nil {
			fmt.Println(err)
		}

	}
	return nil

}
func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int
	var err error
	if len(cmd.Args) == 0 {
		limit = 2
	} else {
		limit, err = strconv.Atoi(cmd.Args[0])
		if err != nil {
			limit = 2
		}
	}

	posts, err := s.db.GetPostesForUser(context.Background(), database.GetPostesForUserParams{ReferenceID: user.ID, Limit: int32(limit)})
	if err != nil {
		return err
	}
	for i := 0; i < limit; i++ {
		fmt.Println(posts[i].Title + "\n" + posts[i].Description + "\n" + posts[i].Url + "\n" + posts[i].PublishedAt.String() )
	}
	return nil
}
