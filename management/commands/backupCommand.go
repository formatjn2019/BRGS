package commands

import (
	"BRGS/conf"
	"BRGS/management"
	"fmt"
)

// BackupCommand 备份命令
type BackupCommand struct {
}

func (b *BackupCommand) Execute() bool {
	fmt.Println("同步中转执行")
	return management.MonServer.Backup()
}

func (b *BackupCommand) String() string {
	return conf.CommandNames.BackupCommand
}

// ResetBackup 重置备份命令
type ResetBackup struct {
}

func (r *ResetBackup) Execute() bool {
	fmt.Println("重置执行")
	return management.MonServer.Scanning()
}

func (r *ResetBackup) String() string {
	return conf.CommandNames.ResetBackup
}

// RestoreBackup 还原命令
type RestoreBackup struct {
}

func (r *RestoreBackup) Execute() bool {
	if management.MonServer.Recover() {
		fmt.Println("同步执行成功")
		return true
	} else {
		fmt.Println("同步执行失败")
		return false
	}
}

func (r *RestoreBackup) String() string {
	return conf.CommandNames.RestoreBackup
}

// StopBackupCommand 停止备份命令
type StopBackupCommand struct {
}

func (b *StopBackupCommand) Execute() bool {
	fmt.Println("停止同步执行")
	return true
}

func (b *StopBackupCommand) String() string {
	return conf.CommandNames.StopBackupCommand
}
