package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

type GithubRepo struct {
	Name            string
	FullName        string `json:"full_name"`
	Description     string
	HtmlUrl         string `json:"html_url"`
	StargazersCount int    `json:"stargazers_count"`
	Topics          []string
}

var repoRegex = regexp.MustCompile(`^https?://github\.com/([A-Za-z0-9_-]+)/([A-Za-z0-9_-]+)$`)

func FetchGithubRepository(repoUrl string) (*GithubRepo, error) {
	matches := repoRegex.FindStringSubmatch(repoUrl)

	if len(matches) != 3 {
		return nil, fmt.Errorf("invalid repo url: %s", repoUrl)
	}

	owner := matches[1]
	name := matches[2]

	res, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, name))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get repo metadata: %s", res.Status)
	}

	var repo GithubRepo
	if err := json.NewDecoder(res.Body).Decode(&repo); err != nil {
		return nil, err
	}

	return &repo, nil
}

type GithubGist struct {
	Description string `json:"description"`
	ID          string `json:"id"`
	Owner       struct {
		Login string `json:"login"`
	} `json:"owner"`
}

var gistRegex = regexp.MustCompile(`^https?://gist\.github\.com/[A-Za-z0-9_-]+/([A-Za-z0-9_-]+)$`)

func FetchGithubGist(rawGistUrl string) (*GithubGist, error) {
	matches := gistRegex.FindStringSubmatch(rawGistUrl)
	if len(matches) != 2 {
		return nil, fmt.Errorf("invalid gist url: %s", rawGistUrl)
	}

	gistId := matches[1]
	if gistId == "" {
		return nil, fmt.Errorf("invalid gist url: %s", rawGistUrl)
	}

	res, err := http.Get(fmt.Sprintf("https://api.github.com/gists/%s", gistId))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get gist metadata: %s", res.Status)
	}

	var gist GithubGist
	if err := json.NewDecoder(res.Body).Decode(&gist); err != nil {
		return nil, err
	}

	return &gist, nil
}

type GitCommit struct {
	Sha string `json:"sha"`
}

func GetLastGitCommit(repoUrl string) (*GitCommit, error) {
	matches := repoRegex.FindStringSubmatch(repoUrl)

	if len(matches) != 3 {
		return nil, fmt.Errorf("invalid repo url: %s", repoUrl)
	}

	owner := matches[1]
	name := matches[2]

	apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", owner, name)

	resp, err := http.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get latest release: %s", resp.Status)
	}

	var commits []*GitCommit
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return nil, err
	}

	if len(commits) == 0 {
		return nil, fmt.Errorf("no commits found")
	}

	return commits[0], nil
}

type GistCommit struct {
	Version string `json:"version"`
}

func GetLastGistCommit(rawGistUrl string) (*GistCommit, error) {
	matches := gistRegex.FindStringSubmatch(rawGistUrl)
	if len(matches) != 2 {
		return nil, fmt.Errorf("invalid gist url: %s", rawGistUrl)
	}

	gistId := matches[1]
	if gistId == "" {
		return nil, fmt.Errorf("invalid gist url: %s", rawGistUrl)
	}

	res, err := http.Get(fmt.Sprintf("https://api.github.com/gists/%s/commits", gistId))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get gist metadata: %s", res.Status)
	}

	var commits []*GistCommit
	if err := json.NewDecoder(res.Body).Decode(&commits); err != nil {
		return nil, err
	}

	if len(commits) == 0 {
		return nil, fmt.Errorf("no commits found")
	}

	return commits[0], nil
}
