package test

import (
	"BRGS/tools"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
)

const INPUT_DIR, OUTPUT_DIR, TEMP_DIR = "D:\\testDir\\input", "D:\\testDir\\output", "D:\\testDir\\temp"

//文件树监控测试
func TestFsTree(t *testing.T) {
	fmt.Println(strings.Repeat("+", 200))
	createFsTreeRoot()
	<-make(chan struct{})
}

//文件树方法测试
func TestFsTreeBackup(t *testing.T) {
	fmt.Println(strings.Repeat("*", 200))
	ftr := createFsTreeRoot()
	syncTestBackup(ftr, []string{"cf", filepath.Join(INPUT_DIR, "test.txt")}, []string{"cd", filepath.Join(INPUT_DIR, "dir1")})
	syncTestBackup(ftr, []string{"rm", filepath.Join(INPUT_DIR, "test.txt")}, []string{"cd", filepath.Join(INPUT_DIR, "dir1")}, []string{"rm", filepath.Join(INPUT_DIR, "dir1")})
}

//文件树方法测试
func TestFsTreeRecoverBackup(t *testing.T) {
	fmt.Println(strings.Repeat("*", 200))
	ftr := createFsTreeRoot()
	syncTestRecover(ftr, []string{"cf", filepath.Join(INPUT_DIR, "test.txt")}, []string{"cd", filepath.Join(INPUT_DIR, "dir1")})
	syncTestRecover(ftr, []string{"rm", filepath.Join(INPUT_DIR, "test.txt")}, []string{"cd", filepath.Join(INPUT_DIR, "dir1")}, []string{"rm", filepath.Join(INPUT_DIR, "dir1")})
}

//文件树方法测试
func TestFsTreeRandomBackup(t *testing.T) {
	f, _ := os.OpenFile("log.log", os.O_CREATE|os.O_WRONLY, 0644)
	log.SetOutput(f)
	fmt.Print(2)
	ftr := createFsTreeRoot()
	showTrees(ftr)
	for rol := 0; rol < 1; rol++ {
		for i := 0; i < 20; i++ {
			randomOperation(INPUT_DIR)
			ftr.Show(true)
		}
		fmt.Printf("tools.CompareDirs(INPUT_DIR, OUTPUT_DIR): %v\n", tools.CompareDirs(INPUT_DIR, OUTPUT_DIR))
		if 0 > 0 {
			fmt.Printf("ftr.BackupFiles(): %v\n", ftr.BackupFiles())
		} else {
			fmt.Printf("ftr.RecoverFiles(): %v\n", ftr.RecoverFiles())
		}
		fmt.Printf("tools.CompareDirs(INPUT_DIR, OUTPUT_DIR): %v\n", tools.CompareDirs(INPUT_DIR, OUTPUT_DIR))
	}
}
func TestOs(t *testing.T) {
	// err := os.RemoveAll("D:\\testDir\\input\\dir2")
	// fmt.Println(err)
	fmt.Println(time.Now())
	os.RemoveAll(TEMP_DIR)
	os.MkdirAll(TEMP_DIR, os.ModeDir)
	OsOperation("cf", filepath.Join(TEMP_DIR, "test.txt"))
	OsOperation("up", filepath.Join(TEMP_DIR, "test.txt"))
	OsOperation("md", filepath.Join(TEMP_DIR, "tdr"))
	OsOperation("rm", filepath.Join(TEMP_DIR, "tdr"))
	OsOperation("rm", filepath.Join(TEMP_DIR, "test.txt"))
	fmt.Println(time.Now())
	fmt.Println(time.Now())
}

func TestOtherMethods(t *testing.T) {
	path := "D:\\testDir\\input\\tdr"
	fmt.Println(strings.Index(path, "D"))
	println("232")
	fmt.Printf("tools.CompareDirs(INPUT_DIR, OUTPUT_DIR): %v\n", tools.CompareDirs(INPUT_DIR, OUTPUT_DIR))
}

func TestFsNotify(t *testing.T) {
	fmt.Println(strings.Repeat("*", 200))
	OsOperation("rm", TEMP_DIR)
	OsOperation("md", TEMP_DIR)
	rootPath, _ := filepath.Abs(TEMP_DIR)
	watcher, _ := fsnotify.NewWatcher()
	watcher.Add(rootPath)
	go func() {
		for event := range watcher.Events {
			log.Println("event:", event)
		}
	}()
	OsOperation("md", filepath.Join(TEMP_DIR, "tp2"))
	watcher.Add(filepath.Join(TEMP_DIR, "tp2"))
	OsOperation("mv", filepath.Join(TEMP_DIR, "tp2"), filepath.Join(TEMP_DIR, "tp3"))
	// watcher.Add("D:\\testDir\\temp\\tp3")
	err := watcher.Remove(filepath.Join(TEMP_DIR, "tp3"))
	fmt.Println(err)
	err = watcher.Add(filepath.Join(TEMP_DIR, "tp3"))
	fmt.Println(err)
	OsOperation("md", filepath.Join(TEMP_DIR, "ttw"))
	fmt.Println(rootPath)
	fmt.Println(rootPath)
}

