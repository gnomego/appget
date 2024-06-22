package github

import (
	"context"
	ctx "context"
	"net/http"
)

type GithubContext struct {
	Context *ctx.Context
	Client  *http.Client
}

func (ctx *GithubContext) GetContext() context.Context {
	if ctx == nil {
		return context.Background()
	}

	if ctx.Context == nil {
		*ctx.Context = context.Background()
	}

	return *ctx.Context
}
