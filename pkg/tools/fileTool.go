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
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

// 大于指定大小分片哈希
// 20M
const maxSilceSize = 10 * 1024 * 1024

// SyncFile 文件同步
func SyncFile(pathSource, pathTarget string, addDic, delDic map[string]bool) (synced map[string]bool, err error) {
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
			fi, err := os.Stat(filepath.Join(pathSource, file))
			if err != nil {
				err = e.TranslateToError(e.ErrorRead, file, err.Error())
				goto errorLog
			}
			cont, err := os.ReadFile(filepath.Join(pathSource, file))
			if err != nil {
				err = e.TranslateToError(e.ErrorRead, file, err.Error())
				goto errorLog
			}
			err = os.WriteFile(filepath.Join(pathTarget, file), cont, fi.Mode())
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
// 注 使用新方法比原方法慢了近十倍
func WalkDir(root string) map[string]string {
	result := map[string]string{}
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			rl, _ := filepath.Rel(root, path)
			result[path] = strings.ReplaceAll(rl, "\\", "/")
		}
		return err
	})
	if err != nil {
		return nil
	}
	return result
}

// CompareDirs 文件对比
func CompareDirs(pathL, pathR string) bool {
	result := true
	lst, _ := os.Stat(pathL)
	rst, _ := os.Stat(pathR)
	//文件不相同情况
	if !utils.CompareFile(pathL, pathR) {
		fmt.Printf("对比失败 %s %s\n", lst, rst)
		fmt.Printf("路径分别为 %s %s\n", pathL, pathR)
		return false
	} else {
		fmt.Printf("%s %s 对比成功\n", pathL, pathR)
	}
	if lst != nil && rst != nil && lst.IsDir() && rst.IsDir() {
		next := map[string]bool{}
		files, _ := os.ReadDir(pathL)
		for _, file := range files {
			next[file.Name()] = true
		}
		files, _ = os.ReadDir(pathR)
		for _, file := range files {
			next[file.Name()] = true
		}
		for name := range next {
			if !CompareDirs(filepath.Join(pathL, name), filepath.Join(pathR, name)) {
				result = false
			}
		}
	}
	return result
}

