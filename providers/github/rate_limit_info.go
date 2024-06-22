package github

import (
	"strconv"

	gh "github.com/google/go-github/v62/github"
)

type RateLimitInfo struct {
	Limit     int
	Remaining int
	Reset     string
	Resource  string
}

func GetRateLimitInfo(resp *gh.Response) *RateLimitInfo {
	rawLimit := resp.Header.Get("X-RateLimit-Limit")
	rawRemaining := resp.Header.Get("X-RateLimit-Remaining")
	rawReset := resp.Header.Get("X-RateLimit-Reset")
	rawResource := resp.Header.Get("X-RateLimit-Resource")

	rl := &RateLimitInfo{
		Limit:     0,
		Remaining: 0,
		Reset:     rawReset,
		Resource:  rawResource,
	}

	if rawLimit != "" {
		l, err := strconv.Atoi(rawLimit)
		if err == nil {
			rl.Limit = l
		}
	}

	if rawRemaining != "" {
		r, err := strconv.Atoi(rawRemaining)
		if err == nil {
			rl.Remaining = r
		}
	}

	return rl
}
