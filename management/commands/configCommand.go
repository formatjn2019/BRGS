package commands

import (
	"BRGS/conf"
	"BRGS/management"
	"BRGS/pkg/e"
	"BRGS/pkg/tools"
	"BRGS/pkg/utils"
	"errors"
	"fmt"
	"log"
	"strconv"
)

// ReadConfigCommand 读取命令
type ReadConfigCommand struct {
	*management.ShareData
}

func (r *ReadConfigCommand) Execute() bool {
	configs, err := utils.ReadCsvAsDictAndTranslate("config.csv", tools.DictReverse(management.ExcelHeadTranslateDic))

	// 非空核验
	notNullCheck := tools.GenerateNotNullCheck()
	// 路径核验
	pathCheck := tools.GeneratePathCheck()
	//通用检查
	checkMap := map[string][]tools.Check{
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
	rules := make([]management.BackupArchive, 0)
	names := map[string]bool{}
	if err == nil {
		for line, config := range configs {
			fmt.Println(line+1, "\t", config)
			//通用规则检查
			for row := 0; row < len(management.ExcelHeadOrder); row++ {
				text := config[management.ExcelHeadOrder[row]]
				for _, check := range checkMap[management.ExcelHeadOrder[row]] {
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
				err = e.TranslateToError(e.ErrorInterval, fmt.Sprintf("检查%d行出错", line+2))
				goto errorLog
			}
			rule := management.BackupArchive{
				Name:            config["name"],
				WatchDir:        config["watchDir"],
				TempDir:         config["tempDir"],
				ArchiveDir:      config["archiveDir"],
				ArchiveInterval: archiveInterval * 60,
				SyncInterval:    syncInterval * 60,
			}
			fmt.Println(rule.String())
			if ok, _ := names[rule.Name]; ok {
				err = e.TranslateToError(e.ErrorSameName, fmt.Sprintf("检查%d行出错", line+2))
				goto errorLog
			}
			names[rule.Name] = true
			rules = append(rules, rule)
		}
		configStr := make([]string, 0)
		for _, rule := range rules {
			configStr = append(configStr, rule.String())
		}
		if selectIndex := tools.CommandMenu(true, configStr...); selectIndex != -1 {
			r.BackupArchive = rules[selectIndex]
		} else {
			goto errorLog
		}
		fmt.Println("当前规则为：\n", r.BackupArchive)
	} else {
		goto errorLog
	}
	return true

errorLog:
	log.Fatal(err)
	return false
}

func (r *ReadConfigCommand) String() string {
	return conf.CommandNames.ReadConfigCommand
}

// GenerateConfigDefaultCommand 生成默认配置文件命令
type GenerateConfigDefaultCommand struct {
	*management.ShareData
}

func (r *GenerateConfigDefaultCommand) Execute() bool {
	tipDic := map[string]string{}
	for k, v := range management.ExcelHeadTranslateDic {
		tipDic[v] = management.ExcelTipDic[k]
	}
	translatedHead, _ := tools.TranslateList(management.ExcelHeadOrder, management.ExcelHeadTranslateDic)
	err := utils.WriteCsvWithDict("config_default(需要改名为config才能使用).csv", []map[string]string{tipDic}, translatedHead...)
	if err != nil {
		fmt.Println("写入默认配置文件识别")
		return false
	}
	fmt.Println("写入默认配置文件成功")
	return true
}

func (r *GenerateConfigDefaultCommand) String() string {
	return conf.CommandNames.GenerateConfigDefaultCommand
}
