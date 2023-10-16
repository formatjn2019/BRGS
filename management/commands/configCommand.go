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
	"os"
	"strconv"
	"strings"
)

// ReadConfigCommand 读取命令
type ReadConfigCommand struct {
}

func (r *ReadConfigCommand) Execute() bool {
	configs, err := utils.ReadCsvAsDictAndTranslate("config.csv", tools.DictReverse(management.ExcelHeadTranslateDic))
	rules := make([]management.BackupArchive, 0)
	names := map[string]bool{}
	if err == nil {
		for line, config := range configs {
			fmt.Println(line+1, "\t", config)
			// 结构转换
			archiveInterval, te := strconv.Atoi(config["archiveInterval"])
			if te != nil {
				err = errors.New("存档间隔非数字")
				goto errorLog
			}
			syncInterval, te := strconv.Atoi(config["syncInterval"])
			if te != nil {
				err = errors.New("同步间隔非数字")
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

			if checkErrors := rule.CheckConfig(); len(checkErrors) > 0 {
				err = errors.New(fmt.Sprintf("检查%d行出错\t", line+2) + strings.Join(checkErrors, "\n"))
				goto errorLog
			}
			if ok, _ := names[rule.Name]; ok {
				err = e.TranslateToError(e.ErrorSameName, fmt.Sprintf("检查%d行出错", line+2))
				goto errorLog
			}
			names[rule.Name] = true
			rules = append(rules, rule)
		}
		//生成命令
		for _, rule := range rules {
			hd := `.\BRGS.exe `
			tail := fmt.Sprintf(`-n %s -wd %s -td %s -ad %s -ai %d -si %d`, rule.Name, rule.WatchDir, rule.TempDir, rule.ArchiveDir, rule.ArchiveInterval/60, rule.SyncInterval/60)
			nameCommandDic := map[string]string{
				rule.Name + "_web.cmd":        hd + " -s " + tail,
				rule.Name + "_manual.cmd":     hd + " -m " + tail,
				rule.Name + "_web_manual.cmd": hd + " -s -m " + tail,
			}
			for name, context := range nameCommandDic {
				err = os.WriteFile(name, []byte(context+"\npause"), 0755)
				if err != nil {
					goto errorLog
				}
			}
		}
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
