package tools

import (
	"BRGS/pkg/e"
)

// TranslateList 翻译数组
func TranslateList(origin []string, dict map[string]string) ([]string, error) {
	result := make([]string, 0, len(origin))
	for _, word := range origin {
		if translated, ok := dict[word]; ok {
			result = append(result, translated)
		} else {
			return nil, e.TranslateToError(e.ErrorTranslate)
		}
	}
	return result, nil
}

// DictReverse 字典翻转
func DictReverse(dict map[string]string) map[string]string {
	result := map[string]string{}
	for k, v := range dict {
		result[v] = k
	}
	return result
}
