package main

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGit(t *testing.T) {

	cfg := GitConfig{
		log:                log.Default(),
		InsecureSkipVerify: true,
		Token:              os.Getenv("GIT_TOKEN"),
	}
	assert.NotNil(t, cfg.log)
	assert.NotEmpty(t, cfg.Token)

	git := NewGit(context.Background(), cfg)
	assert.NotNil(t, git)

	user, err := git.User()
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.GetEmail())

	repos, err := git.Repos(user)
	assert.NoError(t, err)
	assert.NotNil(t, repos)
	assert.Len(t, repos, 1)

	assert.GreaterOrEqual(t, len(repos), 1)
}
