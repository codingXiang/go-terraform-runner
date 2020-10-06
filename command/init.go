package command

import "github.com/codingXiang/go-workflow"

type InitCommand struct {
	BaseCommand
}

func NewInitCommand(meta *Meta) Command {
	cmd := new(InitCommand)
	cmd.init(meta)
	return cmd
}

func (c *InitCommand) init(meta *Meta) {
	c.BaseCommand.init(meta)
	c.Label = "初始化"
}

func (c *InitCommand) Run(w workflow.Context) error {
	return runCommand(c.Path, INIT, c.Commands...)
}
