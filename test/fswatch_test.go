package test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
)

func TestWatch2(f *testing.T) {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add("D:\\testDir\\input")
	if err != nil {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}

func TestWatch(t *testing.T) {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	fmt.Println(watcher)
	watcher.Add("D:\\testDir\\input")
	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	// Add a path.
	rootPath := "D:\\testDir\\input"
	fl, err := os.OpenFile(filepath.Join("D:\\testDir", "fff.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer fl.Close()
	if err == nil {
		log.SetOutput(fl)
	}
	err = watcher.Add(rootPath)
	testContent(rootPath)
	if err != nil {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}

func testContent(rootPath string) {
	//监控内文件夹操作
	// if _, err := os.Stat(rootPath); err == nil {
	// 	os.RemoveAll(rootPath)

	// 	return
	// }
	// os.Mkdir(rootPath, os.ModeDir)
	fmt.Print("alter")
	if err := os.Mkdir(filepath.Join(rootPath, "dir1"), os.ModeDir); err != nil {
		fmt.Println("新建文件夹错误1")
		fmt.Println(err)
	}
	if fl, err := os.Create(filepath.Join(rootPath, "test.txt")); err == nil {
		defer fl.Close()
		fl.Write([]byte("fffffffffff"))
		time.Sleep(1e9 * 3)
		fl.Write([]byte("ggggggggggggg"))
	}
	if err := os.Mkdir(filepath.Join(rootPath, "dir2"), os.ModeDir); err != nil {
		fmt.Println("新建文件夹错误2")
		fmt.Println(err)
	}
	if fl, err := os.Open(filepath.Join(rootPath, "test.txt")); err == nil {
		defer fl.Close()
		fl.Write([]byte("fffffffffff"))
		time.Sleep(1e9 * 3)
		fl.Write([]byte("ggggggggggggg"))
	}
	//监控内子文件夹操作
	if err := os.Mkdir(filepath.Join(rootPath, "dir1", "dir3"), os.ModeDir); err != nil {
		fmt.Println("新建文件夹错误3")
		fmt.Println(err)
	}
	if fl, err := os.Create(filepath.Join(rootPath, "dir1", "test.txt")); err == nil {
		defer fl.Close()
		fl.Write([]byte("fffffffffff"))
		time.Sleep(1e9 * 3)
		fl.Write([]byte("ggggggggggggg"))
	}
}

func prepareTestDirTree(tree string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		return "", fmt.Errorf("error creating temp directory: %v\n", err)
	}

	err = os.MkdirAll(filepath.Join(tmpDir, tree), 0755)
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", err
	}

	return tmpDir, nil
}

// func TestWalk() {
// 	tmpDir, err := prepareTestDirTree("./dir/to/walk/skip")
// 	if err != nil {
// 		fmt.Printf("unable to create test dir tree: %v\n", err)
// 		return
// 	}
// 	defer os.RemoveAll(tmpDir)
// 	os.Chdir(tmpDir)

// 	subDirToSkip := "skip"

// 	fmt.Println("On Unix:")
// 	err = filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
// 		if err != nil {
// 			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
// 			return err
// 		}
// 		if info.IsDir() && info.Name() == subDirToSkip {
// 			fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
// 			return filepath.SkipDir
// 		}
// 		fmt.Printf("visited file or dir: %q\n", path)
// 		return nil
// 	})
// 	if err != nil {
// 		fmt.Printf("error walking the path %q: %v\n", tmpDir, err)
// 		return
// 	}
// }
