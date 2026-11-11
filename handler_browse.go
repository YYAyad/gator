package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ayad/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2

	if len(cmd.args) > 1 {
		return fmt.Errorf("browse accepts at most one limit argument")
	}

	if len(cmd.args) == 1 {
		parsedLimit, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}

		if parsedLimit <= 0 {
			return fmt.Errorf("limit must be greater than zero")
		}

		limit = parsedLimit
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts: %w", err)
	}

	for _, post := range posts {
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("URL: %s\n", post.Url)
		fmt.Printf("Description: %s\n", post.Description)
		fmt.Println()
	}

	return nil
}
