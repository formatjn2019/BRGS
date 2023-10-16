package manage

import (
	"BRGS/management"
	"BRGS/pkg/tools"
	"BRGS/pkg/utils"
	"BRGS/routers"
	"fmt"
	"log"
	"path"
	"strconv"
	"testing"
	"time"
)

const inputDir = "D:\\testDir\\input"
const tempDir = "D:\\testDir\\temp"
const outputDir = "D:\\testDir\\output"

func getMonitor() *management.Monitor {
	ba := management.BackupArchive{
		Name:            "test",
		WatchDir:        inputDir,
		TempDir:         tempDir,
		ArchiveDir:      outputDir,
		SyncInterval:    10,
		ArchiveInterval: 20,
	}
	monitor := management.CreateMonitor(ba)
	return monitor
}

func TestMonitor(t *testing.T) {
	monitor := getMonitor()
	monitor.Run()
	go func() {
		time.Sleep(20 * time.Second)
		println(monitor.Pause())
		time.Sleep(30 * time.Second)
		println(monitor.Continue())
	}()
	for i := 0; i < 100; i++ {
		time.Sleep(9 * time.Second)
		utils.OsOperation("cf", path.Join(inputDir, strconv.Itoa(i)+".txt"))
	}
	<-make(chan struct{})
}

func TestMonitorWeb(t *testing.T) {
	management.MonServer = getMonitor()
	management.MonServer.Run()
	go func() {
		err := routers.StartServer()
		if err != nil {
			log.Fatal(err)
		}
	}()
	<-make(chan struct{})
}

func TestTg(t *testing.T) {
	fmt.Println(tools.GenerateRule("test"))
}