// CalculateAllUid 计算路径下所以文件的uid
func CalculateAllUid(rootPath string, abstractPath bool, skipRule ...string) (map[string]string, error) {
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
			fileName := filepath.Base(path)
			//规则过滤
			for _, rule := range skipRule {
				if ok, _ := filepath.Match(rule, fileName); ok {
					return nil
				}
			}
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				c := make(chan string, 0)
				calculateChan <- c
				c <- path
				uid := <-c
				if uid == "" {
					err = e.TranslateToError(e.ErrorCalculate, "Fail to calculate", path)
					return
				}
				lock.Lock()
				defer lock.Unlock()
				if abstractPath {
					result[path] = uid
				} else {
					pathKey, _ := filepath.Rel(rootPath, path)
					result[pathKey] = uid
				}
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
		wg.Add(1)
		go func(path, uid string) {
			defer wg.Done()
			c := make(chan string, 0)
			calculateChan <- c
			c <- path
			uidN := <-c
			if uid == "" {
				err = e.TranslateToError(e.ErrorCalculate, "Fail to calculate", path)
				return
			}
			lock.Lock()
			defer lock.Unlock()
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

func ScanFilesToDic(paths ...string) (linkedFilesDic map[string]map[string]fs.FileInfo, err error) {
	// 运算操作通道
	calculateChan := make(chan chan string, 0)
	defer close(calculateChan)
	linkedFilesDic = map[string]map[string]fs.FileInfo{}
	// 运算操作线程池
	go CalculateUidPool(calculateChan)
	var lock sync.Mutex
	var wg sync.WaitGroup
	for _, rootPath := range paths {
		err = filepath.Walk(rootPath, func(path string, info fs.FileInfo, err error) error {
			if !info.IsDir() {
				if err != nil {
					return err
				}
				wg.Add(1)
				go func(path string, info fs.FileInfo) {
					defer wg.Done()
					c := make(chan string, 0)
					calculateChan <- c
					c <- path
					uid := <-c
					if uid == "" {
						err = e.TranslateToError(e.ErrorCalculate, "Fail to calculate", path)
						return
					}
					lock.Lock()
					defer lock.Unlock()
					if linkedFiles, ok := linkedFilesDic[uid]; !ok {
						//第一次出现，记录
						linkedFilesDic[uid] = map[string]fs.FileInfo{path: info}
					} else {
						//判断是否和已有文件是同一文件的硬链接
						for _, linkedInfo := range linkedFiles {
							if os.SameFile(linkedInfo, info) {
								return
							}
						}
						//非已记录文件，记录
						linkedFilesDic[uid][path] = info
					}
				}(path, info)
			}
			return err
		})
		wg.Wait()
		if err != nil {
			return nil, err
		}
	}
	return
}

// SyncAllFileWithHardLink 通过硬链接方式同步所有文件
func SyncAllFileWithHardLink(src, dst string, filesDic map[string]string, linkedFilesDic map[string]map[string]fs.FileInfo) error {
	var lock sync.Mutex
	for path, uid := range filesDic {
		srcPath, dstPath := filepath.Join(src, path), filepath.Join(dst, path)
		err := SyncFileWithHardLink(uid, srcPath, dstPath, linkedFilesDic, &lock)
		if err != nil {
			return err
		}
	}
	return nil
}

// SyncFileWithHardLink 同步文件
// 非新文件时，采用硬链接形式记录
func SyncFileWithHardLink(uid, oldFilePath, newFilePath string, linkedFilesDic map[string]map[string]fs.FileInfo, lock sync.Locker) (e error) {
	lock.Lock()
	defer lock.Unlock()
	createParentDir(newFilePath)
	//判断是否第一次出现该文件
	if linkedFiles, ok := linkedFilesDic[uid]; !ok {
		//第一次出现，拷贝
		_, err := copyFile(oldFilePath, newFilePath)
		info, _ := os.Stat(newFilePath)
		linkedFilesDic[uid] = map[string]fs.FileInfo{newFilePath: info}
		return err
	} else {
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
				return os.Rename(tempPath, newFilePath)
			}
		}
		//尝试生成硬链接失败,将该文件加入列表
		_, err := copyFile(oldFilePath, newFilePath)
		info, _ := os.Stat(newFilePath)
		linkedFiles[newFilePath] = info
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

func createParentDir(path string) {
	dstDir := filepath.Dir(path)
	if stat, err := os.Stat(dstDir); err != nil || stat == nil {
		err := os.MkdirAll(dstDir, 0755)
		if err != nil {
			panic(err)
		}
	}
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

// 移除子文件
func removeSub(path string) error {
	subs, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, sub := range subs {
		if err := os.RemoveAll(filepath.Join(path, sub.Name())); err != nil {
			return err
		}
	}
	return nil
}

func CloneDir(src, dst string) error {
	stat, err := os.Stat(dst)
	// 目标文件夹下如存在文件，则清空子文件夹
	if err == nil && stat.IsDir() {
		err = removeSub(dst)
		if err != nil {
			return err
		}
	}
	return cloneFileTree(src, dst)
}

// CloneFileTree 递归拷贝
func cloneFileTree(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if info.IsDir() {
		err = os.MkdirAll(dst, os.ModePerm)
		if err != nil {
			return err
		}
		subs, err := os.ReadDir(src)
		if err != nil {
			return err
		}
		for _, info := range subs {
			s, d := filepath.Join(src, info.Name()), filepath.Join(dst, info.Name())
			if info.IsDir() {
				if err = cloneFileTree(s, d); err != nil {
					return err
				}
			} else {
				if _, err = copyFile(s, d); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
