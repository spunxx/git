package main

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/v58/github"
)

func main() {

	var (
		token = getToken()
		ctx   = context.Background()
		log   = log.Logger{}
	)

	client := github.NewClient(nil).WithAuthToken(token)

	user, _, err := client.Users.Get(ctx, "rsuther")
	if err != nil {
		log.Print("error getting user: ", err)
	}

	log.Print(user.Email)

}

func getToken() string {
	token, ok := os.LookupEnv("GIT_TOKEN")
	if !ok {
		panic("GIT_TOKEN env var does not exist")
	}

	if len(token) == 0 {
		panic("GIT_TOKEN is empty")
	}

	return token
}
