package menu

import (
	"BRGS/management"
	"BRGS/management/commands"
	"BRGS/pkg/tools"
)

type step struct {
	prefix    *[]string
	next      *[]string
	cmd       management.Command
	consoleDo func()
}

func ConfigMenu() {
	exit := &commands.ExitCommand{}
	generateConfigDefaultCommand := &commands.GenerateConfigDefaultCommand{}
	readConfigCommand := &commands.ReadConfigCommand{}
	menuConfig := []string{readConfigCommand.String(), generateConfigDefaultCommand.String(), exit.String()}
	cmdDic := map[string]step{
		exit.String():                         {prefix: &menuConfig, next: nil, cmd: exit},
		generateConfigDefaultCommand.String(): {prefix: &menuConfig, next: &menuConfig, cmd: generateConfigDefaultCommand},
		readConfigCommand.String():            {prefix: &menuConfig, next: &menuConfig, cmd: readConfigCommand},
	}
	menu := &menuConfig
	for {
		if index := tools.CommandMenu(false, *menu...); index == -1 {
			menu = &menuConfig
		} else {
			cmd := cmdDic[(*menu)[index]]
			// 根据命令执行成败决定下级菜单
			if cmd.cmd.Execute() {
				menu = cmd.next
			} else {
				menu = cmd.prefix
				if cmd.consoleDo != nil {
					cmd.consoleDo()
				}
			}
		}
	}
}

// ControlMenu 菜单树
func ControlMenu() {
	backupCommand := &commands.BackupCommand{}
	restoreBackup := &commands.RestoreBackup{}
	compressedArchive := &commands.CompressedArchiveCommand{}
	hardlinkArchive := &commands.BackupFilesWithHardLinkCommand{}
	restoreFileFromArchive := &commands.RestoreFileFromArchive{}
	resetBackup := &commands.ResetBackup{}
	exit := &commands.ExitCommand{}

	menuControl := []string{backupCommand.String(), compressedArchive.String(), restoreBackup.String(), restoreFileFromArchive.String(), resetBackup.String(), hardlinkArchive.String()}

	cmdDic := map[string]step{
		backupCommand.String():          {prefix: &menuControl, next: &menuControl, cmd: backupCommand},
		compressedArchive.String():      {prefix: &menuControl, next: &menuControl, cmd: compressedArchive},
		exit.String():                   {prefix: &menuControl, next: nil, cmd: exit},
		hardlinkArchive.String():        {prefix: &menuControl, next: &menuControl, cmd: hardlinkArchive},
		restoreBackup.String():          {prefix: &menuControl, next: &menuControl, cmd: restoreBackup},
		restoreFileFromArchive.String(): {prefix: &menuControl, next: &menuControl, cmd: restoreFileFromArchive},
		resetBackup.String():            {prefix: &menuControl, next: &menuControl, cmd: resetBackup},
	}

	menu := &menuControl
	for {
		if index := tools.CommandMenu(false, *menu...); index == -1 {
			menu = &menuControl
		} else {
			cmd := cmdDic[(*menu)[index]]
			// 根据命令执行成败决定下级菜单
			if cmd.cmd.Execute() {
				menu = cmd.next
			} else {
				menu = cmd.prefix
				if cmd.consoleDo != nil {
					cmd.consoleDo()
				}
			}
		}
	}
}
