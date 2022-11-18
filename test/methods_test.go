package test

import (
	"BRGS/management"
	"BRGS/pkg/util"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestCommand(t *testing.T) {
	management.CommandMenu(false, "菜单1", "菜单2")
}

func TestWalk(t *testing.T) {
	util.WalkDir("D:\\testDir")
}

func TestWrite(t *testing.T) {
	context := []map[string]string{{"标题1": "332", "标题2": "4422"}, {"标题1": "33", "标题2": "44"}}
	util.WriteCsvWithDict("config.csv", context)
}

func TestReadConfig(t *testing.T) {
	// result, err := tools.ReadCsvAsDict("config.csv")
	// fmt.Println(err)
	// for k, v := range result {
	// 	fmt.Println(k, v)
	// }
	command := management.ReadConfigCommand{}
	util.ReadCsvAsDict("config.csv")
	command.Execute()
}

func TestSaveConfig(t *testing.T) {
	temp := os.TempDir()
	fmt.Println(temp)

	time := time.Now()
	fmt.Println(time.GoString())
	fmt.Println(time.Format("2006-01-02 15:04:05"))
	fmt.Println(time.Format("20060102_150405"))
}

func TestZip(t *testing.T) {
	start := time.Now().UnixNano()
	dict := util.WalkDir("D:\\testDir\\input")
	fmt.Printf("tools.WriteZip(D:\\testDir\\archive\\test.zip, dict): %v\n", util.WriteZip("D:\\testDir\\archive\\test.zip", dict))
	end := time.Now().UnixNano()
	fmt.Println(start)
	fmt.Println(end)
	fmt.Println(end - start)
}

func TestRecover(t *testing.T) {
	// start := time.Now().UnixNano()
	// tools.RecoverFromArchive("D:\\testDir\\archive\\test.zip", "D:\\testDir\\output")
	// end := time.Now().UnixNano()
	// fmt.Println(start)
	// fmt.Println(end)
	// fmt.Println(end - start)
	// rule := regexp.MustCompile("test1" + "_20\\d{}6_\\d{6}.zip")
	// rule := regexp.MustCompile("test1" + )
	// // rule := regexp.MustCompile("test1" + ".*?zip")
	// archDir := "D:\\testDir\\archive"
	// if files, err := ioutil.ReadDir(archDir); err == nil {
	// 	fmt.Println(files)
	// 	fmt.Printf("%T\n", files)
	// 	for _, info := range files {
	// 		fmt.Println(info)
	// 		fmt.Println(info.Name())
	// 		if rule.MatchString(info.Name()) {
	// 			fmt.Println(info.Name(), "匹配成功")
	// 		}
	// 		// key := filepath.Join(parent, info.Name())
	// 		// path := filepath.Join(path, info.Name())
	// 		// root.CreateNode(key, parent, info.IsDir())
	// 		// if info.IsDir() {
	// 		// 	root.ScanFolder(path, key)
	// 		// }
	// 	}
	// } else {
	// 	log.Printf("扫描压缩文件夹%s失败:\t%v\n", archDir, err)
	// }
}
