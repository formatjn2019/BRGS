package management

import (
	"BRGS/models"
	"fmt"
	"os"
	"strconv"

	"github.com/go-ini/ini"
)

var EXCEL_HEAD_TRANSLATE_DIC = map[string]string{
	"name":            "name",
	"watchDir":        "watchDir",
	"tempDir":         "tempDir",
	"archiveDir":      "archiveDir",
	"archiveInterval": "archiveInterval",
	"syncInterval":    "syncInterval",
}

var EXCEL_HEAD_ORDER = []string{
	"name",
	"watchDir",
	"tempDir",
	"archiveDir",
	"archiveInterval",
	"syncInterval",
}

type BackupArchive struct {
	Root            models.FSTreeRoot
	ArchiveDir      string
	Name            string
	TempDir         string
	WatchDir        string
	ArchiveInterval int
	SyncInterval    int
}

func (b *BackupArchive) GetConfigDic() map[string]string {
	return map[string]string{
		// "名称":         b.name,
		"name": b.Name,
		// "存档目录":       b.watchDir,
		"watchDir": b.WatchDir,
		// "中转文件目录":     b.tempDir,
		"tempDir": b.TempDir,
		// "压缩文件存储目录":   b.archiveDir,
		"archiveDir": b.ArchiveDir,
		// "自动存档间隔(分钟)": strconv.Itoa(b.archiveInterval),
		"archiveInterval": strconv.Itoa(b.ArchiveInterval / 60),
		// "自动同步间隔(分钟)": strconv.Itoa(b.syncInterval),
		"syncInterval": strconv.Itoa(b.SyncInterval / 60),
	}
}

func (b *BackupArchive) String() string {
	return fmt.Sprintf("名称:%s\n存档目录:%s\n中转文件目录:%s\n压缩文件目录:%s\n文件同步间隔:%d分钟\n压缩存档间隔:%d分钟\n", b.Name, b.WatchDir, b.TempDir, b.ArchiveDir, b.SyncInterval/60, b.ArchiveInterval/60)
}

func init() {
	cfg, err := ini.Load("config.ini")

	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	for k := range EXCEL_HEAD_TRANSLATE_DIC {
		EXCEL_HEAD_TRANSLATE_DIC[k] = cfg.Section("excel_head_ch").Key(k).String()
	}

}
