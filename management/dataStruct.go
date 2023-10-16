package management

import (
	"BRGS/conf"
)

var ExcelHeadTranslateDic = map[string]string{
	"name":            conf.ExcelTranslateConf.Name,
	"watchDir":        conf.ExcelTranslateConf.WatchDir,
	"tempDir":         conf.ExcelTranslateConf.TempDir,
	"archiveDir":      conf.ExcelTranslateConf.ArchiveDir,
	"backupDir":       conf.ExcelTranslateConf.BackupDir,
	"archiveInterval": conf.ExcelTranslateConf.ArchiveInterval,
	"syncInterval":    conf.ExcelTranslateConf.SyncInterval,
}

var ExcelTipDic = map[string]string{
	"name":            conf.ExcelTipConf.Name,
	"watchDir":        conf.ExcelTipConf.WatchDir,
	"tempDir":         conf.ExcelTipConf.TempDir,
	"archiveDir":      conf.ExcelTipConf.ArchiveDir,
	"backupDir":       conf.ExcelTipConf.BackupDir,
	"archiveInterval": conf.ExcelTipConf.ArchiveInterval,
	"syncInterval":    conf.ExcelTipConf.SyncInterval,
}

var ExcelHeadOrder = []string{
	"name",
	"watchDir",
	"tempDir",
	"archiveDir",
	"backupDir",
	"archiveInterval",
	"syncInterval",
}
