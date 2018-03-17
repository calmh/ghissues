package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	ExitSuccess = iota
	ExitError
	ExitUsage
)

var tpl = template.Must(template.New("issue.md").ParseFiles("issue.md"))

func main() {
	flag.Parse()

	if flag.NArg() != 3 {
		fmt.Println("ghissues - Create offline copy of GitHub issues list.")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Printf("	%s <owner> <repo> <output directory>\n", os.Args[0])
		fmt.Println()
		fmt.Println("Example:")
		fmt.Printf("	%s syncthing/syncthing /tmp/syncthing-issues\n", os.Args[0])
		fmt.Println()
		fmt.Println("The output directory will be created if it does not exist.")
		os.Exit(ExitUsage)
	}

	owner := flag.Arg(0)
	repo := flag.Arg(1)
	outDir := flag.Arg(2)
	if err := os.MkdirAll(outDir, 0777); err != nil {
		fmt.Println("Output dir:", err)
		os.Exit(ExitError)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	issues, err := loadIssues(client, owner, repo)
	if err != nil {
		fmt.Println("Loading issues:", err)
		os.Exit(ExitError)
	}

	for _, issue := range issues {
		issue, _, err := client.Issues.Get(context.TODO(), owner, repo, issue.GetNumber())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		out, err := os.Create(filepath.Join(outDir, fmt.Sprintf("%d.md", issue.GetNumber())))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = tpl.Execute(out, map[string]interface{}{
			"issue": issue,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = out.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func loadIssues(client *github.Client, owner, repo string) ([]*github.Issue, error) {
	var issues []*github.Issue
	opts := &github.IssueListByRepoOptions{
		State:       "all",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	for {
		is, resp, err := client.Issues.ListByRepo(context.TODO(), owner, repo, opts)
		if err != nil {
			return issues, err
		}
		issues = append(issues, is...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return issues, nil
}
