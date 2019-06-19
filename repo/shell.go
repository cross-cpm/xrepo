package repo

import (
	"os"
	"os/exec"
)

func run_shell(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func run_shell_in_dir(dir string, name string, args ...string) error {
	cwd, err := os.Getwd()
	// log.Println("cwd", cwd, err)
	if err != nil {
		return err
	}

	err = os.Chdir(dir)
	if err != nil {
		return err
	}

	err = run_shell(name, args...)
	if err != nil {
		return err
	}

	err = os.Chdir(cwd)
	if err != nil {
		return err
	}

	return nil
}
