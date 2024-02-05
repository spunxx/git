package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
)

type Config struct {
	Debug     bool      `yaml:"debug"`
	GitConfig GitConfig `yaml:"gitConfig"`
}

func main() {

	var (
		config  Config
		LogOpts slog.HandlerOptions
	)

	flag.BoolVar(&config.Debug, "debug", false, "enable/disable debug logging")
	flag.BoolVar(&config.GitConfig.InsecureSkipVerify, "skipVerify", true, "skip certificate verification")
	flag.StringVar(&config.GitConfig.Token, "token", "", "github authentication token")
	flag.Parse()

	if config.Debug {
		LogOpts.Level = slog.LevelDebug
	} else {
		LogOpts.Level = slog.LevelInfo
	}

	var (
		ctx = context.Background()
		log = slog.New(slog.NewTextHandler(os.Stdout, &LogOpts))
	)

	log.Debug("debug logging enabled")

	token, err := getToken()
	if err != nil {
		log.Error("get token failed", "error", err.Error())
		return
	}

	client := NewGit(ctx, GitConfig{
		Log:                log,
		InsecureSkipVerify: true,
		Token:              token,
	})

	user, err := client.User()
	if err != nil {
		log.Error("get user failed", "error", err.Error())
		return
	}

	log.Info("user", "email", user.GetEmail())
	log.Debug("user", "user", user)

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
