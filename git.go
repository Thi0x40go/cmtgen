package main

import (
	"os"
	"os/exec"
)

func getGitDiff() (string, error) {
	out, err := exec.Command("git", "diff", "--staged").Output()
	return string(out), err
}

func executeCommit(msg string) error {
	cmd := exec.Command("git", "commit", "-m", msg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
