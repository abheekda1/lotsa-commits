package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func main() {
	r, err := git.PlainOpen("repo")
	if err != nil {
		createR, err := git.PlainInit("repo", false)
		if err != nil {
			fmt.Println("An error with opening the local repo and creating a new repo occurred:", err)
			return
		} else {
			r = createR
		}
	}

	cCount := 0
	ref, err := r.Head()
	if err != nil {
		if err.Error() != "reference not found" {
			panic(err)
		}
	} else {
		cIter, err := r.Log(&git.LogOptions{
			From: ref.Hash(),
		})

		if err != nil {
			fmt.Println("An error with getting the log occurred:", err)
			return
		}

		err = cIter.ForEach(func(c *object.Commit) error {
			cCount++

			return nil
		})

		if err != nil {
			fmt.Println("An error with looping through past commits occurred:", err)
			return
		}
	}

	fmt.Println(cCount, "commits up till now")

	w, err := r.Worktree()
	if err != nil {
		fmt.Println("An error with getting the worktree occurred:", err)
	}

	auth := &http.BasicAuth{
		Username: os.Getenv("GH_USERNAME"),
		Password: os.Getenv("GH_PASSWORD"),
	}

	authorName := os.Getenv("AUTHOR_NAME")
	authorEmail := os.Getenv("AUTHOR_EMAIL")

	if authorName == "" {
		authorName = "Lotsa Commits by abheekda1"
	}

	for i := cCount; true; i++ {
		if i%1000000 == 0 && i > 0 {
			fmt.Println(i, "commits")
			err = r.Push(&git.PushOptions{
				RemoteURL: os.Getenv("GH_REMOTE_URL"),
				Auth:      auth,
			})
			if err != nil {
				fmt.Println("An error with pushing occurred:", err)
				return
			}
		}
		_, err := w.Commit(fmt.Sprintf("Commit #%d", i+1), &git.CommitOptions{
			Author: &object.Signature{
				Name:  authorName,
				Email: authorEmail,
				When:  time.Now(),
			},
			AllowEmptyCommits: true,
		})

		if err != nil {
			fmt.Println("An error with committing occurred:", err)
			return
		}
	}
}
