package models

import (
	"BRGS/pkg/e"
	"sync"
)

// 末尾数字为方向
// 小到大 末尾为1
// 大到小 末尾为0
const (
	FsTreeWatch = 2 << iota
	FsTreeBackup
	FsTreeArchive
	FsTreeRecover
	FsTreeWatchToBackup  = FsTreeWatch | FsTreeBackup | 1
	FsTreeBackupToWatch  = FsTreeWatch | FsTreeBackup
	FsTreeWatchToRecover = FsTreeWatch | FsTreeRecover | 1
	FsTreeRecoverToWatch = FsTreeWatch | FsTreeRecover
	FsTreeWatchToArchive = FsTreeWatch | FsTreeArchive | 1
	FsTreeArchiveToWatch = FsTreeWatch | FsTreeArchive
)

type FsTreeType struct {
	state int
	mx    sync.Mutex
}

// 判断状态是否可以进行该状态的转换
func (f *FsTreeType) canChange(from, to int) bool {
	coverNum := from | to
	if from < to {
		coverNum |= 1
	}
	switch coverNum {
	case FsTreeWatchToBackup, FsTreeWatchToArchive, FsTreeWatchToRecover,
		FsTreeBackupToWatch, FsTreeRecoverToWatch, FsTreeArchiveToWatch:
		return true
	default:
		return false
	}
}

func (f *FsTreeType) State() int {
	return f.state
}

func (f *FsTreeType) changeState(target int) (int, bool) {
	f.mx.Lock()
	defer f.mx.Unlock()
	if f.canChange(f.state, target) {
		f.state = target
		return f.state, true
	}
	return f.state, false
}

// ChangeToBackup 切换为备份状态
func (f *FsTreeType) ChangeToBackup() (int, bool) {
	return f.changeState(FsTreeBackup)
}

// ChangeToWatch 切换为监控状态
func (f *FsTreeType) ChangeToWatch() (int, bool) {
	return f.changeState(FsTreeWatch)
}

// ChangeToRecover 切换为还原状态
func (f *FsTreeType) ChangeToRecover() (int, bool) {
	return f.changeState(FsTreeRecover)
}

// ChangeToArchive 切换为存档状态
func (f *FsTreeType) ChangeToArchive() (int, bool) {
	return f.changeState(FsTreeArchive)
}

// Translate 翻译当前状态
func (f *FsTreeType) Translate(state int) string {
	switch state {
	case FsTreeArchive:
		return "FsTreeArchive"
	case FsTreeBackup:
		return "FsTreeBackup"
	case FsTreeRecover:
		return "FsTreeRecover"
	case FsTreeWatch:
		return "FsTreeWatch"
	default:
		panic(e.ErrorOperationUnsupport)
	}
}
