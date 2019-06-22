package main

import (
	"os"
	"os/exec"
)

func shell_run(name string, args ...string) error {
	return shell_run_in_dir("", name, args...)
}

func shell_run_in_dir(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func shell_get_in_dir(dir string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return "", nil
	}

	return string(output), nil
}
