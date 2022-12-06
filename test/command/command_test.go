package test

import (
	"BRGS/management"
	"BRGS/pkg/tools"
	"testing"
)

func TestWriteDefaultConfig(t *testing.T) {
	generateConfigDefaultCommand := &management.GenerateConfigDefaultCommand{}
	generateConfigDefaultCommand.Execute()
}

func TestReadConfigCommand(t *testing.T) {
	tools.DictReverse(management.EXCEL_HEAD_TRANSLATE_DIC)
	readConfigCommand := &management.ReadConfigCommand{}
	readConfigCommand.Execute()
}
