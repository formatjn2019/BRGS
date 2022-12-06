package e

import (
	"errors"
	"fmt"
	"strings"
)

// 功能划分
const (
	ERROR_FUNCTION = iota + 1
	ERROR_IO
	ERROR_SYNC
	ERROR_INPUT_CHECK
	ERROR_OPERATION_UNSUPPORT
)

// 功能错误
const (
	ERROR_TRANSLATE = (iota+1)<<3 | ERROR_FUNCTION
)

// IO错误
const (
	ERROR_READ = (iota+1)<<3 | ERROR_IO
	ERROR_WRITE
	ERROR_CREATE
	ERROR_DELETE
	ERROR_MOVE
	ERROR_UPDATE
)

// 同步错误
const (
	ERROR_SYNC_INIT = (iota+1)<<3 | ERROR_SYNC
	ERROR_SYNC_SCAN
	ERROR_SYNC_WATCH
	ERROR_SYNC_CHANGE_STATS
	ERROR_SYNC_FAIL_BY_FILEUPDATE
	ERROR_SYNC_FAIL_COMPRESS
	ERROR_SYNC_FAIL_DECOMPRESS
)

// 输入错误
const (
	ERROR_EMPTY_STRING = (iota+1)<<3 | ERROR_INPUT_CHECK
	ERROR_NOT_DIR
	ERROR_INPUT_NUM_TOO_LARGE
	ERROR_INPUT_NUM_TOO_SMALL
	ERROR_INTERVAL
	ERROR_SAME_NAME
)

var translateDic = map[int]string{
	ERROR_FUNCTION:                "功能错误",
	ERROR_IO:                      "IO错误",
	ERROR_SYNC:                    "同步错误",
	ERROR_INPUT_CHECK:             "输入错误",
	ERROR_OPERATION_UNSUPPORT:     "不支持的操作",
	ERROR_TRANSLATE:               "翻译",
	ERROR_READ:                    "读取",
	ERROR_WRITE:                   "写入",
	ERROR_CREATE:                  "创建",
	ERROR_DELETE:                  "删除",
	ERROR_MOVE:                    "移动",
	ERROR_UPDATE:                  "更新",
	ERROR_SYNC_INIT:               "同步初始化",
	ERROR_SYNC_SCAN:               "扫描文件夹",
	ERROR_SYNC_WATCH:              "监控文件夹",
	ERROR_SYNC_CHANGE_STATS:       "更改同步状态",
	ERROR_SYNC_FAIL_BY_FILEUPDATE: "同步时更新文件",
	ERROR_SYNC_FAIL_COMPRESS:      "压缩文件",
	ERROR_SYNC_FAIL_DECOMPRESS:    "解压文件",
	ERROR_EMPTY_STRING:            "空字符串",
	ERROR_NOT_DIR:                 "非文件夹",
	ERROR_INPUT_NUM_TOO_LARGE:     "数字过大",
	ERROR_INPUT_NUM_TOO_SMALL:     "数字过小",
	ERROR_INTERVAL:                "同步间隔大于归档间隔",
	ERROR_SAME_NAME:               "名称重复",
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
