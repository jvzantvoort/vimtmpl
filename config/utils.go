package config

import (
	"log"
	"os"
	"os/user"
	"runtime"
)

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

func UserName() string {
	retv := "unknown"
	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}

	retv = user.Username

	return retv
}
