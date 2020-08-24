package client

import (
	"context"
	"errors"
	"net/http"

	"golang.org/x/oauth2"
)

func NewOAuthClient(accessToken string) *http.Client {
	if accessToken == "" {
		panic(errors.New("GitHub access token was not provided"))
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	return oauth2.NewClient(ctx, ts)
}
