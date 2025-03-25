package config

import (
	"os"
	"path"
	"runtime"
)

const (
	ReadWriteExecutePermission = 0755
)

var (
	ProjectPath string
	EnvPath     string
)

func InitPaths() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	
	SetPaths(home)
	return nil
}

func SetPaths(home string) {
	isWindows := runtime.GOOS == "windows"
	
	if isWindows {
		ProjectPath = path.Join(home, "denv")
	} else {
		ProjectPath = path.Join(home, ".config", "denv")
	}
	
	EnvPath = path.Join(ProjectPath, ".env")
}