package tools

import (
	"BRGS/pkg/e"
	"BRGS/pkg/utils"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// 文件同步
func SyncFile(pathSorce, pathTarget string, addDic, delDic map[string]bool) (synced map[string]bool, err error) {
	log.Println("同步开始")
	// 文件及文件夹删除
	for file := range delDic {
		if stat, err := os.Stat(filepath.Join(pathTarget, file)); err == nil && stat != nil {
			err := os.RemoveAll(filepath.Join(pathTarget, file))
			if err != nil {
				err = e.TranslateToError(e.ERROR_DELETE, file, err.Error())
				goto errorLog
			}
			synced[file] = true
		} else if stat == nil {
			synced[file] = true
		}
	}
	//文件夹创建
	for dir, isDir := range addDic {
		//只对文件夹进行操作
		if isDir {
			err := os.MkdirAll(filepath.Join(pathTarget, dir), os.ModeDir)
			if err != nil {
				err = e.TranslateToError(e.ERROR_DELETE, dir, err.Error())
				goto errorLog
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
				err = e.TranslateToError(e.ERROR_READ, file, err.Error())
				goto errorLog
			}
			cont, err := ioutil.ReadFile(filepath.Join(pathSorce, file))
			if err != nil {
				err = e.TranslateToError(e.ERROR_READ, file, err.Error())
				goto errorLog
			}
			err = ioutil.WriteFile(filepath.Join(pathTarget, file), cont, fi.Mode())
			if err != nil {
				err = e.TranslateToError(e.ERROR_WRITE, file, err.Error())
				goto errorLog
			}
			log.Println("写入文件成功", file)
			synced[file] = true
		}
	}
	return synced, nil

errorLog:
	log.Println(err)
	return synced, err
}

// 遍历目录获取相对路径
func WalkDir(root string) map[string]string {
	result := map[string]string{}
	var walkDir func(string, string)
	walkDir = func(root string, prefix string) {
		fl, err := ioutil.ReadDir(root)
		if err != nil {
			panic(err)
		} else {
			for _, file := range fl {
				if file.IsDir() {
					walkDir(filepath.Join(root, file.Name()), prefix+"/"+file.Name())
				} else {
					result[filepath.Join(root, file.Name())] = (prefix + "/" + file.Name())[1:]
				}
			}
		}
	}
	walkDir(root, "")
	for key, value := range result {
		log.Printf("%-70s\t--------->\t%-50s", key, value)
	}
	return result
}

// 文件对比
func CompareDirs(pathl, pathr string) bool {
	result := true
	lst, _ := os.Stat(pathl)
	rst, _ := os.Stat(pathr)
	//文件不相同情况
	if !utils.CompareFile(pathl, pathr) {
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
