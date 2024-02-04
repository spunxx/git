package main

import (
	"context"
	"fmt"
	"log"
	"os"
)

func main() {

	var (
		ctx = context.Background()
		log = log.New(os.Stdout, "GIT: ", log.LstdFlags)
	)

	token, err := getToken()
	if err != nil {
		log.Print("get token failed: ", err.Error())
		return
	}

	client := NewGit(ctx, GitConfig{
		log:                log,
		InsecureSkipVerify: true,
		Token:              token,
	})

	user, err := client.User()
	if err != nil {
		log.Print("error getting user: ", err)
		return
	}

	log.Print("email: ", user.GetEmail())

}

func getToken() (string, error) {
	token, ok := os.LookupEnv("GIT_TOKEN")
	if !ok {
		return "", fmt.Errorf("GIT_TOKEN is missing")
	}

	if len(token) == 0 {
		return "", fmt.Errorf("GIT_TOKEN is empty")
	}

	return token, nil
}
