package tools

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// 文件对比
func CompareDirs(pathl, pathr string) bool {
	result := true
	lst, _ := os.Stat(pathl)
	rst, _ := os.Stat(pathr)
	//文件不相同情况
	if lst == nil || rst == nil || lst.IsDir() != rst.IsDir() || (!lst.IsDir() && rst.Size() != lst.Size()) {
		fmt.Printf("对比失败 %s %s\n", lst, rst)
		fmt.Printf("路径分别为 %s %s\n", pathl, pathr)
		return false
	} else {
		fmt.Printf("%s %s 对比成功\n", pathl, pathr)
	}
	if lst != nil && rst != nil && lst.IsDir() && rst.IsDir() {
		next := map[string]bool{}
		files, _ := ioutil.ReadDir(pathl)
		for _, file := range files {
			next[file.Name()] = true
		}
		files, _ = ioutil.ReadDir(pathr)
		for _, file := range files {
			next[file.Name()] = true
		}
		for name := range next {
			if !CompareDirs(filepath.Join(pathl, name), filepath.Join(pathr, name)) {
				result = false
			}
		}
	}
	return result
}

// 文件同步
func SyncFile(pathSorce, pathTarget string, addDic, delDic map[string]bool) (map[string]bool, error) {
	fmt.Println("同步开始")
	log.Println("同步开始")
	synced := map[string]bool{}
	// 文件及文件夹删除
	for file := range delDic {
		if stat, err := os.Stat(filepath.Join(pathTarget, file)); err == nil && stat != nil {
			err := os.RemoveAll(filepath.Join(pathTarget, file))
			if err != nil {
				log.Fatal("删除文件失败", file)
				return synced, err
			}
			synced[file] = true
		}
	}
	//文件夹创建
	for dir, isDir := range addDic {
		//只对文件夹进行操作
		if isDir {
			err := os.MkdirAll(filepath.Join(pathTarget, dir), os.ModeDir)
			if err != nil {
				log.Fatal("文件夹创建出错", err)
				return synced, err
			}
			log.Println("新建文件夹:\t", dir)
			synced[dir] = true
		}
	}
	//文件复制
	for file, isDir := range addDic {
		//只对文件进行操作
		if !isDir {
			fi, err := os.Stat(filepath.Join(pathSorce, file))
			if err != nil {
				log.Fatal("获取文件信息错误", file)
				return synced, err
			}
			cont, err := ioutil.ReadFile(filepath.Join(pathSorce, file))
			if err != nil {
				log.Fatal("读取错误", file)
				return synced, err
			}
			err = ioutil.WriteFile(filepath.Join(pathTarget, file), cont, fi.Mode())
			if err != nil {
				log.Fatal("写入文件错误", file)
				return synced, err
			}
			log.Println("写入文件成功", file)
			synced[file] = true
		}
	}
	return synced, nil
}
