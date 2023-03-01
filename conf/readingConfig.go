package conf

import (
	"BRGS/pkg/e"
	"fmt"
	"os"

	"github.com/go-ini/ini"
)

type ExcelConfig struct {
	ArchiveDir      string
	ArchiveInterval string
	Name            string
	SyncInterval    string
	TempDir         string
	WatchDir        string
}

var ExcelTranslateConf = &ExcelConfig{}
var ExcelTipConf = &ExcelConfig{}

type CommandName struct {
	// 配置命令
	ReadConfigCommand            string
	GenerateConfigDefaultCommand string
	// 压缩与解压命令
	CompressedArchive      string
	RestoreFileFromArchive string
	// 备份命令
	BackupCommand     string
	ResetBackup       string
	RestoreBackup     string
	StopBackupCommand string
	// 备份控制命令
	AutoBackupCommand          string
	ManualAndAutoBackupCommand string
	ManualBackupCommand        string
	// 控制命令
	ExitCommand        string
	StartServerCommand string
	StopServerCommand  string
}

var CommandNames = &CommandName{}

type Server struct {
	Port string
}

var ServerConf = &Server{}

func init() {
	cfg, err := ini.Load("config.ini")

	if err != nil {
		fmt.Printf(e.TranslateToError(e.ErrorReadConfig, "Fail to read file").Error())
		os.Exit(1)
	}
	cfg.Section("excel_head_ch").MapTo(ExcelTranslateConf)
	cfg.Section("excel_default_tip_ch").MapTo(ExcelTipConf)
	cfg.Section("command_name").MapTo(CommandNames)
	cfg.Section("server").MapTo(ServerConf)
}
