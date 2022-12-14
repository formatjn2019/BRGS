package commands

import (
	"BRGS/management"
	"fmt"
)

// 停止备份命令
type StopBackupCommond struct {
	*management.ShareData
}

func (b *StopBackupCommond) Execute() bool {
	fmt.Println("停止同步执行")
	return true
}

func (b *StopBackupCommond) String() string {
	return fmt.Sprintf("停止同步文件夹内容")
}

// 备份命令
type BackupCommand struct {
	*management.ShareData
}

func (b *BackupCommand) Execute() bool {
	if b.Tree.BackupFiles() {
		fmt.Println("同步执行成功")
		return true
	} else {
		fmt.Println("同步执行失败")
		return false
	}
}

func (b *BackupCommand) String() string {
	return fmt.Sprintf("同步文件夹内容")
}
