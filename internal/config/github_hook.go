package config

import "github.com/go-playground/webhooks/v6/github"

func NewGithubHook() (hook *github.Webhook, err error) {

	hook, err = github.New(github.Options.Secret("NqDw9wTeyp"))
	if err != nil {
		return nil, err
	}

	return hook, nil
}
