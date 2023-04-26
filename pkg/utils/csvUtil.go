package utils

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

// ReadCsvAsDictAndTranslate 读取csv文件并翻译头
func ReadCsvAsDictAndTranslate(filePath string, translateHeadDic map[string]string) (result []map[string]string, e error) {
	originDicList, err := ReadCsvAsDict(filePath)
	if err != nil {
		return originDicList, err
	}
	for _, oriDic := range originDicList {
		tmpDic := map[string]string{}
		for k, v := range oriDic {
			tmpDic[translateHeadDic[k]] = v
		}
		result = append(result, tmpDic)
	}
	return result, nil
}

// ReadCsvAsDict 读取csv文件
func ReadCsvAsDict(filePath string) (result []map[string]string, e error) {
	file, _ := os.OpenFile(filePath, os.O_RDONLY, 438)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			e = err
		}
	}(file)
	_, err := file.Seek(3, 0)
	if err != nil {
		return nil, err
	}
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

// WriteCsvWithDict 将字典写入csv
func WriteCsvWithDict(filePath string, content []map[string]string, order ...string) (e error) {
	matrix := make([][]string, 0)
	headLine := make([]string, 0)
	if len(order) != 0 {
		headLine = order
	} else {
		for k := range content[0] {
			headLine = append(headLine, k)
		}
	}
	matrix = append(matrix, headLine)
	for _, lineDic := range content {
		nLine := make([]string, 0)
		for _, head := range headLine {
			if content, ok := lineDic[head]; ok {
				nLine = append(nLine, content)
			} else {
				nLine = append(nLine, "")
			}
		}
		matrix = append(matrix, nLine)
	}
	file, _ := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 438)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			e = err
		}
	}(file)
	//excel 乱码问题，插入头
	// file.WriteString("\xEF\xBB\xBF")
	_, err := file.WriteString("\uFEFF")
	if err != nil {
		return err
	}

	w := csv.NewWriter(file)
	err = w.WriteAll(matrix)
	if err != nil {
		return err
	}
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
	return nil
}
