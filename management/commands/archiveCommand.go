package commands

import (
	"BRGS/management"
	"BRGS/models"
	"BRGS/pkg/tools"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

// 压缩存档命令
type CompressedArchive struct {
	*management.ShareData
}

func (c *CompressedArchive) Execute() bool {
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

func (c *CompressedArchive) String() string {
	return fmt.Sprintf("将同步文件夹压缩为压缩包")
}

// 退出命令
type Exit struct {
	*management.ShareData
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
	*management.ShareData
}

func (r *ResetBackup) Execute() bool {
	r.Tree = *models.CreateFsTreeRoot(r.BackupArchive.WatchDir, r.BackupArchive.TempDir)
	fmt.Println("重置执行")
	return true
}

func (r *ResetBackup) String() string {
	return fmt.Sprintf("重置备份文件夹")
}

// 还原命令
type ResoreBackup struct {
	*management.ShareData
}

func (r *ResoreBackup) Execute() bool {
	if r.Tree.RecoverFiles() {
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
	return fmt.Sprintf("从压缩文件还原执行文件夹")
}
