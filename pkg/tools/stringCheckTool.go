package tools

import (
	"BRGS/pkg/e"
	"os"
	"strconv"
)

type Check func(string) error

// 核验方法生成器
// 非空核验
func GenerateNotNullCheck() Check {
	return func(s string) error {
		if len(s) == 0 {
			return e.TranslateToError(e.ERROR_EMPTY_STRING)
		}
		return nil
	}
}

// 路径核验
func GeneratePathCheck() Check {
	return func(str string) error {
		//空字符串，卫语句
		if str == "" {
			return nil
		}
		if info, err := os.Stat(str); err != nil {
			return err
		} else if !info.IsDir() {
			return e.TranslateToError(e.ERROR_NOT_DIR)
		}
		return nil
	}
}

// 字符数值核验
func GenerateRangeCheck(zero bool, min, max int) Check {
	return func(s string) error {
		if num, err := strconv.Atoi(s); err != nil {
			return err
		} else if zero && num == 0 {
			return nil
		} else if num < min {
			return e.TranslateToError(e.ERROR_INPUT_NUM_TOO_SMALL)
		} else if num > max {
			return e.TranslateToError(e.ERROR_INPUT_NUM_TOO_LARGE)
		}
		return nil
	}
}
