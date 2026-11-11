package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ayad/gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func parsePublishedAt(date string) sql.NullTime {
	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC3339,
	}

	for _, format := range formats {
		t, err := time.Parse(format, date)
		if err == nil {
			return sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}
	}

	return sql.NullTime{Valid: false}
}
func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get next feed: %w", err)
	}

	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("couldn't mark feed as fetched: %w", err)
	}

	fmt.Printf("Fetching feed: %s\n", feed.Name)

	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %w", err)
	}

	for _, item := range rssFeed.Channel.Item {
		now := time.Now()

		_, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   now,
			UpdatedAt:   now,
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: parsePublishedAt(item.PubDate),
			FeedID:      feed.ID,
		})

		if err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) && pqErr.Code == "23505" {
				// This post already exists. Skip it.
				continue
			}

			fmt.Printf("couldn't save post %q: %v\n", item.Title, err)
		}
	}

	return nil
}
func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("agg command requires a time_between_reqs argument")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't parse duration: %w", err)
	}

	if timeBetweenRequests <= 0 {
		return fmt.Errorf("duration must be greater than zero")
	}

	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			fmt.Printf("couldn't scrape feeds: %v\n", err)
		}
	}
}
