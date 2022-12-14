package test

import (
	"BRGS/management"
	"BRGS/management/commands"
	"BRGS/pkg/tools"
	"testing"
)

func TestWriteDefaultConfig(t *testing.T) {
	generateConfigDefaultCommand := &commands.GenerateConfigDefaultCommand{}
	generateConfigDefaultCommand.Execute()
}

func TestReadConfigCommand(t *testing.T) {
	tools.DictReverse(management.EXCEL_HEAD_TRANSLATE_DIC)
	readConfigCommand := &commands.ReadConfigCommand{}
	readConfigCommand.Execute()
}
