package management

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// 菜单树
func MenuTree() {
	//共享数据
	sd := shareData{}
	backupCommand := &BackupCommand{shareData: &sd}
	compressedArchive := &CompressedArchive{shareData: &sd}
	exit := &Exit{shareData: &sd}
	generateConfigDefaultCommand := &GenerateConfigDefaultCommand{shareData: &sd}
	readConfigCommand := &ReadConfigCommand{shareData: &sd}
	resoreBackup := &ResoreBackup{shareData: &sd}
	resoreFileFromArchive := &ResoreFileFromArchive{shareData: &sd}
	resetBackup := &ResetBackup{shareData: &sd}
	stopBackupCommond := &StopBackupCommond{shareData: &sd}
	manualAndAutoBackupCommand := &ManualAndAutoBackupCommand{shareData: &sd}
	autoBackupCommand := &AutoBackupCommand{shareData: &sd}
	manualBackupCommand := &ManualBackupCommand{shareData: &sd}

	menuConfig := []string{readConfigCommand.String(), generateConfigDefaultCommand.String(), exit.String()}
	menuControl := []string{backupCommand.String(), compressedArchive.String(), resoreBackup.String(), resoreFileFromArchive.String(), resetBackup.String(), stopBackupCommond.String()}
	menuType := []string{autoBackupCommand.String(), manualBackupCommand.String(), manualAndAutoBackupCommand.String()}
	menuAutoBackup := []string{stopBackupCommond.String()}
	type Step struct {
		prefix *[]string
		next   *[]string
		cmd    Command
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
		stopBackupCommond.String():            {prefix: &menuControl, next: &menuConfig, cmd: stopBackupCommond},
		manualAndAutoBackupCommand.String():   {prefix: &menuConfig, next: &menuControl, cmd: manualAndAutoBackupCommand},
		autoBackupCommand.String():            {prefix: &menuConfig, next: &menuAutoBackup, cmd: autoBackupCommand},
		manualBackupCommand.String():          {prefix: &menuConfig, next: &menuControl, cmd: manualBackupCommand},
	}

	menu := &menuConfig
	for {
		if index := CommandMenu(false, *menu...); index == -1 {
			menu = &menuConfig
		} else {
			cmd := cmdDic[(*menu)[index]]
			// 根据命令执行成败决定下级菜单
			if cmd.cmd.Execute() {
				menu = cmd.next
			} else {
				menu = cmd.prefix
			}
		}
	}
}

// 命令行菜单
func CommandMenu(canBack bool, menuItems ...string) int {
	if canBack {
		menuItems = append(menuItems, "返回上级")
	}
	menuSelect := func(menuItems ...string) int {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			for index, menuItem := range menuItems {
				fmt.Printf("%5d.%80s\n", index+1, menuItem)
			}
			scanner.Scan()
			text := scanner.Text()
			if index, err := strconv.Atoi(text); err == nil {
				if index <= len(menuItems) && index > 0 {
					return index - 1
				} else {
					fmt.Println("输入无效，请重新输入")
				}
			}
		}
	}
	//索引转换
	index := menuSelect(menuItems...)
	if canBack {
		return (index+1)%len(menuItems) - 1
	}
	return index
}
