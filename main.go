package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

const (
	ExitSuccess = iota
	ExitError
	ExitUsage
)

type issue struct {
	Number int
	Title  string
	User   struct {
		Login string
	}
	Labels []struct {
		Name string
	}
	State string
	Body  string
}

func main() {
	flag.Parse()

	if flag.NArg() != 2 {
		fmt.Println("ghissues - Create offline copy of GitHub issues list.")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Printf("	%s <repo> <output directory>\n", os.Args[0])
		fmt.Println()
		fmt.Println("Example:")
		fmt.Printf("	%s syncthing/syncthing /tmp/syncthing-issues\n", os.Args[0])
		fmt.Println()
		fmt.Println("The output directory will be created if it does not exist.")
		os.Exit(ExitUsage)
	}

	repo := flag.Arg(0)
	outDir := flag.Arg(1)
	if err := os.MkdirAll(outDir, 0777); err != nil {
		fmt.Println("Output dir:", err)
		os.Exit(ExitError)
	}

	issues, err := loadIssues(repo)
	if err != nil {
		fmt.Println("Loading issues:", err)
		os.Exit(ExitError)
	}

	if err := writeIndex(issues, outDir); err != nil {
		fmt.Println("Write index:", err)
		os.Exit(ExitError)
	}

	for _, issue := range issues {
		if err := writeIssue(issue, outDir); err != nil {
			fmt.Println("Write issue:", err)
			os.Exit(ExitError)
		}
	}
}

func loadIssues(repo string) ([]issue, error) {
	// TODO: Load issue comments. Requires a request per issue, will be
	// subject to rate limits...

	var issues []issue

	link := "https://" + path.Join("api.github.com/repos", repo, "issues")
	for link != "" {
		fmt.Println("Loading", link, "...")

		resp, err := http.Get(link)
		if err != nil {
			return nil, err
		}

		var is []issue
		err = json.NewDecoder(resp.Body).Decode(&is)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		issues = append(issues, is...)

		link = parseRel(resp.Header.Get("Link"), "next")
	}

	return issues, nil
}

func parseRel(link, rel string) string {
	exp := regexp.MustCompile(`<([^>]+)>;\s+rel="` + rel + `"`)
	match := exp.FindStringSubmatch(link)
	if len(match) == 2 {
		return match[1]
	}
	return ""
}

func writeIndex(issues []issue, outDir string) error {
	fd, err := os.Create(filepath.Join(outDir, "index.html"))
	if err != nil {
		return err
	}

	err = indexTpl.Execute(fd, map[string]interface{}{
		"issues": issues,
	})
	if err != nil {
		return err
	}
	return fd.Close()
}

func writeIssue(issue issue, outDir string) error {
	fd, err := os.Create(filepath.Join(outDir, fmt.Sprintf("issue-%d.html", issue.Number)))
	if err != nil {
		return err
	}
	err = issueTpl.Execute(fd, issue)
	if err != nil {
		return err
	}
	return fd.Close()
}
