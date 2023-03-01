package tools

import (
	"BRGS/pkg/e"
	"archive/zip"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// RecoverFromArchive 从压缩包中还原文件
func RecoverFromArchive(zipPath, targetPath string) (err error) {
	var logMessage string

	zipReader, errs := zip.OpenReader(zipPath)
	defer zipReader.Close()
	errt := os.RemoveAll(targetPath)

	if err = errt; err != nil {
		logMessage = "解压路径错误"
		goto errorLog
	}
	if err = errs; err != nil {
		logMessage = "打开压缩包错误"
		goto errorLog
	}
	for _, file := range zipReader.File {
		path := filepath.Join(targetPath, filepath.Join(strings.Split(file.Name, "/")...))
		if file.FileInfo().IsDir() {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				logMessage = "创建文件夹失败"
				goto errorLog
			}
			continue
		}
		//确保上层目录存在
		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			logMessage = "创建文件夹失败"
			goto errorLog
		}
		dstFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			logMessage = "创建文件失败"
			goto errorLog
		}

		zipedFile, err := file.Open()
		if err != nil {
			logMessage = "读取压缩包文件失败"
			goto errorLog
		}

		if _, err := io.Copy(dstFile, zipedFile); err != nil {
			logMessage = "解压文件失败"
			goto errorLog
		}
		dstFile.Close()
		zipedFile.Close()
	}

	return nil

errorLog:
	log.Fatal(e.TranslateToError(e.ErrorSyncFailDecompress, logMessage, err.Error()))
	return err
}

// WriteZip 将字典中的文件写入压缩包
func WriteZip(fileName string, mapDic map[string]string) (err error) {

	file, err := os.Create(fileName)
	var logMessage string
	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	if err != nil {
		logMessage = "压缩文件创建错误"
		goto errorLog
	}
	defer file.Close()

	for path, target := range mapDic {
		iowriter, err := zipWriter.Create(target)
		if err != nil {
			if os.IsPermission(err) {
				logMessage = "权限不足: "
				goto errorLog
			}
			logMessage = "创建文件失败 %s error: %s\n"
			goto errorLog
		}
		path, err := filepath.Abs(path)
		if err != nil {
			logMessage = "压缩文件路径错误"
			goto errorLog
		}
		content, err := ioutil.ReadFile(path)
		if err != nil {
			logMessage = "读取文件失败" + path
			goto errorLog
		} else {
			iowriter.Write(content)
		}
	}
	return nil
errorLog:
	log.Fatal(e.TranslateToError(e.ErrorSyncFailCompress, logMessage, err.Error()))
	return err
}
