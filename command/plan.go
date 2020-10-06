package command

import (
	"errors"
	"github.com/codingXiang/go-workflow"
)

type PlanCommand struct {
	BaseCommand
}

func NewPlanCommand(meta *Meta) Command {
	cmd := new(PlanCommand)
	cmd.init(meta)
	return cmd
}

func (c *PlanCommand) init(meta *Meta) {
	c.BaseCommand.init(meta)
	c.Label = "預覽結果"
	c.Meta.Commands = append(c.Meta.Commands, "-input=false", "-out=plan")
}

func (c *PlanCommand) Run(w workflow.Context) error {
	err := runCommand(c.Path, PLAN, c.Commands...)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}
