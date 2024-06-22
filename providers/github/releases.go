package github

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/gnomego/avm/paths"
	gh "github.com/google/go-github/v62/github"
)

type SingleRelease struct {
	Owner         string
	Repo          string
	RateLimitInfo *RateLimitInfo
	ReleaseInfo   *ReleaseResult
}

type ReleaseResult struct {
	TagName      string
	HitRateLimit bool
	Assets       []*gh.ReleaseAsset
}

type Releases struct {
	Owner         string
	Repo          string
	RateLimitInfo *RateLimitInfo
	Releases      []*ReleaseResult
}

func (ctx *GithubContext) ListReleases(owner string, repo string, refreshCache *bool) (*Releases, error) {
	var releases []*gh.RepositoryRelease
	readCache := true
	if refreshCache != nil {
		readCache = !(*refreshCache)
	}
	if readCache {
		cacheDir, err := paths.GetUserCacheDir()
		if err != nil {
			return nil, err
		}

		repoCacheDir := filepath.Join(cacheDir, "github", owner, repo)
		if _, err := os.Stat(repoCacheDir); os.IsNotExist(err) {
			os.MkdirAll(repoCacheDir, os.ModePerm)
		}
		cachePath := filepath.Join(repoCacheDir, "releases.json")

		if stat, err := os.Stat(cachePath); err == nil {

			now := time.Now()
			duration := now.Sub(stat.ModTime())
			if duration.Hours() > 24 {
				os.Remove(cachePath)
			} else {
				data, err := os.ReadFile(cachePath)
				if err != nil {
					return nil, err
				}

				releases = make([]*gh.RepositoryRelease, 0)

				json.Unmarshal(data, &releases)
			}
		}
	}

	result := &Releases{
		Owner:    owner,
		Repo:     repo,
		Releases: make([]*ReleaseResult, 0),
	}

	if releases == nil {

		opts := &gh.ListOptions{
			Page:    1,
			PerPage: 100,
		}
		releases, resp, err := gh.NewClient(ctx.Client).Repositories.ListReleases(ctx.GetContext(), owner, repo, opts)
		rl := false
		if _, ok := err.(*gh.RateLimitError); ok {
			rl = true
		} else if err != nil {
			return nil, err
		}

		rateInfo := GetRateLimitInfo(resp)
		result.RateLimitInfo = rateInfo

		if releases != nil {
			cacheDir, err := paths.GetUserCacheDir()
			if err != nil {
				return nil, err
			}

			repoCacheDir := filepath.Join(cacheDir, "github", owner, repo)
			if _, err := os.Stat(repoCacheDir); os.IsNotExist(err) {
				os.MkdirAll(repoCacheDir, os.ModePerm)
			}
			cachePath := filepath.Join(repoCacheDir, "releases.json")

			data, err := json.Marshal(releases)
			if err != nil {
				return nil, err
			}

			err = os.WriteFile(cachePath, data, os.ModePerm)
			if err != nil {
				return nil, err
			}
		}

		set := result.Releases

		for _, release := range releases {
			releaseInfo := &ReleaseResult{
				TagName:      release.GetTagName(),
				HitRateLimit: rl,
				Assets:       release.Assets,
			}

			set = append(set, releaseInfo)
		}

		result.Releases = set
	}

	return result, nil
}

func (ctx *GithubContext) GetReleaseVersion(owner string, repo string, version string) (*SingleRelease, error) {

	release, resp, err := gh.NewClient(ctx.Client).Repositories.GetReleaseByTag(ctx.GetContext(), owner, repo, version)
	rl := false
	if _, ok := err.(*gh.RateLimitError); ok {
		rl = true
	} else if err != nil {
		return nil, err
	}

	rateInfo := GetRateLimitInfo(resp)

	releaseInfo := &ReleaseResult{
		TagName:      release.GetTagName(),
		HitRateLimit: rl,
		Assets:       release.Assets,
	}

	result := &SingleRelease{
		Owner:         owner,
		Repo:          repo,
		RateLimitInfo: rateInfo,
		ReleaseInfo:   releaseInfo,
	}

	return result, nil
}

func (ctx *GithubContext) GetLastestReleaseVersion(owner string, repo string) (*SingleRelease, error) {
	release, resp, err := gh.NewClient(ctx.Client).Repositories.GetLatestRelease(ctx.GetContext(), owner, repo)
	rl := false
	if _, ok := err.(*gh.RateLimitError); ok {
		rl = true
	} else if err != nil {
		return nil, err
	}

	rateInfo := GetRateLimitInfo(resp)

	releaseInfo := &ReleaseResult{
		TagName:      release.GetTagName(),
		HitRateLimit: rl,
		Assets:       release.Assets,
	}

	result := &SingleRelease{
		Owner:         owner,
		Repo:          repo,
		RateLimitInfo: rateInfo,
		ReleaseInfo:   releaseInfo,
	}

	return result, nil
}
