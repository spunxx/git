package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/go-github/v58/github"
)

type Git struct {
	client *http.Client
	ctx    context.Context
	token  string
}

func NewGit(ctx context.Context, token string) *Git {
	return &Git{
		client: http.DefaultClient,
		ctx:    ctx,
		token:  token,
	}
}

func (git *Git) User() (*github.User, error) {
	request, err := git.newRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}

	response, err := git.client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode/100 != 2 {
		return nil, fmt.Errorf("%s: %d", response.Status, response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	user := github.User{}

	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (git *Git) Repos(user *github.User) ([]github.Repository, error) {
	request, err := git.newRequest(http.MethodGet, user.GetReposURL(), nil)
	if err != nil {
		return nil, err
	}

	response, err := git.client.Do(request)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	repos := []github.Repository{}

	if err := json.Unmarshal(body, &repos); err != nil {
		return nil, err
	}

	return repos, nil
}

func (git *Git) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Accept", "application/vnd.github+json")
	request.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	request.Header.Set("Authorization", "Bearer "+git.token)

	return request, nil
}
