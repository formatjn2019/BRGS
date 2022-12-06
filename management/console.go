package management

import (
	"BRGS/models"
	"BRGS/pkg/e"
	"BRGS/pkg/tools"
	"BRGS/pkg/utils"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/go-ini/ini"
)

func (b *BackupArchive) GetConfigDic() map[string]string {
	return map[string]string{
		// "名称":         b.name,
		"name": b.name,
		// "存档目录":       b.watchDir,
		"watchDir": b.watchDir,
		// "中转文件目录":     b.tempDir,
		"tempDir": b.tempDir,
		// "压缩文件存储目录":   b.archiveDir,
		"archiveDir": b.archiveDir,
		// "自动存档间隔(分钟)": strconv.Itoa(b.archiveInterval),
		"archiveInterval": strconv.Itoa(b.archiveInterval / 60),
		// "自动同步间隔(分钟)": strconv.Itoa(b.syncInterval),
		"syncInterval": strconv.Itoa(b.syncInterval / 60),
	}
}

func (b *BackupArchive) String() string {
	return fmt.Sprintf("名称:%s\n存档目录:%s\n中转文件目录:%s\n压缩文件目录:%s\n文件同步间隔:%d分钟\n压缩存档间隔:%d分钟\n", b.name, b.watchDir, b.tempDir, b.archiveDir, b.syncInterval/60, b.archiveInterval/60)
}

// 命令接口
type Command interface {
	Execute() bool
	String() string
}

type shareData struct {
	ba        BackupArchive
	tree      models.FSTreeRoot
	matchRule *regexp.Regexp
}

// 读取命令
type ReadConfigCommand struct {
	*shareData
}

func (r *ReadConfigCommand) Execute() bool {
	configs, err := utils.ReadCsvAsDictAndTranslate("config.csv", tools.DictReverse(EXCEL_HEAD_TRANSLATE_DIC))

	// 非空核验
	notNullCheck := tools.GenerateNotNullCheck()
	// 路径核验
	pathCheck := tools.GeneratePathCheck()
	//通用检查
	checkMap := map[string]([]tools.Check){
		// "名称"
		"name": {notNullCheck},
		// "存档目录"
		"watchDir": {notNullCheck, pathCheck},
		// "中转文件目录"
		"tempDir": {pathCheck},
		// "压缩文件存储目录"
		"archiveDir": {notNullCheck, pathCheck},
		// "自动存档间隔(分钟)"
		"archiveInterval": {notNullCheck, tools.GenerateRangeCheck(true, 2, 120)},
		// "自动同步间隔(分钟)"
		"syncInterval": {notNullCheck, tools.GenerateRangeCheck(true, 1, 30)},
	}
	println(checkMap)
	rules := make([]BackupArchive, 0)
	names := map[string]bool{}
	if err == nil {
		for line, config := range configs {
			fmt.Println(line+1, "\t", config)
			//通用规则检查
			for row := 0; row < len(EXCEL_HEAD_ORDER); row++ {
				text := config[EXCEL_HEAD_ORDER[row]]
				for _, check := range checkMap[EXCEL_HEAD_ORDER[row]] {
					if err = check(text); err != nil {
						err = errors.New(fmt.Sprintf("检查%d行\t%d列出错,内容为:%s,错误为：%s", line+2, row+1, text, err))
						goto errorLog
					}
				}
			}
			//隐藏规则检查
			archiveInterval, _ := strconv.Atoi(config["archiveInterval"])
			syncInterval, _ := strconv.Atoi(config["syncInterval"])
			if archiveInterval < syncInterval {
				err = e.TranslateToError(e.ERROR_INTERVAL, fmt.Sprintf("检查%d行出错", line+2))
				goto errorLog
			}
			rule := BackupArchive{
				name:            config["name"],
				watchDir:        config["watchDir"],
				tempDir:         config["tempDir"],
				archiveDir:      config["archiveDir"],
				archiveInterval: archiveInterval * 60,
				syncInterval:    syncInterval * 60,
			}
			fmt.Println(rule.String())
			if ok, _ := names[rule.name]; ok {
				err = e.TranslateToError(e.ERROR_SAME_NAME, fmt.Sprintf("检查%d行出错", line+2))
				goto errorLog
			}
			names[rule.name] = true
			rules = append(rules, rule)
		}
		configStr := make([]string, 0)
		for _, rule := range rules {
			configStr = append(configStr, rule.String())
		}
		if selectIndex := CommandMenu(true, configStr...); selectIndex != -1 {
			r.ba = rules[selectIndex]
		} else {
			goto errorLog
		}
		fmt.Println("当前规则为：\n", r.ba)
	} else {
		goto errorLog
	}
	return true

errorLog:
	log.Fatal(err)
	return false
}

func (r *ReadConfigCommand) String() string {
	return fmt.Sprintf("读取配置文件")
}

// 生成默认配置文件命令
type GenerateConfigDefaultCommand struct {
	*shareData
}

func (r *GenerateConfigDefaultCommand) Execute() bool {
	cfg, err := ini.Load("config.ini")

	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	tipDic := map[string]string{}
	for k, v := range EXCEL_HEAD_TRANSLATE_DIC {
		tipDic[v] = cfg.Section("excel_default_tip_ch").Key(k).String()
	}

	translatedHead, _ := tools.TranslateList(EXCEL_HEAD_ORDER, EXCEL_HEAD_TRANSLATE_DIC)
	utils.WriteCsvWithDict("config_default(需要改名为config才能使用).csv", []map[string]string{tipDic}, translatedHead...)
	fmt.Println("写入默认配置文件成功")
	return true
}

