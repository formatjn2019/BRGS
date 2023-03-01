package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"os"
)

// CompareFile 采用文件大小，md5,sha1方式比较文件是否相同
// 注 空不等于空
func CompareFile(lPath, rPath string) bool {
	lst, _ := os.Stat(lPath)
	rst, _ := os.Stat(rPath)
	// 文件不相同情况
	// 简单判定
	// 其中一个文件不存在 一为文件，一为文件夹	皆为文件，文件大小不同
	if lst == nil || rst == nil || lst.IsDir() != rst.IsDir() || (!lst.IsDir() && (rst.Size() != lst.Size())) {
		return false
	} else if !lst.IsDir() && !rst.IsDir() {
		// 信息摘要判定
		lContext, _ := os.ReadFile(lPath)
		rContext, _ := os.ReadFile(rPath)
		return md5.Sum(lContext) == md5.Sum(rContext) && sha1.Sum(lContext) == sha1.Sum(rContext)
	} else {
		//同为文件夹则不进行对比
		return true
	}
}
