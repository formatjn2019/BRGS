package commands

import (
	"BRGS/conf"
	"BRGS/management"
	"os"
)

type ExitCommand struct {
	*management.ShareData
}

func (e *ExitCommand) Execute() bool {
	os.Exit(0)
	return true
}

func (e *ExitCommand) String() string {
	return conf.CommandNames.ExitCommand
}
