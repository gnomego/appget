//go:build windows
// +build windows

package paths

import (
	"os"
	"os/user"
	"path/filepath"
)

func GetMachineConfigDir() (string, error) {
	dir := os.Getenv("ALLUSERSPROFILE")
	if dir == "" {
		dir = "C:\\ProgramData"
	}

	return filepath.Join(dir, "avm", "etc"), nil
}

func GetMachineCacheDir() (string, error) {
	dir := os.Getenv("ALLUSERSPROFILE")
	if dir == "" {
		dir = "C:\\ProgramData"
	}

	return filepath.Join(dir, "avm", "cache"), nil
}

func GetUserConfigDir() (string, error) {
	dir := os.Getenv("APPDATA")
	if dir == "" {
		userProfile := os.Getenv("USERPROFILE")
		if userProfile == "" {
			usr, err := user.Current()
			if err != nil {
				return "", err
			}
			userProfile = usr.HomeDir
		}

		return filepath.Join(userProfile, "AppData", "Roaming", "avm", "etc"), nil
	}

	return filepath.Join(dir, "avm", "etc"), nil
}

func GetUserCacheDir() (string, error) {
	dir := os.Getenv("APPDATALOCAL")
	if dir == "" {
		userProfile := os.Getenv("USERPROFILE")
		if userProfile == "" {
			usr, err := user.Current()
			if err != nil {
				return "", err
			}
			userProfile = usr.HomeDir
		}

		return filepath.Join(userProfile, "AppData", "Local", "avm", "cache"), nil
	}

	return filepath.Join(dir, "avm", "cache"), nil
}
