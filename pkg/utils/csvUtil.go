package utils

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

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
	// fmt.Println(marix)
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
