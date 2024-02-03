package main

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGit(t *testing.T) {

	token := os.Getenv("GIT_TOKEN")
	assert.NotEmpty(t, token)

	git := NewGit(context.Background(), token)
	assert.NotNil(t, git)

	user, err := git.User()
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.GetEmail())

	t.Log("repos url: ", user.GetReposURL())

	repos, err := git.Repos(user)
	assert.NoError(t, err)
	assert.NotNil(t, repos)
	assert.Len(t, repos, 1)

	t.Log("repo name: ", repos[0].GetName())
	t.Log("repo branches: ", repos[0].GetBranchesURL())

}
