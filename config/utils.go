
// Package config provides utility functions for user and file operations.
package config

import (
	"log"
	"os"
	"os/user"
	"runtime"
)

// UserHomeDir returns the home directory of the current user, handling platform differences.
func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

// UserName returns the username of the current user.
func UserName() string {
	user, err := user.Current()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	retv := user.Username

	return retv
}

// TargetExists returns true if the target file or directory exists.
func TargetExists(targetpath string) bool {
	_, err := os.Stat(targetpath)
	if err != nil {
		return false
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}
