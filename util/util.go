package util

import (
	"github.com/codingXiang/go-logger/v2"
	"os/exec"
)

func RunCommand(base, dir, command string, commands ...string) error {
	commands = append([]string{command}, commands...)
	cmd := exec.Command(base, commands...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	logger.Log.Info(string(out))
	return err
}
