package util

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
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

// 从压缩包中还原文件
func RecoverFromArchive(zipPath, targetPath string) error {
	err := os.RemoveAll(targetPath)
	if err != nil {
		log.Panicln("解压路径错误", err)
		return err
	}
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		log.Println("打开压缩包错误", err.Error())
		return err
	}
	defer zipReader.Close()
	for _, file := range zipReader.File {
		path := filepath.Join(targetPath, filepath.Join(strings.Split(file.Name, "/")...))
		if file.FileInfo().IsDir() {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				log.Println("创建文件夹失败", err.Error())
				return err
			}
			continue
		}
		//确保上层目录存在
		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			log.Println("创建文件夹失败", err.Error())
			return err
		}
		dstFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			log.Println("创建文件失败", err.Error())
			return err
		}

		zipedFile, err := file.Open()
		if err != nil {
			log.Println("读取压缩包文件失败", err.Error())
			return err
		}

		if _, err := io.Copy(dstFile, zipedFile); err != nil {
			log.Println("解压文件失败", err.Error())
			return err
		}
		dstFile.Close()
		zipedFile.Close()
	}

	return nil
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
		} else if stat == nil {
			synced[file] = true
			fmt.Println("file", file)
			fmt.Println("error", err)
			fmt.Println("stat", stat)
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

// 将字典中的文件写入压缩包
func WriteZip(fileName string, mapDic map[string]string) error {
	file, err := os.Create(fileName)
	fmt.Println(fileName)
	if err != nil {
		log.Println("压缩文件创建错误", err.Error())
		return err
	}
	defer file.Close()
	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()
	for path, target := range mapDic {
		iowriter, err := zipWriter.Create(target)
		if err != nil {
			if os.IsPermission(err) {
				log.Println("权限不足: ", err.Error())
				return err
			}
			log.Printf("创建文件失败 %s error: %s\n", target, err.Error())
			return err
		}
		path, err := filepath.Abs(path)
		if err != nil {
			log.Println("压缩文件路径错误", err.Error())
			return err
		}
		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.Printf("读取文件失败 %s error: %s\n", path, err.Error())
			return err
		} else {
			iowriter.Write(content)
		}
	}
	return nil
}

// 读取csv文件并翻译头
func ReadCsvAsDictAndTranslate(filePath string, tralslateHeadDic map[string]string) (result []map[string]string, e error) {
	originDicList, err := ReadCsvAsDict(filePath)
	if err != nil {
		return originDicList, err
	}
	for _, oriDic := range originDicList {
		tmpDic := map[string]string{}
		for k, v := range oriDic {
			tmpDic[tralslateHeadDic[k]] = v
		}
		result = append(result, tmpDic)
	}
	return result, nil
}

// 读取csv文件
func ReadCsvAsDict(filePath string) (result []map[string]string, e error) {
	file, _ := os.OpenFile(filePath, os.O_RDONLY, 438)
	defer file.Close()
	file.Seek(3, 0)
	r := csv.NewReader(file)
	// 行首
	heads, err := r.Read()

	if err != nil {
		return nil, err
	}
	// 内容
	for record, err := r.Read(); err != io.EOF; record, err = r.Read() {
		if err != nil {
			return
		}
		dic := map[string]string{}
		for index, key := range heads {
			dic[key] = record[index]
		}
		result = append(result, dic)
	}
	return result, nil
}

// 将字典写入csv
func WriteCsvWithDict(filePath string, content []map[string]string, order ...string) error {
	marix := make([][]string, 0)
	headLine := make([]string, 0)
	if len(order) != 0 {
		headLine = order
	} else {
		for k := range content[0] {
			headLine = append(headLine, k)
		}
	}
	marix = append(marix, headLine)
	for _, lineDic := range content {
		nLine := make([]string, 0)
		for _, head := range headLine {
			if content, ok := lineDic[head]; ok {
				nLine = append(nLine, content)
			} else {
				nLine = append(nLine, "")
			}
		}
		marix = append(marix, nLine)
	}
	fmt.Println(marix)
	file, _ := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 438)
	defer file.Close()
	//excel 乱码问题，插入头
	// file.WriteString("\xEF\xBB\xBF")
	file.WriteString("\uFEFF")

	w := csv.NewWriter(file)
	w.WriteAll(marix)
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
	return nil
}
