package main

import (
	"flag"
	"log"
	"os"
	"os/exec"

	"path"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Clones all public gists for user
// Requires an auth token be provided
// If repo exists, pull instead
// Doesn't clean up deleted gists

var allGists []github.Gist

var (
	flagUser = flag.String("user", "", "username for gists")
    flagDir = flag.String("dir", "", "output directory")
    flagToken = flag.String("token", "", "github api auth token")
)

func main() {
    flag.Parse()
	if (*flagUser == "" || *flagDir == "" || *flagToken == "") {
        log.Fatal("all arguments must be provided")
    } 
    ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *flagToken})
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	c := github.NewClient(tc)

	listOpts := &github.GistListOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	for {
		gists, resp, err := c.Gists.List(*flagUser, listOpts)
		if err != nil {
			log.Fatal(err)
		}
		allGists = append(allGists, gists...)
		if resp.NextPage == 0 {
			log.Println(resp.Rate.String())
			break
		}
		listOpts.ListOptions.Page = resp.NextPage
	}

	for _, gist := range allGists {
		if !*gist.Public {
			continue
		}
		dir := path.Join(*flagDir, *gist.ID)
		cmdArgs := []string{"-C", dir, "pull"}
		if _, err := os.Stat(dir); err != nil {
			if os.IsNotExist(err) {
				cmdArgs = []string{"clone", *gist.GitPullURL, dir}
			}
		}
		log.Println(cmdArgs)
		cmd := exec.Command("git", cmdArgs...)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}
