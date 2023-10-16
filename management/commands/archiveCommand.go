package commands

import (
	"BRGS/conf"
	"BRGS/management"
	"BRGS/pkg/tools"
	"fmt"
	"log"
)

// CompressedArchiveCommand 压缩存档命令
type CompressedArchiveCommand struct {
}

func (c *CompressedArchiveCommand) Execute() bool {
	fmt.Println("压缩执行")
	return management.MonServer.ZipArchive()
}

func (c *CompressedArchiveCommand) String() string {
	return conf.CommandNames.CompressedArchive
}

// BackupFilesWithHardLinkCommand  硬链接备份命令
type BackupFilesWithHardLinkCommand struct {
}

func (c *BackupFilesWithHardLinkCommand) Execute() bool {
	fmt.Println("硬链接执行")
	return management.MonServer.HardLinkArchive()
}

func (c *BackupFilesWithHardLinkCommand) String() string {
	return conf.CommandNames.HardLinkArchive
}

// RestoreFileFromArchive 从存档还原命令
type RestoreFileFromArchive struct {
}

func (r *RestoreFileFromArchive) Execute() bool {
	fmt.Println("从压缩文件还原执行")
	archives := management.MonServer.LoadArchives()
	archiveInfoList := make([]string, 0, len(archives))
	for _, archive := range archives {
		archiveInfoList = append(archiveInfoList, archive["type"]+"\t"+archive["name"])
	}
	if len(archives) == 0 {
		fmt.Println("无压缩文件")
		return true
	} else {
		if index := tools.CommandMenu(true, archiveInfoList...); index >= 0 {
			if ok := management.MonServer.RecoverFormArchive(archives[index]["name"]); ok {
				log.Printf("还原%s失败", archives[index])
				return false
			} else {
				log.Printf("还原%s成功", archives[index])
				return true
			}
		}
	}
	return true
}

func (r *RestoreFileFromArchive) String() string {
	return conf.CommandNames.RestoreFileFromArchive
}