func (r *GenerateConfigDefaultCommand) String() string {
	return fmt.Sprintf("生成默认配置文件")
}

// 停止备份命令
type StopBackupCommond struct {
	*shareData
}

func (b *StopBackupCommond) Execute() bool {
	fmt.Println("停止同步执行")
	return true
}

func (b *StopBackupCommond) String() string {
	return fmt.Sprintf("停止同步文件夹内容")
}

// 备份命令
type BackupCommand struct {
	*shareData
}

func (b *BackupCommand) Execute() bool {
	if b.tree.BackupFiles() {
		fmt.Println("同步执行成功")
		return true
	} else {
		fmt.Println("同步执行失败")
		return false
	}
}

func (b *BackupCommand) String() string {
	return fmt.Sprintf("同步文件夹内容")
}

// 压缩存档命令
type CompressedArchive struct {
	*shareData
}

func (c *CompressedArchive) Execute() bool {
	fmt.Println("压缩执行")
	nowTime := time.Now()
	// 根据时间生成压缩包名称
	fileName := fmt.Sprintf("%s_%s.zip", c.ba.name, nowTime.Format("20060102_150405"))
	zipPath := filepath.Join(c.ba.archiveDir, fileName)
	fmt.Println(fileName)
	fmt.Println(zipPath)
	if err := tools.WriteZip(zipPath, tools.WalkDir(c.ba.tempDir)); err != nil {
		return false
	}
	return true
}

func (c *CompressedArchive) String() string {
	return fmt.Sprintf("将同步文件夹压缩为压缩包")
}

// 退出命令
type Exit struct {
	*shareData
}

func (e *Exit) Execute() bool {
	log.Println("退出执行")
	os.Exit(0)
	return true
}

func (e *Exit) String() string {
	return fmt.Sprintf("退出")
}

// 重置备份命令
type ResetBackup struct {
	*shareData
}

func (r *ResetBackup) Execute() bool {
	r.tree = *models.CreateFsTreeRoot(r.ba.watchDir, r.ba.tempDir)
	fmt.Println("重置执行")
	return true
}

func (r *ResetBackup) String() string {
	return fmt.Sprintf("重置备份文件夹")
}

// 还原命令
type ResoreBackup struct {
	*shareData
}

func (r *ResoreBackup) Execute() bool {
	if r.tree.RecoverFiles() {
		fmt.Println("同步执行成功")
		return true
	} else {
		fmt.Println("同步执行失败")
		return false
	}
}

func (r *ResoreBackup) String() string {
	return fmt.Sprintf("从同步文件夹还原文件")
}

// 从存档还原命令
type ResoreFileFromArchive struct {
	*shareData
}

func (r *ResoreFileFromArchive) Execute() bool {
	fmt.Println("从压缩文件还原执行")
	archives := make([]string, 0)
	if files, err := ioutil.ReadDir(r.ba.archiveDir); err == nil {
		for _, info := range files {
			if r.matchRule.Match([]byte(info.Name())) {
				archives = append(archives, info.Name())
			}
		}
		if len(archives) == 0 {
			fmt.Println("无压缩文件")
			return true
		} else {
			if index := CommandMenu(true, archives...); index >= 0 {
				if err := tools.RecoverFromArchive(filepath.Join(r.ba.archiveDir, archives[index]), r.ba.watchDir); err != nil {
					log.Printf("还原%s失败", archives[index])
					return false
				} else {
					log.Printf("还原%s成功", archives[index])
					return true
				}
			}
		}

	} else {
		log.Printf("扫描压缩文件夹%s失败:\t%v\n", r.ba.archiveDir, err)
	}
	// models.RecoverFromArchive("")
	return true
}

func (r *ResoreFileFromArchive) String() string {
	return fmt.Sprintf("从压缩文件还原执行文件夹")
}

// 自动备份命令
type AutoBackupCommand struct {
	*shareData
}

func (r *AutoBackupCommand) Execute() bool {
	r.tree = *models.CreateFsTreeRoot(r.ba.watchDir, r.ba.tempDir)
	return true
}

func (r *AutoBackupCommand) String() string {
	return fmt.Sprintf("根据配置自动备份")
}

// 手动备份命令
type ManualBackupCommand struct {
	*shareData
}

func (r *ManualBackupCommand) Execute() bool {
	r.tree = *models.CreateFsTreeRoot(r.ba.watchDir, r.ba.tempDir)
	r.matchRule = regexp.MustCompile(r.ba.name + "_20\\d{6}_\\d{6}\\.zip")
	fmt.Println("手动进行操作执行")
	return true
}

func (r *ManualBackupCommand) String() string {
	return fmt.Sprintf("手动备份")
}

// 手动与自动备份命令
type ManualAndAutoBackupCommand struct {
	*shareData
}

func (r *ManualAndAutoBackupCommand) Execute() bool {
	r.tree = *models.CreateFsTreeRoot(r.ba.watchDir, r.ba.tempDir)
	r.matchRule = regexp.MustCompile(r.ba.name + "_20\\d{6}_\\d{6}\\.zip")
	fmt.Println("双重模式执行")
	return true
}

func (r *ManualAndAutoBackupCommand) String() string {
	return fmt.Sprintf("手动与自动双模式")
}
