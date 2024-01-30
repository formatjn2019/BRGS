package main

import (
	"BRGS/conf"
	"BRGS/management"
	"BRGS/management/menu"
	"BRGS/routers"
	"embed"
	_ "embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	serverFlag bool
	manualFlag bool
	rule       *management.BackupArchive
)

// 初始化模式
func init() {
	rule = &management.BackupArchive{}
	flag.BoolVar(&serverFlag, "s", false, "web服务")
	flag.BoolVar(&manualFlag, "m", false, "手动操作")
	flag.StringVar(&rule.Name, "n", "", "名称")
	flag.StringVar(&rule.WatchDir, "wd", "", "监控文件夹")
	flag.StringVar(&rule.TempDir, "td", "", "中转文件夹")
	flag.StringVar(&rule.ArchiveDir, "ad", "", "存储文件夹")
	flag.IntVar(&rule.ArchiveInterval, "ai", 0, "存档间隔")
	flag.IntVar(&rule.SyncInterval, "si", 0, "同步间隔")
}

func showFlags() {
	fmt.Println("serverFlag: ", serverFlag)
	fmt.Println("manualFlag: ", manualFlag)
	fmt.Println("rule: ", rule)
}

//go:embed web
var app embed.FS

func starWebServer() error {
	dist, err := fs.Sub(app, "web/dist")
	if err != nil {
		return err
	}
	http.Handle("/", http.FileServer(http.FS(dist)))
	err = http.ListenAndServe(":"+conf.ServerConf.WebPort, nil)
	return err
}

func main() {
	var wg sync.WaitGroup
	flag.Parse()
	showFlags()
	var err error

	management.MonServer = management.CreateMonitor(*rule)

	log.Println(err)

	// 不开启手动和自动模式，则启动命令行
	if !(serverFlag || manualFlag) {
		menu.ConfigMenu()
	} else {
		// 解析配置
		if errors := rule.CheckConfig(); len(errors) > 0 {
			fmt.Println(errors)
			panic(errors)
		}
		// 开启web服务
		if serverFlag {
			// 前端
			go func() {
				wg.Add(1)
				err = starWebServer()
				if err != nil {
					log.Fatal(err)
				}
				wg.Done()
			}()
			// 后端
			go func() {
				wg.Add(1)
				err = routers.StartServer()
				if err != nil {
					log.Fatal(err)
				}
				wg.Done()
			}()
		}
		// 开启手动控制
		if manualFlag {
			wg.Add(1)
			menu.ControlMenu()
			wg.Done()
		}
	}
	// 等待其它线程启动
	time.Sleep(10 * time.Second)
	wg.Wait()
}
