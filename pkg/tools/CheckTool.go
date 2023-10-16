package tools

import (
	"BRGS/pkg/e"
	"os"
	"strconv"
)

// Check 核验方法生成器
type Check func(string) error

// GenerateNotNullCheck 非空核验
func GenerateNotNullCheck() Check {
	return func(s string) error {
		if len(s) == 0 {
			return e.TranslateToError(e.ErrorEmptyString)
		}
		return nil
	}
}

// GeneratePathCheck 路径核验
func GeneratePathCheck() Check {
	return func(str string) error {
		//空字符串，卫语句
		if str == "" {
			return nil
		}
		if info, err := os.Stat(str); err != nil {
			return err
		} else if !info.IsDir() {
			return e.TranslateToError(e.ErrorNotDir)
		}
		return nil
	}
}

// GenerateRangeCheck 字符数值核验
func GenerateRangeCheck(zero bool, min, max int) Check {
	return func(s string) error {
		if num, err := strconv.Atoi(s); err != nil {
			return err
		} else if zero && num == 0 {
			return nil
		} else if num < min {
			return e.TranslateToError(e.ErrorInputNumTooSmall)
		} else if num > max {
			return e.TranslateToError(e.ErrorInputNumTooLarge)
		}
		return nil
	}
}
