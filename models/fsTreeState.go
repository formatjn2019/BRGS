package models

import (
	"BRGS/pkg/e"
	"sync"
)

// 末尾数字为方向
// 小到大 末尾为1
// 大到小 末尾为0
const (
	FSTREE_WATCH = 2 << iota
	FSTREE_BACKUP
	FSTREE_ARCHIVE
	FSTREE_RECOVER
	FSTREE_WATCH_TO_BACKUP  = FSTREE_WATCH | FSTREE_BACKUP | 1
	FSTREE_BACKUP_TO_WATCH  = FSTREE_WATCH | FSTREE_BACKUP
	FSTREE_WATCH_TO_RECOVER = FSTREE_WATCH | FSTREE_RECOVER | 1
	FSTREE_RECOVER_TO_WATCH = FSTREE_WATCH | FSTREE_RECOVER
	FSTREE_WATCH_TO_ARCHIVE = FSTREE_WATCH | FSTREE_ARCHIVE | 1
	FSTREE_ARCHIVE_TO_WATCH = FSTREE_WATCH | FSTREE_ARCHIVE
)

type FstreeType struct {
	state int
	mx    sync.Mutex
}

// 判断状态是否可以进行该状态的转换
func (f *FstreeType) canChange(from, to int) bool {
	COVER_NUM := from | to
	if from < to {
		COVER_NUM |= 1
	}
	switch COVER_NUM {
	case FSTREE_WATCH_TO_BACKUP, FSTREE_WATCH_TO_ARCHIVE, FSTREE_WATCH_TO_RECOVER,
		FSTREE_BACKUP_TO_WATCH, FSTREE_RECOVER_TO_WATCH, FSTREE_ARCHIVE_TO_WATCH:
		return true
	default:
		return false
	}
}

func (f *FstreeType) State() int {
	return f.state
}

func (f *FstreeType) changeState(target int) (int, bool) {
	f.mx.Lock()
	defer f.mx.Unlock()
	if f.canChange(f.state, target) {
		f.state = target
		return f.state, true
	}
	return f.state, false
}

// 切换为备份状态
func (f *FstreeType) ChangeToBackup() (int, bool) {
	return f.changeState(FSTREE_BACKUP)
}

// 切换为监控状态
func (f *FstreeType) ChangeToWatch() (int, bool) {
	return f.changeState(FSTREE_WATCH)
}

// 切换为还原状态
func (f *FstreeType) ChangeToRecover() (int, bool) {
	return f.changeState(FSTREE_RECOVER)
}

// 翻译当前状态
func (f *FstreeType) Translate(state int) string {
	switch state {
	case FSTREE_ARCHIVE:
		return "FSTREE_ARCHIVE"
	case FSTREE_BACKUP:
		return "FSTREE_BACKUP"
	case FSTREE_RECOVER:
		return "FSTREE_RECOVER"
	case FSTREE_WATCH:
		return "FSTREE_WATCH"
	default:
		panic(e.ERROR_OPERATION_UNSUPPORT)
	}
}
