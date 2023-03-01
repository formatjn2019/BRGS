package models

import (
	"fmt"
)

// FsTreeNode 树的结点
type FsTreeNode struct {
	path     string      //相对路径
	parent   *FsTreeNode //父目录
	exist    bool        //是否存在
	isDir    bool        //是不是文件夹
	isAlter  bool        //是否修改
	synced   bool        //同步文件状态
	syncType bool        //同步文件类型(是否文件夹)
	subs     map[string]*FsTreeNode
}

// Delete 文件删除操作
func (ftn *FsTreeNode) Delete() {
	ftn.exist = false
	ftn.isAlter = true
	delete(ftn.parent.subs, ftn.path)
	if ftn.isDir {
		for _, subNode := range ftn.subs {
			subNode.Delete()
		}
	}
}

// Update 文件更新操作
func (ftn *FsTreeNode) Update() {
	ftn.exist = true
	ftn.isAlter = true
}

// ReverseSync 逆向文件同步操作
func (ftn *FsTreeNode) ReverseSync() {
	ftn.isAlter = false
	ftn.exist = ftn.synced
	ftn.isDir = ftn.exist
}

// Subs 输出
func (ftn *FsTreeNode) Subs() map[string]*FsTreeNode {
	return ftn.subs
}

// 代码格式化输出
func (ftn *FsTreeNode) String() string {
	return fmt.Sprintf("路径%40s\t 文件夹：%t\t更改状态%t\t 是否存在%t\t同步状态%t\t同步类型%t\t 子文件%d\n", ftn.path, ftn.isDir, ftn.isAlter, ftn.exist, ftn.synced, ftn.syncType, len(ftn.subs))
}

// Sync 文件同步操作
func (ftn *FsTreeNode) Sync() {
	ftn.isAlter = false
	ftn.synced = ftn.exist
	ftn.syncType = ftn.isDir
}
