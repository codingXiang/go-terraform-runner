package command

import (
	"fmt"
	"github.com/codingXiang/go-workflow"
)

type UpgradeCommand struct {
	BaseCommand
	upgradeVersion float64
}

func NewUpgradeCommand(meta *Meta, upgradeVersion float64) Command {
	cmd := new(UpgradeCommand)
	cmd.init(meta)
	cmd.upgradeVersion = upgradeVersion
	return cmd
}

func (c *UpgradeCommand) init(meta *Meta) {
	c.BaseCommand.init(meta)
	c.Label = fmt.Sprintf("升級至版本 %f", c.upgradeVersion)
	c.Meta.Commands = append(c.Meta.Commands, "-yes")
}

func (u UpgradeCommand) Run(w workflow.Context) error {
	switch u.upgradeVersion {
	case 0.12:
		return runCommand(u.Path, UPGRADE_0_12, u.Commands...)
	default:
		return runCommand(u.Path, UPGRADE_0_13, u.Commands...)
	}
}
