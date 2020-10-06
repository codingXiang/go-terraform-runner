package command

import (
	"github.com/codingXiang/go-workflow"
	"github.com/codingXiang/go-terraform-runner    /util"
)

const (
	COMMAND      = "terraform"
	INIT         = "init"
	UPGRADE_0_12 = "0.12upgrade"
	UPGRADE_0_13 = "0.13upgrade"
	PLAN         = "plan"
	APPLY        = "apply"
	DESTROY      = "destroy"
)

type Command interface {
	Run(c workflow.Context) error
	GetMeta() *Meta
	SetMeta(meta *Meta)
}

func runCommand(dir string, command string, commands ...string) error {
	return util.RunCommand(COMMAND, dir, command, commands...)
}

type BaseCommand struct {
	*Meta
}

func (c *BaseCommand) init(meta *Meta) {
	if meta == nil {
		c.SetMeta(NewMeta())
	} else {
		c.SetMeta(meta)
	}
}

func (c *BaseCommand) GetMeta() *Meta {
	return c.Meta
}
func (c *BaseCommand) SetMeta(meta *Meta) {
	c.Meta = meta
}
