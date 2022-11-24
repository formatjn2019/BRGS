package models

import "sync"

// 末尾数字为方向
// 小到大 末尾为1
// 大到小 末尾为0
const (
	FSTREE_WATCH = 2 << iota
	FSTREE_BACKUP
	FSTREE_RECOVER
	FSTREE_WATCH_TO_BACKUP  = FSTREE_WATCH | FSTREE_BACKUP | 1
	FSTREE_BACKUP_TO_WATCH  = FSTREE_WATCH | FSTREE_BACKUP
	FSTREE_WATCH_TO_RECOVER = FSTREE_WATCH | FSTREE_RECOVER | 1
	FSTREE_RECOVER_TO_WATCH = FSTREE_WATCH | FSTREE_RECOVER
)

const (
	FSTREE_OP_BACKUP_PREPERE = iota
	FSTREE_OP_BACKUP_END
	FSTREE_OP_RECOVER_PREPERE
	FSTREE_OP_RECOVER_END
)

type FstreeType struct {
	state int
	mx    sync.Mutex
}

func (f *FstreeType) changeState(nstat int) {
	f.mx.Lock()
	if f.state != nstat && f.canChange(f.state, nstat) {
		f.state = nstat
	}
	f.mx.Unlock()
}

//判断状态是否可以进行该状态的转换
func (f *FstreeType) canChange(from, to int) bool {
	COVER_NUM := from | to
	if from < to {
		COVER_NUM |= 1
	}
	return COVER_NUM == FSTREE_BACKUP_TO_WATCH || COVER_NUM == FSTREE_WATCH_TO_BACKUP ||
		COVER_NUM == FSTREE_WATCH_TO_RECOVER || COVER_NUM == FSTREE_RECOVER_TO_WATCH
}

func (f *FstreeType) State() int {
	return f.state
}

// 切换为备份状态
func (f *FstreeType) Backup() bool {
	f.changeState(FSTREE_BACKUP)
	return true
}

// 切换为监控状态
func (f *FstreeType) Watch() bool {
	f.changeState(FSTREE_WATCH)
	return true
}

// 切换为还原状态
func (f *FstreeType) Recover() bool {
	f.changeState(FSTREE_RECOVER)
	return true
}
