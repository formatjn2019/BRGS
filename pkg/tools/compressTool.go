package tools

import (
	"BRGS/pkg/e"
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func GenerateNameByTime(name string) string {
	return fmt.Sprintf("%s_%s", name, time.Now().Format("20060102_150405"))
}

func GenerateRule(name string) *regexp.Regexp {
	return regexp.MustCompile(name + `_20\d{6}_\d{6}(.zip)?`)
}

// RecoverFromArchive 从压缩包中还原文件
func RecoverFromArchive(zipPath, targetPath string) (err error) {
	var logMessage string

	zipReader, errs := zip.OpenReader(zipPath)
	defer func(zipReader *zip.ReadCloser) {
		err = zipReader.Close()
	}(zipReader)
	errN := removeSub(targetPath)

	if err = errN; err != nil {
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
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				logMessage = "创建文件夹失败"
				goto errorLog
			}
			continue
		}
		//确保上层目录存在
		if err = os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			logMessage = "创建文件夹失败"
			goto errorLog
		}
		dstFile, errN := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err = errN; err != nil {
			logMessage = "创建文件失败"
			goto errorLog
		}

		zipedFile, errN := file.Open()
		if err = errN; err != nil {
			logMessage = "读取压缩包文件失败"
			goto errorLog
		}
		_, errN = io.Copy(dstFile, zipedFile)
		if err = errN; err != nil {
			logMessage = "解压文件失败"
			goto errorLog
		}
		err = dstFile.Close()
		if err != nil {
			logMessage = "关闭目标文件失败"
			goto errorLog
		}
		err = zipedFile.Close()
		if err != nil {
			logMessage = "关闭zip文件失败"
			goto errorLog
		}
	}
	return nil

errorLog:
	log.Fatal(e.TranslateToError(e.ErrorSyncFailDecompress, logMessage, err.Error()))
	return err
}

// WriteZip 将字典中的文件写入压缩包
func WriteZip(fileName string, mapDic map[string]string) (err error) {
	file, err := os.Create(fileName)
	defer func(file *os.File) {
		err = file.Close()
	}(file)
	var logMessage string
	zipWriter := zip.NewWriter(file)
	defer func(zipWriter *zip.Writer) {
		err = zipWriter.Close()
	}(zipWriter)

	if err != nil {
		logMessage = "压缩文件创建错误"
		goto errorLog
	}

	for path, target := range mapDic {
		ioWriter, errN := zipWriter.Create(target)
		if err = errN; err != nil {
			if os.IsPermission(err) {
				logMessage = "权限不足: "
				goto errorLog
			}
			logMessage = "创建文件失败 %s error: %s\n"
			goto errorLog
		}
		path, errN := filepath.Abs(path)
		if err = errN; err != nil {
			logMessage = "压缩文件路径错误"
			goto errorLog
		}
		content, errN := os.ReadFile(path)
		if err = errN; err != nil {
			logMessage = "读取文件失败" + path
			goto errorLog
		} else {
			_, errN := ioWriter.Write(content)
			if err = errN; err != nil {
				logMessage = "写入失败" + path
				return err
			}
		}
	}
	return nil
errorLog:
	log.Fatal(e.TranslateToError(e.ErrorSyncFailCompress, logMessage, err.Error()))
	return err
}
