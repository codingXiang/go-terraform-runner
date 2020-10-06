package command

import (
	"errors"
	"github.com/codingXiang/go-workflow"
)

type ApplyCommand struct {
	BaseCommand
}

func NewApplyCommand(meta *Meta) Command {
	cmd := new(ApplyCommand)
	cmd.init(meta)
	return cmd
}

func (c *ApplyCommand) init(meta *Meta) {
	c.BaseCommand.init(meta)
	c.Label = "執行腳本"
	c.Meta.Commands = append(c.Meta.Commands, "plan")
}

func (c *ApplyCommand) Run(w workflow.Context) error {
	err := runCommand(c.Path, APPLY, c.Commands...)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}
