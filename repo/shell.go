package repo

import (
	"os"
	"os/exec"
)

func run_shell_in_dir(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func run_shell(name string, args ...string) error {
	return run_shell_in_dir("", name, args...)
}
