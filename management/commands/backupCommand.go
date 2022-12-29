package commands

import (
	"BRGS/conf"
	"BRGS/management"
	"BRGS/models"
	"fmt"
)

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
	return conf.CommandNames.BackupCommand
}

// 重置备份命令
type ResetBackup struct {
	*management.ShareData
}

func (r *ResetBackup) Execute() bool {
	r.Tree = *models.CreateFsTreeRoot(r.BackupArchive.WatchDir, r.BackupArchive.TempDir)
	fmt.Println("重置执行")
	return true
}

func (r *ResetBackup) String() string {
	return conf.CommandNames.ResetBackup
}

// 还原命令
type ResoreBackup struct {
	*management.ShareData
}

func (r *ResoreBackup) Execute() bool {
	if r.Tree.RecoverFiles() {
		fmt.Println("同步执行成功")
		return true
	} else {
		fmt.Println("同步执行失败")
		return false
	}
}

func (r *ResoreBackup) String() string {
	return conf.CommandNames.ResoreBackup
}

// 停止备份命令
type StopBackupCommond struct {
	*management.ShareData
}

func (b *StopBackupCommond) Execute() bool {
	fmt.Println("停止同步执行")
	return true
}

func (b *StopBackupCommond) String() string {
	return conf.CommandNames.StopBackupCommond
}
