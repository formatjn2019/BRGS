package management

import "BRGS/models"

var EXCEL_HEAD_ORDER = []string{
	"名称",
	"存档目录",
	"中转文件目录",
	"压缩文件存储目录",
	"自动存档间隔(分钟)",
	"自动同步间隔(分钟)",
}

type BackupArchive struct {
	root            models.FSTreeRoot
	archiveDir      string
	name            string
	tempDir         string
	watchDir        string
	archiveInterval int
	syncInterval    int
}
