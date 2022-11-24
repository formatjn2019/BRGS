package models

import (
	"BRGS/pkg/util"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/fsnotify/fsnotify"
)

func CreateFsTreeRoot(inputdir, outputdir string) *FSTreeRoot {
	ftn, err := InitScanFolder(inputdir)
	if err != nil {
		log.Println("读取文件树错误")
		log.Println("error: ", err)
		return nil
	}
	err = ftn.CreateTargetDir(outputdir)
	if err != nil {
		log.Println("创建输出文件夹错误")
		log.Println("error: ", err)
		return nil
	}

	ftn.WatchDirs()
	ftn.Watch()
	b := ftn.BackupFiles()
	fmt.Println(b)
	return ftn
}

// 树的根
type FSTreeRoot struct {
	appendDics chan map[string]bool
	deleteDics chan map[string]bool
	dic        map[string]*FsTreeNode
	endTag     chan struct{}
	eventCount int64
	events     chan fsnotify.Event
	source     string
	state      *FstreeType
	syncedDics chan map[string]bool
	syncOp     chan int
	syncTag    chan struct{}
	target     string
	watcher    *fsnotify.Watcher
}

// 根据输入文件夹创建根节点
func InitScanFolder(rootPath string) (*FSTreeRoot, error) {
	info, err := os.Stat(rootPath)
	//不是文件夹或路径错误
	if err != nil || !info.IsDir() {
		log.Println("扫描路径出错")
		log.Printf("info: %v\n", info)
		log.Printf("err: %v\n", err)
		return nil, err
	} else {
		rootPath, _ = filepath.Abs(rootPath)
	}
	treeMap := map[string]*FsTreeNode{}
	result := &FSTreeRoot{
		appendDics: make(chan map[string]bool),
		deleteDics: make(chan map[string]bool),
		dic:        treeMap,
		endTag:     make(chan struct{}),
		eventCount: 0,
		events:     make(chan fsnotify.Event, 10000),
		source:     rootPath,
		state:      &FstreeType{state: FSTREE_WATCH},
		syncedDics: make(chan map[string]bool),
		syncOp:     make(chan int),
		syncTag:    make(chan struct{}),
	}
	treeMap[""] = &FsTreeNode{path: "", exist: true, isDir: true, isAlter: false, synced: true, syncType: true, subs: map[string]*FsTreeNode{}}
	result.ScanFolder(rootPath, "")
	return result, nil
}

// 正向同步文件
func (root *FSTreeRoot) BackupFiles() bool {
	root.syncOp <- FSTREE_OP_BACKUP_PREPERE
	appendDic := <-root.appendDics
	deleteDic := <-root.deleteDics
	log.Printf("appendDic: %v\n", appendDic)
	log.Printf("deleteDic: %v\n", deleteDic)
	//同步完成
	synced, err := util.SyncFile(root.source, root.target, appendDic, deleteDic)
	root.syncOp <- FSTREE_OP_BACKUP_END
	root.syncedDics <- synced
	return err == nil
}

// 创建文件节点
func (root *FSTreeRoot) CreateNode(path, parent string, isDir bool) *FsTreeNode {
	if node, ok := root.dic[path]; !ok {
		newNode := &FsTreeNode{path: path, parent: root.dic[parent], exist: true, isDir: isDir, isAlter: true, synced: false, syncType: isDir, subs: map[string]*FsTreeNode{}}
		root.dic[path] = newNode
		newNode.parent.subs[path] = newNode
		return newNode
	} else {
		node.exist = true
		node.isAlter = true
		node.isDir = isDir
		node.parent.subs[path] = node
		if node.isDir {
			node.subs = map[string]*FsTreeNode{}
		}
		return node
	}
}

// 设置目标文件夹
func (root *FSTreeRoot) CreateTargetDir(target string) error {
	root.target = target
	if _, err := os.Stat(target); err == nil {
		err := os.RemoveAll(target)
		if err != nil {
			return err
		}
	}
	return os.MkdirAll(target, fs.ModeDir)
}

// 获取节点
func (root *FSTreeRoot) GetNode(key string) *FsTreeNode {
	return root.dic[key]
}

// 逆向同步文件
func (root *FSTreeRoot) RecoverFiles() bool {
	root.syncOp <- FSTREE_OP_RECOVER_PREPERE
	appendDic := <-root.appendDics
	deleteDic := <-root.deleteDics
	log.Printf("appendDic: %v\n", appendDic)
	log.Printf("deleteDic: %v\n", deleteDic)
	//同步完成
	synced, err := util.SyncFile(root.target, root.source, appendDic, deleteDic)
	root.syncOp <- FSTREE_OP_RECOVER_END
	root.syncedDics <- synced
	return err == nil
}

// 扫描文件
func (root *FSTreeRoot) ScanFolder(path, parent string) {
	if files, err := ioutil.ReadDir(path); err == nil {
		for _, info := range files {
			key := filepath.Join(parent, info.Name())
			path := filepath.Join(path, info.Name())
			root.CreateNode(key, parent, info.IsDir())
			if info.IsDir() {
				root.ScanFolder(path, key)
			}
		}
	} else {
		log.Printf("扫描文件夹%s失败:\t%v\n", path, err)
	}
}

// 展示文件节点
func (root *FSTreeRoot) Show(changed bool) {
	fmt.Println(strings.Repeat("*", 200))
	for k, v := range root.dic {
		if changed {
			if v.isAlter {
				log.Println(k, v)
			}
		} else {
			log.Println(k, v)
		}

	}
	fmt.Println(strings.Repeat("*", 200))
}

// 控制文件树的操作

