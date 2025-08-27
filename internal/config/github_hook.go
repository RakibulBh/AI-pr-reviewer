package config

import "github.com/go-playground/webhooks/v6/github"

func NewGithubHook(token string) (hook *github.Webhook, err error) {

	hook, err = github.New(github.Options.Secret(token))
	if err != nil {
		return nil, err
	}

	return hook, nil
}
