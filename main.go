package main

import (
	"context"
	"os"
	"strings"

	"github.com/google/go-github/v50/github"
	"github.com/samber/lo"
	"golang.org/x/oauth2"
	log "unknwon.dev/clog/v2"
)

const unfollowedEventFileName = "unfollowed.txt"

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
	var currentFollowerIDs []string
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
			currentFollowerIDs = append(currentFollowerIDs, follower.GetLogin())
		}
		currentPage++
	}

	log.Info("Total followers: %d", len(currentFollowerIDs))

	// Load the gist follower database.
	gistID := os.Getenv("GIST_ID")
	gist, _, err := client.Gists.Get(ctx, gistID)
	if err != nil {
		log.Fatal("Failed to get gist: %v", err)
	}
	// Get the first file as the gist follower database.
	var gistFollowerFile *github.GistFile
	for _, file := range gist.Files {
		gistFollowerFile = &file
		break
	}
	fileName := gistFollowerFile.GetFilename()
	content := gistFollowerFile.GetContent()

	previousFollowerIDs := strings.Split(content, "\n")
	// Diff the two lists.
	unfollowed, newFollowed := lo.Difference(previousFollowerIDs, currentFollowerIDs)
	if len(newFollowed) > 0 {
		log.Trace("Got %d new followers! 🎉", len(newFollowed))
	}
	if len(unfollowed) > 0 {
		log.Warn("Got %d unfollowers! 😢", len(unfollowed))
		// Save it!
		unfollowedEventFile := gist.Files[unfollowedEventFileName]
		unfollowedListContent := unfollowedEventFile.GetContent()
		unfollowedListContent = strings.Join(unfollowed, "\n") + "\n" + unfollowedListContent
		unfollowedEventFile.Content = &unfollowedListContent
		gist.Files[unfollowedEventFileName] = unfollowedEventFile
	}

	// Save the followers to gist database.
	newContent := strings.Join(currentFollowerIDs, "\n")
	gistFollowerFile.Content = &newContent
	gist.Files[github.GistFilename(fileName)] = *gistFollowerFile
	_, _, err = client.Gists.Edit(ctx, gistID, gist)
	if err != nil {
		log.Fatal("Failed to update gist: %v", err)
	}
}
