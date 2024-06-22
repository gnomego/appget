//go:build !windows
// +build !windows

package paths

import (
	"os"
	"os/user"
	"path/filepath"
)

func GetMachineConfigDir() (string, error) {
	return filepath.Join("/etc", "avm"), nil
}

func GetMachineCacheDir() (string, error) {
	return filepath.Join("/var", "cache", "avm"), nil
}

func GetUserConfigDir() (string, error) {
	homeConfigDir := os.Getenv("XDG_CONFIG_HOME")
	if homeConfigDir == "" {
		homeDir := os.Getenv("HOME")
		if homeDir == "" {
			user, err := user.Current()
			if err != nil {
				return "", err
			}

			homeDir = user.HomeDir
		}

		homeConfigDir = filepath.Join(homeDir, ".config")
	}

	return filepath.Join(homeConfigDir, "avm"), nil
}

func GetUserCacheDir() (string, error) {
	homeCacheDir := os.Getenv("XDG_CACHE_HOME")
	if homeCacheDir == "" {
		homeDir := os.Getenv("HOME")
		if homeDir == "" {
			user, err := user.Current()
			if err != nil {
				return "", err
			}

			homeDir = user.HomeDir
		}

		homeCacheDir = filepath.Join(homeDir, ".cache")
	}

	return filepath.Join(homeCacheDir, "avm"), nil
}
