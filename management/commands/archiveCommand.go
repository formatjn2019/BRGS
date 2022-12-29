package commands

import (
	"BRGS/conf"
	"BRGS/management"
	"BRGS/pkg/tools"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"
)

// 压缩存档命令
type CompressedArchiveCommand struct {
	*management.ShareData
}

func (c *CompressedArchiveCommand) Execute() bool {
	fmt.Println("压缩执行")
	nowTime := time.Now()
	// 根据时间生成压缩包名称
	fileName := fmt.Sprintf("%s_%s.zip", c.BackupArchive.Name, nowTime.Format("20060102_150405"))
	zipPath := filepath.Join(c.BackupArchive.ArchiveDir, fileName)
	fmt.Println(fileName)
	fmt.Println(zipPath)
	if err := tools.WriteZip(zipPath, tools.WalkDir(c.BackupArchive.TempDir)); err != nil {
		return false
	}
	return true
}

func (c *CompressedArchiveCommand) String() string {
	return conf.CommandNames.CompressedArchive
}

// 从存档还原命令
type ResoreFileFromArchive struct {
	*management.ShareData
}

func (r *ResoreFileFromArchive) Execute() bool {
	fmt.Println("从压缩文件还原执行")
	archives := make([]string, 0)
	if files, err := ioutil.ReadDir(r.BackupArchive.ArchiveDir); err == nil {
		for _, info := range files {
			if r.MatchRule.Match([]byte(info.Name())) {
				archives = append(archives, info.Name())
			}
		}
		if len(archives) == 0 {
			fmt.Println("无压缩文件")
			return true
		} else {
			if index := tools.CommandMenu(true, archives...); index >= 0 {
				if err := tools.RecoverFromArchive(filepath.Join(r.BackupArchive.ArchiveDir, archives[index]), r.BackupArchive.WatchDir); err != nil {
					log.Printf("还原%s失败", archives[index])
					return false
				} else {
					log.Printf("还原%s成功", archives[index])
					return true
				}
			}
		}

	} else {
		log.Printf("扫描压缩文件夹%s失败:\t%v\n", r.BackupArchive.ArchiveDir, err)
	}
	// models.RecoverFromArchive("")
	return true
}

func (r *ResoreFileFromArchive) String() string {
	return conf.CommandNames.ResoreFileFromArchive
}
