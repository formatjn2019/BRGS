package test

import (
	"BRGS/management"
	"BRGS/pkg/util"
	"testing"
)

func TestWriteDefaultConfig(t *testing.T) {
	generateConfigDefaultCommand := &management.GenerateConfigDefaultCommand{}
	generateConfigDefaultCommand.Execute()
}

func TestReadConfigCommand(t *testing.T) {
	util.DictReverse(management.EXCEL_HEAD_TRANSLATE_DIC)
	readConfigCommand := &management.ReadConfigCommand{}
	readConfigCommand.Execute()
}
