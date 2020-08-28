package cmd

import (
	"log"
	"os"
)

var workDir string

func setInitialWorkDir() {
	dir,err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	workDir = dir
}

// GetWorkDir returns the workDir
func GetWorkDir() string {
	if workDir == "" {
		setInitialWorkDir()
	}
	return workDir
}

func setWorkDir(dir string) {
	workDir = dir
	os.Chdir(dir) //TODO: #14 Handle wrong directories on cd
}