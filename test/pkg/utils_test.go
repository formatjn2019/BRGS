package test

import (
	"BRGS/management"
	"BRGS/pkg/util"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-ini/ini"
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

func TestIni(t *testing.T) {
	cfg, err := ini.Load("../test.ini")
	println()
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	// // 典型读取操作，默认分区可以使用空字符串表示
	// fmt.Println("App Mode:", cfg.Section("").Key("app_mode").String())
	// fmt.Println("Data Path:", cfg.Section("paths").Key("data").String())

	// // 我们可以做一些候选值限制的操作
	// fmt.Println("Server Protocol:",
	// 	cfg.Section("server").Key("protocol").In("http", []string{"http", "https"}))
	// // 如果读取的值不在候选列表内，则会回退使用提供的默认值
	// fmt.Println("Email Protocol:",
	// 	cfg.Section("server").Key("protocol").In("smtp", []string{"imap", "smtp"}))

	// // 试一试自动类型转换
	// fmt.Printf("Port Number: (%[1]T) %[1]d\n", cfg.Section("server").Key("http_port").MustInt(9999))
	// fmt.Printf("Enforce Domain: (%[1]T) %[1]v\n", cfg.Section("server").Key("enforce_domain").MustBool(false))

	// // 差不多了，修改某个值然后进行保存
	// cfg.Section("").Key("app_mode").SetValue("production")
	fmt.Println(cfg.Section("excel_default_tip").Key("name").String())
	fmt.Println(cfg.Section("excel_default_tip").Key("watchDir").String())
	fmt.Println(cfg.Section("excel_default_tip").Key("tempDir").String())
	fmt.Println(cfg.Section("excel_default_tip").Key("archiveDir").String())
	fmt.Println(cfg.Section("excel_default_tip").Key("syncInterval").String())
	fmt.Println(cfg.Section("excel_default_tip").Key("archiveInterval").String())
	// cfg.SaveTo("my.ini.local")
}
