package tools

import (
	"BRGS/pkg/e"
	"BRGS/pkg/utils"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
)

// 大于指定大小分片哈希
// 20M
const maxSilceSize = 20 * 1024 * 1024

// SyncFile 文件同步
func SyncFile(pathSorce, pathTarget string, addDic, delDic map[string]bool) (synced map[string]bool, err error) {
	synced = map[string]bool{}
	log.Println("同步开始")
	// 文件及文件夹删除
	for file := range delDic {
		if stat, err := os.Stat(filepath.Join(pathTarget, file)); err == nil && stat != nil {
			err := os.RemoveAll(filepath.Join(pathTarget, file))
			if err != nil {
				err = e.TranslateToError(e.ErrorDelete, file, err.Error())
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
				err = e.TranslateToError(e.ErrorDelete, dir, err.Error())
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
				err = e.TranslateToError(e.ErrorRead, file, err.Error())
				goto errorLog
			}
			cont, err := ioutil.ReadFile(filepath.Join(pathSorce, file))
			if err != nil {
				err = e.TranslateToError(e.ErrorRead, file, err.Error())
				goto errorLog
			}
			err = ioutil.WriteFile(filepath.Join(pathTarget, file), cont, fi.Mode())
			if err != nil {
				err = e.TranslateToError(e.ErrorWrite, file, err.Error())
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

// WalkDir 遍历目录获取相对路径
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

// CompareDirs 文件对比
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

// CalculateAllUid 计算路径下所以文件的uid
func CalculateAllUid(rootPath string) (map[string]string, error) {
	result := map[string]string{}
	// 运算操作通道
	calculateChan := make(chan chan string, 0)
	defer close(calculateChan)
	// 运算操作线程池
	go CalculateUidPool(calculateChan)
	var lock sync.Mutex
	var wg sync.WaitGroup
	err := filepath.Walk(rootPath, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			if err != nil {
				return err
			}
			c := make(chan string, 0)
			calculateChan <- c
			c <- path
			go func(path string) {
				wg.Add(1)
				defer wg.Done()
				lock.Lock()
				defer lock.Unlock()
				uid := <-c
				if uid == "" {
					err = e.TranslateToError(e.ErrorCalculate, "Fail to calculate", path)
					return
				}
				result[path] = uid
			}(path)
		}
		return err
	})
	wg.Wait()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CheckUid 核验uid
func CheckUid(pathUidDic map[string]string) (result map[string]string, err error) {
	result = map[string]string{}
	// 运算操作通道
	calculateChan := make(chan chan string, 0)
	defer close(calculateChan)
	// 运算操作线程池
	go CalculateUidPool(calculateChan)
	var lock sync.Mutex
	var wg sync.WaitGroup
	for path, uid := range pathUidDic {
		c := make(chan string, 0)
		calculateChan <- c
		c <- path
		go func(path, uid string) {
			wg.Add(1)
			defer wg.Done()
			lock.Lock()
			defer lock.Unlock()
			uidN := <-c
			if uid == "" {
				err = e.TranslateToError(e.ErrorCalculate, "Fail to calculate", path)
				return
			}
			if uidN != uid {
				result[path] = uid
			}
		}(path, uid)
	}
	wg.Wait()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func CalculateUidPool(calculateChan chan chan string) {
	var chanGroup = make(chan struct{}, runtime.NumCPU()*2)
	defer close(chanGroup)
	for conChan := range calculateChan {
		go func(conChan chan string) {
			chanGroup <- struct{}{}
			path := <-conChan
			uid, err := CalculateUid(path)
			if err != nil {
				log.Print(path, err)
				conChan <- ""
				return
			}
			conChan <- uid
			<-chanGroup
		}(conChan)
	}
}

// SyncFileWithHardLink 同步文件
// 非新文件时，采用硬链接形式记录
func SyncFileWithHardLink(uid, oldFilePath, newFilePath string, newInfo fs.FileInfo, linkedFilesDic map[string]map[string]fs.FileInfo, lock sync.Locker) (e error) {
	lock.Lock()
	defer lock.Unlock()
	//判断是否第一次出现该文件
	if linkedFiles, ok := linkedFilesDic[uid]; !ok {
		//第一次出现，记录
		linkedFilesDic[uid] = map[string]fs.FileInfo{newFilePath: newInfo}
		_, err := copyFile(oldFilePath, newFilePath)
		return err
	} else {
		//判断是重复文件的硬链接
		for _, oldInfo := range linkedFiles {
			if os.SameFile(newInfo, oldInfo) {
				return nil
			}
		}
		tempPath, idx := newFilePath+"_0", 1
		//生成中转文件名
		for _, err := os.Stat(tempPath); err == nil; idx++ {
			tempPath = tempPath[:len(tempPath)-len(strconv.Itoa(idx))] + strconv.Itoa(idx+1)
			_, err = os.Stat(tempPath)
		}
		//生成硬链接
		for linkedFilePath, _ := range linkedFiles {
			//链接成功，进行操作，返回
			if err := os.Link(linkedFilePath, tempPath); err == nil {
				println("relink", newFilePath)
				err = os.Remove(newFilePath)
				if err != nil {
					return err
				}
				return os.Rename(tempPath, newFilePath)
			}
		}
		//尝试生成硬链接失败,将该文件加入列表
		linkedFiles[newFilePath] = newInfo
		_, err := copyFile(oldFilePath, newFilePath)
		return err
	}
}

// CalculateUid 通过md5和sha1生成字符串
func CalculateUid(filePath string) (string, error) {
	stat, err := os.Stat(filePath)
	if err != nil {
		return "", err
	}
	// 大小文件用不同计算方式
	if stat.Size() < maxSilceSize {
		return calculateSmallFileUid(filePath)
	} else {
		return calculateLargeFileUid(filePath)
	}
}

// 计算小文件Uid
// 读取所有文件内容到内存中
func calculateSmallFileUid(filePath string) (uid string, err error) {
	context, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	md5Sum, sha1Sum := md5.Sum(context), sha1.Sum(context)
	return hex.EncodeToString(append(md5Sum[:], sha1Sum[:]...)), nil
}

// 计算大文件Uid
// 逐渐读取复制
func calculateLargeFileUid(filePath string) (uid string, e error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			e = err
		}
	}(file)
	md5Hash, sha1Hash := md5.New(), sha1.New()
	_, _ = io.Copy(md5Hash, file)
	_, _ = io.Copy(sha1Hash, file)
	return hex.EncodeToString(append(md5Hash.Sum(nil)[:], sha1Hash.Sum(nil)[:]...)), e
}

// 文件复制
func copyFile(src, dst string) (written int64, err error) {
	file1, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	file2, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return 0, err
	}
	defer func(file1 *os.File) {
		err = file1.Close()
	}(file1)
	defer func(file2 *os.File) {
		err = file2.Close()
	}(file2)
	return io.Copy(file2, file1)
}