func (root *FSTreeRoot) watchTree() {
	for {
		select {
		//取得操作映射
		case _, ok := <-root.syncTag:
			for ; ok && atomic.LoadInt64(&root.eventCount) > 0; atomic.AddInt64(&root.eventCount, -1) {
				event, ok := <-root.events
				log.Println("event:", event)
				path := event.Name[len(root.source)+1:]
				node, ok := root.dic[path]
				if event.Op == fsnotify.Remove || event.Op == fsnotify.Rename {
					if ok {
						fmt.Println("mark", event, path, node)
						node.Delete()
						if node.isDir {
							err := root.watcher.Remove(event.Name)
							if err != nil {
								fmt.Println("移除监控失败", err)
								fmt.Println(node)
								fmt.Println(len(node.subs))
							}

						}
						fmt.Println("删除测试成功")
					}
				} else {
					//字典中不存在则创建节点
					if (event.Op == fsnotify.Write && node == nil) || event.Op == fsnotify.Create {
						info, err := os.Stat(event.Name)
						if err != nil {
							log.Println(path, err)
						}
						fmt.Println(path)
						parent := event.Name[len(root.source):strings.LastIndex(event.Name, "\\")]
						if len(parent) > 0 {
							parent = parent[1:]
						}
						log.Println(path, parent, info.IsDir())
						root.CreateNode(path, parent, info.IsDir())
						if info.IsDir() {
							root.ScanFolder(event.Name, path)
							root.WatchDirs()
						}
					} else if event.Op == fsnotify.Write && ok {
						//状态校验
						if _, err := os.Stat(event.Name); err == nil {
							node.Update()
						} else {
							node.Delete()
						}

					} else {
						log.Println("其它状态")
						log.Println(event)
						log.Println(event.Op)
					}
				}
			}
			//取得操作映射
		case op, ok := <-root.syncOp:
			log.Println(op, ok)
			if !ok {
				return
			}
			switch op {
			case FSTREE_OP_BACKUP_PREPERE:
				root.getSyncFiles()
			case FSTREE_OP_BACKUP_END:
				for key := range <-root.syncedDics {
					root.dic[key].Sync()
				}
				root.cleanup()
			case FSTREE_OP_RECOVER_PREPERE:
				root.getReverseSyncFiles()
			case FSTREE_OP_RECOVER_END:
				for key := range <-root.syncedDics {
					root.dic[key].ReverseSync()
				}
				root.cleanup()
			default:
				panic("unsupport operation")
			}

		}
	}
}

// 文件监控操作
func (root *FSTreeRoot) Watch() {
	//负责fsnotify的事件监听
	go func() {
		for {
			select {
			// 事件转移
			case event, ok := <-root.watcher.Events:
				if !ok {
					return
				}
				log.Println(event, "监听文件改变")
				root.events <- event
				atomic.AddInt64(&root.eventCount, 1)
				log.Println(event, root.state.state)
				if root.state.state == FSTREE_WATCH {
					root.syncTag <- struct{}{}
				}
			case err, ok := <-root.watcher.Errors:

				if ok { // Channel was closed (i.e. Watcher.Close() was called).
					if !ok { // Channel was closed (i.e. Watcher.Close() was called).
						log.Println("监控异常", err)
					}
				}
			}
		}
	}()
	go root.watchTree()
}

// 添加文件夹监控
func (root *FSTreeRoot) WatchDirs(exceptRule ...string) error {
	if root.watcher == nil {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return err
		}
		root.watcher = watcher
	}
	for k, v := range root.dic {
		if v.isDir && v.exist {
			fullPath := filepath.Join(root.source, k)
			log.Println("监控目录", k, fullPath)
			root.watcher.Remove(fullPath)
			root.watcher.Add(fullPath)
		}
	}
	return nil
}

// 清理无效节点(新建后未同步就已经删除的节点)
func (root *FSTreeRoot) cleanup() {
	nodes := make([]string, 0)
	for key, node := range root.dic {
		//（未/已）同步，已删除
		if !node.synced && !node.exist {
			nodes = append(nodes, key)
		}
		// 已经存在的同名文件夹，取消更改标记
		if node.isAlter && node.exist && node.synced && node.isDir && node.syncType {
			node.isAlter = false
		}
	}
	for _, key := range nodes {
		delete(root.dic, key)
	}
}

// 获取逆向同步的文件夹
func (root *FSTreeRoot) getReverseSyncFiles() {
	appendDic, deleteDic := map[string]bool{}, map[string]bool{}
	//清理无效节点
	root.cleanup()
	//标记新增和删除节点
	for path, node := range root.dic {
		if node.isAlter {
			log.Printf("标记同步节点%40s\t%s\n", path, node)
			if node.synced {
				appendDic[path] = node.syncType
				// 文件夹删除又创建后或者 删除文件夹又创建了同名文件 清理目标文件夹指定位置
				if node.exist && (node.syncType || node.syncType != node.isDir) {
					deleteDic[path] = node.isDir
				}
			} else {
				deleteDic[path] = node.isDir
			}
		}
	}
	root.appendDics <- appendDic
	root.deleteDics <- deleteDic
}

// 获取同步的文件夹
func (root *FSTreeRoot) getSyncFiles() {
	appendDic, deleteDic := map[string]bool{}, map[string]bool{}
	//清理无效节点
	root.cleanup()
	//标记新增和删除节点
	for path, node := range root.dic {
		if node.isAlter {
			log.Printf("标记同步节点%40s\t%s\n", path, node)
			if node.exist {
				appendDic[path] = node.isDir
				// 文件夹删除又创建后或者 删除文件夹又创建了同名文件 清理目标文件夹指定位置
				if node.synced && (node.isDir || node.syncType != node.isDir) {
					deleteDic[path] = node.isDir
				}
			} else {
				deleteDic[path] = node.isDir
			}
		}
	}
	root.appendDics <- appendDic
	root.deleteDics <- deleteDic
}
