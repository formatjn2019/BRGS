package menu

import (
	"BRGS/management"
	"BRGS/management/commands"
	"BRGS/pkg/tools"
)

// 菜单树
func MenuTree() {
	//共享数据

	sd := management.ShareData{ServerChan: make(chan struct{})}
	backupCommand := &commands.BackupCommand{ShareData: &sd}
	compressedArchive := &commands.CompressedArchiveCommand{ShareData: &sd}
	exit := &commands.ExitCommand{ShareData: &sd}
	generateConfigDefaultCommand := &commands.GenerateConfigDefaultCommand{ShareData: &sd}
	readConfigCommand := &commands.ReadConfigCommand{ShareData: &sd}
	resoreBackup := &commands.ResoreBackup{ShareData: &sd}
	resoreFileFromArchive := &commands.ResoreFileFromArchive{ShareData: &sd}
	resetBackup := &commands.ResetBackup{ShareData: &sd}
	stopBackupCommond := &commands.StopBackupCommond{ShareData: &sd}
	manualAndAutoBackupCommand := &commands.ManualAndAutoBackupCommand{ShareData: &sd}
	autoBackupCommand := &commands.AutoBackupCommand{ShareData: &sd}
	manualBackupCommand := &commands.ManualBackupCommand{ShareData: &sd}
	startServerCommand := &commands.StartServerCommand{ShareData: &sd}
	stopServerCommand := &commands.StopServerCommand{ShareData: &sd}

	menuConfig := []string{readConfigCommand.String(), generateConfigDefaultCommand.String(), exit.String()}
	menuControl := []string{backupCommand.String(), compressedArchive.String(), resoreBackup.String(), resoreFileFromArchive.String(), resetBackup.String(), stopBackupCommond.String()}
	menuType := []string{autoBackupCommand.String(), manualBackupCommand.String(), manualAndAutoBackupCommand.String(), startServerCommand.String(), stopServerCommand.String()}
	menuAutoBackup := []string{stopBackupCommond.String()}

	closeServer := func() {
		println("close server")
	}
	type Step struct {
		prefix   *[]string
		next     *[]string
		cmd      management.Command
		concelDo func()
	}
	cmdDic := map[string]Step{
		backupCommand.String():                {prefix: &menuControl, next: &menuControl, cmd: backupCommand},
		compressedArchive.String():            {prefix: &menuControl, next: &menuControl, cmd: compressedArchive},
		exit.String():                         {prefix: &menuConfig, next: nil, cmd: exit},
		generateConfigDefaultCommand.String(): {prefix: &menuConfig, next: &menuConfig, cmd: generateConfigDefaultCommand},
		readConfigCommand.String():            {prefix: &menuConfig, next: &menuType, cmd: readConfigCommand},
		resoreBackup.String():                 {prefix: &menuControl, next: &menuControl, cmd: resoreBackup},
		resoreFileFromArchive.String():        {prefix: &menuControl, next: &menuControl, cmd: resoreFileFromArchive},
		resetBackup.String():                  {prefix: &menuConfig, next: &menuControl, cmd: resetBackup},
		stopBackupCommond.String():            {prefix: &menuControl, next: &menuConfig, cmd: stopBackupCommond, concelDo: closeServer},
		manualAndAutoBackupCommand.String():   {prefix: &menuConfig, next: &menuControl, cmd: manualAndAutoBackupCommand, concelDo: closeServer},
		autoBackupCommand.String():            {prefix: &menuConfig, next: &menuAutoBackup, cmd: autoBackupCommand, concelDo: closeServer},
		manualBackupCommand.String():          {prefix: &menuConfig, next: &menuControl, cmd: manualBackupCommand, concelDo: closeServer},
		startServerCommand.String():           {prefix: &menuConfig, next: &menuControl, cmd: startServerCommand, concelDo: closeServer},
		stopServerCommand.String():            {prefix: &menuConfig, next: &menuControl, cmd: stopBackupCommond, concelDo: closeServer},
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
				if cmd.concelDo != nil {
					cmd.concelDo()
				}
			}
		}
	}
}