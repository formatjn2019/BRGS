package basic

import (
	"fmt"
	"io/fs"
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
