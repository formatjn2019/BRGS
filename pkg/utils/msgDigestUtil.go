package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"io/ioutil"
	"os"
)

// 采用文件大小，md5,sha1方式比较文件是否相同
// 注 空不等于空
func CompareFile(lpath, rpath string) bool {
	lst, _ := os.Stat(lpath)
	rst, _ := os.Stat(rpath)
	// 文件不相同情况
	// 简单判定
	// 其中一个文件不存在 一为文件，一为文件夹	皆为文件，文件大小不同
	if lst == nil || rst == nil || lst.IsDir() != rst.IsDir() || (!lst.IsDir() && (rst.Size() != lst.Size())) {
		return false
	} else if !lst.IsDir() && !rst.IsDir() {
		// 信息摘要判定
		ctxl, _ := ioutil.ReadFile(lpath)
		ctxr, _ := ioutil.ReadFile(rpath)
		return md5.Sum(ctxl) == md5.Sum(ctxr) && sha1.Sum(ctxl) == sha1.Sum(ctxr)
	} else {
		//同为文件夹则不进行对比
		return true
	}
}
