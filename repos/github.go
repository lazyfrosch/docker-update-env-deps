package repos

import (
	"context"
	"fmt"
	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type GitHub struct {
	URL     string
	Owner   string
	Repo    string
	Client  *github.Client
	APIData *github.Repository
}

const GitHubCom = "github.com"

func NewGitHubClient() *github.Client {
	ctx := context.Background()

	var c *http.Client

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		c = oauth2.NewClient(ctx, ts)
	}

	return github.NewClient(c)
}

func LoadGitHub(rawurl string) (*GitHub, error) {
	// TODO: handle non http URLs
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	if u.Hostname() != GitHubCom {
		return nil, fmt.Errorf("not a URL to " + GitHubCom)
	}

	// normalize path for our use
	path := u.Path
	path = strings.TrimLeft(path, "/")      // Remove any leading slash coming from the URL
	path = strings.TrimSuffix(path, ".git") // Remove trailing .git extension

	// Parse owner and repo from path
	var owner, repo string

	if parts := strings.Split(path, "/"); len(parts) == 2 {
		owner = parts[0]
		repo = parts[1]
	} else {
		return nil, fmt.Errorf("URL path is not in owner/repo syntax")
	}

	// Load API data
	client := NewGitHubClient()

	details, _, err := client.Repositories.Get(context.Background(), owner, repo)
	if err != nil {
		return nil, fmt.Errorf("could not load repo from GitHub API: %w", err)
	}

	return &GitHub{URL: rawurl, Owner: owner, Repo: repo, Client: client, APIData: details}, nil
}

func (r GitHub) LoadReleases() ([]string, error) {
	data, _, err := r.Client.Repositories.ListReleases(context.Background(), r.Owner, r.Repo, nil)
	if err != nil {
		return nil, fmt.Errorf("could not load releases: %w", err)
	}

	var result []string

	for _, release := range data {
		// TODO: filter beta, rc and such
		result = append(result, *release.TagName)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no releases found")
	}

	return result, nil
}