func OsOperation(operation string, path ...string) {
	fmt.Println(operation, path)
	switch operation {
	//create dir
	case "md":
		os.MkdirAll(path[0], os.ModeDir)
		//create file
	case "cf":
		ioutil.WriteFile(path[0], []byte("create file"), 0644)
		//delete
	case "rm":
		//update
		os.RemoveAll(path[0])
	case "up":
		origin, _ := ioutil.ReadFile(path[0])
		p := make([]byte, 100)
		rand.Read(p)
		origin = append(origin, p[rand.Intn(90):]...)
		ioutil.WriteFile(path[0], origin, 0644)
	case "mv":
		os.Rename(path[0], path[1])
	}
	time.Sleep(1e9)
}

func createFsTreeRoot() *tools.FSTreeRoot {
	ftn, err := tools.InitScanFolder(INPUT_DIR)
	if err != nil {
		log.Println("读取文件树错误")
		log.Println("error: ", err)
		return nil
	}
	err = ftn.CreateTargetDir(OUTPUT_DIR)
	if err != nil {
		log.Println("创建输出文件夹错误")
		log.Println("error: ", err)
		return nil
	}
	b := ftn.BackupFiles()
	ftn.WatchDirs()
	ftn.Watch()
	fmt.Println(b)
	return ftn
}

func randomOperation(inputDir string) {
	osOp := []string{"md", "cf", "rm", "up", "mv"}
	getRandName := func() string {
		var sb strings.Builder
		sb.WriteRune('a' + rune(rand.Intn(26)))
		for rand.Intn(2) > 0 {
			if rand.Intn(2) > 0 {
				sb.WriteRune('a' + rune(rand.Intn(26)))
			} else {
				sb.WriteRune('A' + rune(rand.Intn(26)))
			}
		}
		return sb.String()
	}
	var getRandPath func(string, bool) string
	getRandPath = func(parentPath string, conRoot bool) string {
		entrys, _ := ioutil.ReadDir(parentPath)
		var index int
		if !conRoot {
			index = rand.Intn(len(entrys))
		} else {
			index = rand.Intn(len(entrys) + 1)
		}
		if index == len(entrys) {
			return parentPath
		} else {
			return getRandPath(filepath.Join(parentPath, entrys[index].Name()), true)
		}
	}

	getRandDirOrFile := func(parentPath string, isDir bool) string {
		for i := 0; i < 100000; i++ {
			tempPath := getRandPath(parentPath, true)
			info, _ := os.Stat(tempPath)
			if info.IsDir() == isDir {
				return tempPath
			}
		}
		panic("找不到文件")
	}
	switch operation := osOp[rand.Intn(len(osOp))]; operation {
	case "md", "cf":
		dirPath := getRandDirOrFile(inputDir, true)
		name := getRandName()
		for info, _ := os.Stat(filepath.Join(dirPath, name)); info != nil; {
			name = getRandName()
			info, _ = os.Stat(filepath.Join(dirPath, name))
		}
		OsOperation(operation, filepath.Join(dirPath, name))
	case "up":
		filePath := getRandDirOrFile(inputDir, false)
		OsOperation(operation, filePath)
	case "rm":
		filePath := getRandPath(inputDir, false)
		OsOperation(operation, filePath)
	case "mv":
		filePath := getRandPath(inputDir, false)
		name := getRandName()
		OsOperation(operation, filePath, filePath+name)
	}
}

func showTrees(root *tools.FSTreeRoot) {
	rootNode := root.GetNode("")
	var showNode func(int, *tools.FsTreeNode)
	showNode = func(depth int, node *tools.FsTreeNode) {
		fmt.Printf("%d%100s%v", depth, strings.Repeat(" ", depth), node)
		for _, ftn := range node.Subs() {
			showNode(depth+1, ftn)
		}
	}
	showNode(0, rootNode)
}

func syncTestBackup(ftn *tools.FSTreeRoot, opers ...[]string) {
	println(strings.Repeat("+", 200))
	fmt.Printf("更改前文件夹同步状态%t\n", tools.CompareDirs(INPUT_DIR, OUTPUT_DIR))
	for _, oper := range opers {
		OsOperation(oper[0], oper[1:]...)
	}
	fmt.Printf("更改后文件夹同步状态%t\n", tools.CompareDirs(INPUT_DIR, OUTPUT_DIR))
	ftn.Show(false)
	ftn.BackupFiles()
	ftn.Show(false)
	fmt.Printf("同步后文件夹同步状态%t\n", tools.CompareDirs(INPUT_DIR, OUTPUT_DIR))
	println(strings.Repeat("-", 200))
}

func syncTestRecover(ftn *tools.FSTreeRoot, opers ...[]string) {
	println(strings.Repeat("+", 200))
	fmt.Printf("更改前文件夹同步状态%t\n", tools.CompareDirs(INPUT_DIR, OUTPUT_DIR))
	for _, oper := range opers {
		OsOperation(oper[0], oper[1:]...)
	}
	fmt.Printf("更改后文件夹同步状态%t\n", tools.CompareDirs(INPUT_DIR, OUTPUT_DIR))
	ftn.Show(false)
	ftn.RecoverFiles()
	ftn.Show(false)
	fmt.Printf("同步后文件夹同步状态%t\n", tools.CompareDirs(INPUT_DIR, OUTPUT_DIR))
	println(strings.Repeat("-", 200))
}
