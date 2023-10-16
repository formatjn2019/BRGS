package basic

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestFilepath(t *testing.T) {
	skipList := make([]string, 0)
	filepath.Walk("D:\\testDir", func(path string, info fs.FileInfo, err error) error {
		abs, err := filepath.Abs(path)
		base := filepath.Base(abs)
		if matched, e := filepath.Match("README.md", base); e == nil && matched {
			skipList = append(skipList, abs)
		}
		return nil
	})
	fmt.Println(len(skipList))
	for _, path := range skipList {
		println(path)
	}
}

func TestGetPath(t *testing.T) {
	filepath.Base("D:\\testDir")
	err := os.MkdirAll("D:\\testDir", os.ModePerm)
	fmt.Println(err)
	fmt.Println(os.Stat("d:/22g"))
	println(filepath.Rel("D:\\testDir", "D:\\testDir\\input.text"))
	println(filepath.Dir("D:\\testDir\\input.text"))
	println(filepath.Ext("D:\\testDir\\input.text"))
}
