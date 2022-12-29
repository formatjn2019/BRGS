package commands

import (
	"BRGS/conf"
	"BRGS/management"
	"BRGS/models"
	"fmt"
	"regexp"
)

// 自动备份命令
type AutoBackupCommand struct {
	*management.ShareData
}

func (r *AutoBackupCommand) Execute() bool {
	r.Tree = *models.CreateFsTreeRoot(r.BackupArchive.WatchDir, r.BackupArchive.TempDir)
	return true
}

func (r *AutoBackupCommand) String() string {
	return conf.CommandNames.AutoBackupCommand
}

// 手动备份命令
type ManualBackupCommand struct {
	*management.ShareData
}

func (r *ManualBackupCommand) Execute() bool {
	r.Tree = *models.CreateFsTreeRoot(r.BackupArchive.WatchDir, r.BackupArchive.TempDir)
	r.MatchRule = regexp.MustCompile(r.BackupArchive.Name + "_20\\d{6}_\\d{6}\\.zip")
	fmt.Println("手动进行操作执行")
	return true
}

func (r *ManualBackupCommand) String() string {
	return conf.CommandNames.ManualBackupCommand
}

// 手动与自动备份命令
type ManualAndAutoBackupCommand struct {
	*management.ShareData
}

func (r *ManualAndAutoBackupCommand) Execute() bool {
	r.Tree = *models.CreateFsTreeRoot(r.BackupArchive.WatchDir, r.BackupArchive.TempDir)
	r.MatchRule = regexp.MustCompile(r.BackupArchive.Name + "_20\\d{6}_\\d{6}\\.zip")
	fmt.Println("双重模式执行")
	return true
}

func (r *ManualAndAutoBackupCommand) String() string {
	return conf.CommandNames.ManualAndAutoBackupCommand
}
