package management

import (
	"BRGS/models"
	"regexp"
)

// 命令接口
type Command interface {
	Execute() bool
	String() string
}

// 数据共享
type ShareData struct {
	BackupArchive BackupArchive
	Tree          models.FSTreeRoot
	MatchRule     *regexp.Regexp
	ServerChan    chan struct{}
}
