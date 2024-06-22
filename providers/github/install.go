package github

import (
	"fmt"
	"runtime"
	"slices"
	"strings"
)

type GithubInstall struct {
	Name string 
	Hints *[]string
	Files *[]string
}

type GithubInstallName struct {
	Owner string 
	Repo string 
	Version string 
}

func Install(instructions *GithubInstall) error {

	name := parseName(instructions.Name)

	ctx := &GithubContext{}

	if version == "latest" {
		info, err := ctx.GetLastestReleaseVersion(owner, repo)
		if err != nil {
			return err
		}

		var assets = info.ReleaseInfo.Assets
		set := []gh.Asset
		for _, asset := range assets {
			name := *asset.Name
			name = strings.ToLower(name)
			fragments := strings.Split(name, "-")
			
			if isOs(fragments) {
				if isArch(fragments) {
					set = append(set, asset)
					continue
				}
			}
		}

		l := len(set)
		if (l == 0) {
			return fmt.Errorf("No matches found for blah blah bal")
		}

		var target gh.Asset 
		
		if (l > 1 && instructions.Hints != nil && len(*instructions.Hints) > 0) {
			for _, asset := range set {
				name := *asset.Name
				name = strings.ToLower(name)
				fragments := strings.Split(name, "-")
				isMatch := slices.IndexFunc(fragments, func(sec) (bool) {
					slices.Index(instructions.Hints, sec) > -1;
				}) > -1;

				if (isMatch) {
					target = asset
					break;
				}
			}
		} else {
			target = assets[0]
		}
	}

}

func parseName(name string) (GithubInstallName) {
	version := "latest"
	if (strings.Contains(name, "@")) {
		parts := strings.Split(name, "@")
		name = parts[0]
		version = parts[1]
	}

	parts := strings.Split(name, "/")
	owner := parts[0]
	repo := parts[1]

	result := GithubInstallName {
		Owner: owner,
		Repo: repo,
		Version: version,
	}

	return result
}

func isHint

func isArch(fragments []string) bool {
	// "386", "amd64", "amd64p32", "arm", "arm64", "arm64be", "armbe", "loong64", "mips", "mips64", "mips64le", "mips64p32", "mips64p32le", "mipsle", "ppc", "ppc64", "ppc64le", "riscv", "riscv64", "s390", "s390x", "sparc", "sparc64", "wasm"
	switch runtime.GOARCH {
	case "arm":
		return slices.IndexFunc(fragments, func(sec) (bool) {
			switch(sec) {
			case "arm", "aarch":
				return true
			default:
				return false
			}
		}) > -1
	case "arm64":
		return slices.IndexFunc(fragments, func(sec) (bool) {
			switch(sec) {
			case "arm64", "aarch64":
				return true
			default:
				return false
			}
		}) > -1
	case "amd64":
		return slices.IndexFunc(fragments, func(sec) (bool) {
			switch(sec) {
			case "amd64", "86_64", "x86_64", "x64":
				return true
			default:
				return false
			}
		}) > -1
	case "386":
		return slices.IndexFunc(fragments, func(sec) (bool) {
			switch(sec) {
			case "386", "x86", "x32":
				return true 
			default:
				return false
			}
		}) > -1
	}

	return false;
}

func isOs(fragments []string) bool {
	switch runtime.GOOS {
	case "windows":
		return slices.IndexFunc(fragments, func(sec string) (bool) {
			switch sec {
			case "windows", "win32", "win":
				return true
			default:
				return false
			}
		}) > -1

	case "linux":
		return slices.Index(fragments, "linux") > -1
	case "darwin":
		return slices.IndexFunc(fragments, func(sec string) (bool) {
			switch sec {
			case "darwin", "macos", "mac", "apple", "macosx"
				return true 
			default:
				return false
			}
		}) > -1
		default:
		return false
	}
}
