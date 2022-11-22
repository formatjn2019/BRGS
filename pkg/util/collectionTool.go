package util

import "errors"

// 翻译数组
func TranslateList(origin []string, dict map[string]string) ([]string, error) {
	result := make([]string, 0, len(origin))
	for _, word := range origin {
		if translated, ok := dict[word]; ok {
			result = append(result, translated)
		} else {
			return nil, errors.New("错误，待定")
		}
	}
	return result, nil
}

// 字典翻转
func DictReverse(dict map[string]string) map[string]string {
	result := map[string]string{}
	for k, v := range dict {
		result[v] = k
	}
	return result
}
