package management

import (
	"BRGS/models"
	"io/fs"
	"regexp"
)

// Command 命令接口
type Command interface {
	Execute() bool
	String() string
}

// ShareData 数据共享
type ShareData struct {
	BackupArchive    BackupArchive
	Tree             models.FsTreeRoot
	LinkedBackupFile map[string]map[string]fs.FileInfo
	MatchRule        *regexp.Regexp
	ServerChan       chan struct{}
}
