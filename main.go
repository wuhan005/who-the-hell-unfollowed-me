package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
	log "unknwon.dev/clog/v2"
)

const fileName = "github-followers.txt"

func main() {
	defer log.Stop()
	if err := log.NewConsole(); err != nil {
		panic("init log: " + err.Error())
	}

	githubToken := os.Getenv("GH_TOKEN")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Get current user.
	currentUser, _, err := client.Users.Get(ctx, "")
	if err != nil {
		log.Fatal("Failed to get current user: %v", err)
	}
	log.Info("Current user: %s", currentUser.GetLogin())

	// Get user's follower list.
	var currentFollowers []string
	currentPage := 1
	for {
		log.Trace("Fetching followers, page %d", currentPage)

		followers, _, err := client.Users.ListFollowers(ctx, currentUser.GetLogin(), &github.ListOptions{
			Page:    currentPage,
			PerPage: 100,
		})
		if err != nil {
			log.Fatal("Failed to get user's followers, page: %d: %v", currentPage, err)
		}

		if len(followers) == 0 {
			break
		}
		for _, follower := range followers {
			followerID := follower.GetID()
			followerName := follower.GetLogin()

			currentFollowers = append(currentFollowers, fmt.Sprintf("%s (%d)", followerName, followerID))
		}
		currentPage++
	}

	if len(currentFollowers) == 0 {
		// Something unexpected happened.
		log.Fatal("Failed to get user's followers, empty list.")
	}

	log.Info("Total followers: %d", len(currentFollowers))

	gistID := os.Getenv("GIST_ID")
	gist, _, err := client.Gists.Get(ctx, gistID)
	if err != nil {
		log.Fatal("Failed to get gist: %v", err)
	}
	if len(gist.Files) == 0 {
		log.Fatal("Gist has no files.")
	}

	files := gist.GetFiles()
	content := strings.Join(currentFollowers, "\n")

	// Check if content has changed before updating
	file := files[fileName]
	oldContent := file.GetContent()
	if oldContent == content {
		log.Info("Gist content unchanged, skipping update.")
		return
	}

	*file.Content = content
	files[fileName] = file
	gist.Files = files
	_, _, err = client.Gists.Edit(ctx, gistID, gist)
	if err != nil {
		log.Fatal("Failed to update gist: %v", err)
	}
	log.Info("Gist updated successfully.")
}
