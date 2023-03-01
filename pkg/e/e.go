package e

import (
	"errors"
	"fmt"
	"strings"
)

// 功能划分
const (
	ErrorFunction = iota + 1
	ErrorIo
	ErrorSync
	ErrorInputCheck
	ErrorOperationUnsupport
)

// 功能错误
const (
	ErrorTranslate = (iota+1)<<3 | ErrorFunction
	ErrorReadConfig
)

// IO错误
const (
	ErrorRead = (iota+1)<<3 | ErrorIo
	ErrorWrite
	ErrorCreate
	ErrorDelete
	ErrorMove
	ErrorUpdate
)

// 同步错误
const (
	ErrorSyncInit = (iota+1)<<3 | ErrorSync
	ErrorSyncScan
	ErrorSyncWatch
	ErrorSyncChangeStats
	ErrorSyncFailByFileUpdate
	ErrorSyncFailCompress
	ErrorSyncFailDecompress
)

// 输入错误
const (
	ErrorEmptyString = (iota+1)<<3 | ErrorInputCheck
	ErrorNotDir
	ErrorInputNumTooLarge
	ErrorInputNumTooSmall
	ErrorInterval
	ErrorSameName
)

var translateDic = map[int]string{
	ErrorFunction:             "功能错误",
	ErrorReadConfig:           "读取配置文件错误",
	ErrorIo:                   "IO错误",
	ErrorSync:                 "同步错误",
	ErrorInputCheck:           "输入错误",
	ErrorOperationUnsupport:   "不支持的操作",
	ErrorTranslate:            "翻译",
	ErrorRead:                 "读取",
	ErrorWrite:                "写入",
	ErrorCreate:               "创建",
	ErrorDelete:               "删除",
	ErrorMove:                 "移动",
	ErrorUpdate:               "更新",
	ErrorSyncInit:             "同步初始化",
	ErrorSyncScan:             "扫描文件夹",
	ErrorSyncWatch:            "监控文件夹",
	ErrorSyncChangeStats:      "更改同步状态",
	ErrorSyncFailByFileUpdate: "同步时更新文件",
	ErrorSyncFailCompress:     "压缩文件",
	ErrorSyncFailDecompress:   "解压文件",
	ErrorEmptyString:          "空字符串",
	ErrorNotDir:               "非文件夹",
	ErrorInputNumTooLarge:     "数字过大",
	ErrorInputNumTooSmall:     "数字过小",
	ErrorInterval:             "同步间隔大于归档间隔",
	ErrorSameName:             "名称重复",
}

func TranslateError(code int) string {
	if code < 8 {
		return translateDic[code]
	} else {
		return fmt.Sprintf("%s--%s", translateDic[code%8], translateDic[code])
	}
}

func TranslateToError(code int, context ...string) error {
	if len(context) > 0 {
		return errors.New(TranslateError(code) + "\t: " + strings.Join(context, "\n\t\t"))
	}
	return errors.New(TranslateError(code))
}
