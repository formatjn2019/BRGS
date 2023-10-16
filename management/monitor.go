package management

import (
	"BRGS/models"
	"BRGS/pkg/e"
	"BRGS/pkg/tools"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"time"
)

type Monitor struct {
	ftr *models.FsTreeRoot
	BackupArchive
	linkedFilesDic     map[string]map[string]fs.FileInfo
	running            bool
	state              int
	reloadArchivesFlag bool
	rule               *regexp.Regexp
}

const (
	StateWatch = iota
	stateSync
	stateArchive
	statePause
)

func CreateMonitor(ba BackupArchive) *Monitor {
	root := models.CreateFsTreeRoot(ba.WatchDir, ba.TempDir)
	monitor := &Monitor{
		BackupArchive: ba,
		ftr:           root,
		rule:          tools.GenerateRule(ba.Name),
	}
	var linkedFilesDic map[string]map[string]fs.FileInfo
	var dirs []string
	if files, err := os.ReadDir(ba.ArchiveDir); err == nil {
		for _, info := range files {
			if info.IsDir() && monitor.rule.Match([]byte(info.Name())) {
				dirs = append(dirs, path.Join(ba.ArchiveDir, info.Name()))
			}
		}
	}
	linkedFilesDic, err := tools.ScanFilesToDic(dirs...)
	if err != nil {
		panic(err)
	}
	monitor.linkedFilesDic = linkedFilesDic
	return monitor
}

var MonServer *Monitor

func (m *Monitor) CanOperate() bool {
	fmt.Println(m.ftr.State())
	return m.ftr.State() == models.FsTreeWatch
}

func (m *Monitor) State() int {
	if m.state == statePause {
		return statePause
	} else {
		switch m.ftr.State() {
		case models.FsTreeWatch:
			return StateWatch
		case models.FsTreeArchive:
			return stateArchive
		case models.FsTreeRecover, models.FsTreeBackup:
			return stateSync
		default:
			panic(e.TranslateToError(e.ErrorFunction, "状态异常"))
		}
	}
}

// ZipArchive 使用zip备份至存档
func (m *Monitor) ZipArchive() bool {
	m.reloadArchivesFlag = true
	err := m.ftr.CommandInput(models.FsTreeOpArchive)
	if err != nil {
		return false
	}
	zipDst := path.Join(m.ArchiveDir, tools.GenerateNameByTime(m.Name)) + ".zip"
	if err = tools.WriteZip(zipDst, tools.WalkDir(m.TempDir)); err != nil {
		return false
	}
	err = m.ftr.CommandInput(models.FsTreeOpWatch)
	if err != nil {
		return false
	}
	return true
}

// HardLinkArchive 使用硬链接备份至存档
func (m *Monitor) HardLinkArchive() bool {
	m.reloadArchivesFlag = true
	err := m.ftr.CommandInput(models.FsTreeOpArchive)
	if err != nil {
		return false
	}
	dist := path.Join(m.ArchiveDir, tools.GenerateNameByTime(m.Name))
	// 计算所有Uid
	pathUidDic, err := tools.CalculateAllUid(m.TempDir, false)
	if err != nil {
		return false
	}
	err = tools.SyncAllFileWithHardLink(m.WatchDir, dist, pathUidDic, m.linkedFilesDic)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	err = m.ftr.CommandInput(models.FsTreeOpWatch)
	if err != nil {
		return false
	}
	return true
}

// RecoverFormArchive 从存档还原
func (m *Monitor) RecoverFormArchive(name string) bool {
	fmt.Println(name)
	err := m.ftr.CommandInput(models.FsTreeOpArchive)
	if err != nil {
		return false
	}

	archivePath := filepath.Join(m.ArchiveDir, name)
	if stat, err := os.Stat(archivePath); err != nil {
		log.Fatal("还原失败", archivePath)
		return false
	} else {
		if stat.IsDir() {
			err = tools.CloneDir(archivePath, m.WatchDir)
			if err != nil {
				return false
			}
		} else {
			err = tools.RecoverFromArchive(archivePath, m.WatchDir)
			if err != nil {
				return false
			}
		}
	}
	m.reloadArchivesFlag = true
	err = m.ftr.CommandInput(models.FsTreeOpWatch)
	if err != nil {
		return false
	}
	// 重新扫描改动
	m.Scanning()
	return true
}

func (m *Monitor) NeedReload() bool {
	return m.reloadArchivesFlag
}

func (m *Monitor) LoadArchives() []map[string]string {
	archives := make([]map[string]string, 0)
	if files, err := os.ReadDir(m.ArchiveDir); err == nil {
		for _, info := range files {
			if m.rule.Match([]byte(info.Name())) {
				archiveItem := map[string]string{
					"name": info.Name(),
				}
				if info.IsDir() {
					archiveItem["type"] = "硬链接"
				} else {
					archiveItem["type"] = "zip文件"
				}

				archives = append(archives, archiveItem)
			}
		}
	} else {
		log.Printf("扫描压缩文件夹%s失败:\t%v\n", m.ArchiveDir, err)
		return []map[string]string{}
	}
	m.reloadArchivesFlag = false
	fmt.Println(archives)
	return archives
}

// Scanning 重新扫描
func (m *Monitor) Scanning() bool {
	// 重新扫描改动
	m.ftr.ScanFolder(m.WatchDir, "")
	return true
}

// Backup 备份
func (m *Monitor) Backup() bool {
	err := m.ftr.CommandInput(models.FsTreeOpBackup)
	if err != nil {
		return false
	}
	return true
}

// Recover 还原
func (m *Monitor) Recover() bool {
	err := m.ftr.CommandInput(models.FsTreeOpRecover)
	if err != nil {
		return false
	}
	return true
}

// Pause 暂停
func (m *Monitor) Pause() bool {
	m.state = statePause
	err := m.ftr.CommandInput(models.FsTreeOpArchive)
	if err != nil {
		return false
	}
	return true
}

// Continue 继续
func (m *Monitor) Continue() bool {
	m.state = StateWatch
	err := m.ftr.CommandInput(models.FsTreeOpWatch)
	if err != nil {
		return false
	}
	return true
}

// Run 自动备份
func (m *Monitor) Run() {
	m.running = true
	if m.SyncInterval != 0 {
		go func() {
			for m.running {
				time.Sleep(time.Duration(m.SyncInterval) * time.Second)
				succeedFlag := false
				// 不成功，则重试
				for !succeedFlag {
					if m.CanOperate() {
						succeedFlag = m.Backup()
					}
					time.Sleep(10 * time.Second)
				}
			}
		}()
	}
}

func (m *Monitor) Stop() {

}
