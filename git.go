package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/google/go-github/v58/github"
)

const (
	userURL = `https://api.github.com/user`
)

type GitConfig struct {
	InsecureSkipVerify bool   `json:"insecureSkipVerify"`
	Token              string `json:"token"`

	Log *slog.Logger
}

type Git struct {
	client *http.Client
	ctx    context.Context
	log    *slog.Logger
	token  string
}

func NewGit(ctx context.Context, cfg GitConfig) *Git {
	return &Git{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: cfg.InsecureSkipVerify,
				},
			},
		},
		ctx:   ctx,
		log:   cfg.Log,
		token: cfg.Token,
	}
}

func (git *Git) User() (*github.User, error) {
	request, err := git.newRequest(http.MethodGet, userURL, nil)
	if err != nil {
		return nil, err
	}

	body, err := git.do(request)
	if err != nil {
		return nil, err
	}

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

	body, err := git.do(request)
	if err != nil {
		return nil, err
	}

	repos := []github.Repository{}

	if err := json.Unmarshal(body, &repos); err != nil {
		return nil, err
	}

	return repos, nil
}

func (git *Git) do(request *http.Request) ([]byte, error) {
	response, err := git.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode/100 != 2 {
		return nil, fmt.Errorf("%d: %s", response.StatusCode, response.Status)
	}

	return io.ReadAll(response.Body)
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
