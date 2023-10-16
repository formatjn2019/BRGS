package management

import (
	"BRGS/pkg/e"
	"BRGS/pkg/tools"
	"fmt"
	"strconv"
)

type BackupArchive struct {
	ArchiveDir      string
	BackupDir       string
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
		// "备份目录":   b.backupDir,
		"backupDir": b.BackupDir,
		// "自动存档间隔(分钟)": b.archiveInterval,
		"archiveInterval": strconv.Itoa(b.ArchiveInterval / 60),
		// "自动同步间隔(分钟)":  b.syncInterval,
		"syncInterval": strconv.Itoa(b.SyncInterval / 60),
	}
}

// 同步与备份间隔（分钟）
const (
	ArchiveIntervalAllowZero = true
	MaxArchiveInterval       = 30
	MinSyncInterval          = 2
	SyncIntervalAllowZero    = true
	MaxSyncInterval          = 10
	MinArchiveInterval       = 1
)

var notNullCheck = tools.GenerateNotNullCheck()
var pathCheck = tools.GeneratePathCheck()
var checkMap = map[string][]tools.Check{
	// "名称"
	"name": {notNullCheck},
	// "存档目录"
	"watchDir": {notNullCheck, pathCheck},
	// "中转文件目录"
	"tempDir": {pathCheck},
	// "压缩文件存储目录"
	"archiveDir": {notNullCheck, pathCheck},
	// "自动存档间隔(分钟)"
	"archiveInterval": {notNullCheck, tools.GenerateRangeCheck(ArchiveIntervalAllowZero, MinArchiveInterval, MaxArchiveInterval)},
	// "自动同步间隔(分钟)"
	"syncInterval": {notNullCheck, tools.GenerateRangeCheck(SyncIntervalAllowZero, MinSyncInterval, MaxSyncInterval)},
}

// CheckConfig 配置检查
func (b *BackupArchive) CheckConfig() []string {
	var result []string
	ckDic := b.GetConfigDic()
	//通用规则检查
	for item, value := range ckDic {
		for _, check := range checkMap[item] {
			if err := check(value); err != nil {
				result = append(result, item+"\t"+err.Error())
			}
		}
	}
	if len(result) > 0 {
		return result
	}
	//隐藏规则检查
	if b.ArchiveInterval < b.SyncInterval {
		result = append(result, e.TranslateToError(e.ErrorInterval, "输入间隔小于输出间隔").Error())
	}
	return result
}

func (b *BackupArchive) String() string {
	return fmt.Sprintf("名称:%s\n存档目录:%s\n中转文件目录:%s\n压缩文件目录:%s\n文件同步间隔:%d分钟\n压缩存档间隔:%d分钟\n", b.Name, b.WatchDir, b.TempDir, b.ArchiveDir, b.SyncInterval/60, b.ArchiveInterval/60)
}
