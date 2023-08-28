package merge

import (
	"reflect"
)

const (
	Concatenate                  = "concatenate"
	RemoveDuplicates             = "remove_duplicates"
	Overwrite                    = "overwrite"
	ReplaceByIndexPreferRight    = "replace_by_index_prefer_right"
	ReplaceByIndexPreferLeft     = "replace_by_index_prefer_left"
	ReplaceByIndexPreferMax      = "replace_by_index_prefer_max"
	ReplaceByIndexPreferRightRec = "replace_by_index_prefer_right_rec"
	ReplaceByIndexPreferLeftRec  = "replace_by_index_prefer_left_rec"
	ReplaceByIndexPreferMaxRec   = "replace_by_index_prefer_max_rec"
)

func MergeMapStructure(m1, m2 map[string]interface{}, mergeType string) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range m1 {
		result[k] = v
	}

	for k, v := range m2 {
		switch mergeType {
		case ReplaceByIndexPreferRightRec, ReplaceByIndexPreferLeftRec, ReplaceByIndexPreferMaxRec:
			if reflect.TypeOf(v).Kind() == reflect.Slice && reflect.TypeOf(result[k]).Kind() == reflect.Slice {
				result[k] = mergeSlice(result[k], v, mergeType)
			} else if reflect.TypeOf(v).Kind() == reflect.Map && reflect.TypeOf(result[k]).Kind() == reflect.Map {
				result[k] = MergeMapStructure(result[k].(map[string]interface{}), v.(map[string]interface{}), mergeType)
			} else {
				result[k] = v
			}
		default:
			switch reflect.TypeOf(v).Kind() {
			case reflect.Map:
				if reflect.TypeOf(result[k]).Kind() != reflect.Map {
					result[k] = make(map[string]interface{})
				}
				result[k] = MergeMapStructure(result[k].(map[string]interface{}), v.(map[string]interface{}), mergeType)
			case reflect.Slice:
				if reflect.TypeOf(result[k]).Kind() != reflect.Slice {
					// 创建一个类型为v的slice
					sliceType := reflect.SliceOf(reflect.TypeOf(v).Elem())
					result[k] = reflect.MakeSlice(sliceType, 0, 0).Interface()
				}
				result[k] = mergeSlice(result[k], v, mergeType)
			default:
				result[k] = v
			}
		}
	}

	return result
}

// removeDuplicates 去重
func removeDuplicates(s1, s2 interface{}) []interface{} {
	// 将s1和s2转换为reflect.Value类型
	s1Value := reflect.ValueOf(s1)
	s2Value := reflect.ValueOf(s2)

	// 创建一个map来记录已经出现过的元素
	seen := make(map[interface{}]bool)
	result := make([]interface{}, 0)

	// 遍历s1中的元素，并添加到result中（如果未出现过）
	for i := 0; i < s1Value.Len(); i++ {
		v := s1Value.Index(i).Interface()
		if !seen[v] {
			result = append(result, v)
			seen[v] = true
		}
	}

	// 遍历s2中的元素，并添加到result中（如果未出现过）
	for i := 0; i < s2Value.Len(); i++ {
		v := s2Value.Index(i).Interface()
		if !seen[v] {
			result = append(result, v)
			seen[v] = true
		}
	}

	return result
}

// mergeSlice 合并两个Slice
func mergeSlice(s1, s2 interface{}, mergeType string) interface{} {
	var result interface{}

	switch mergeType {
	case Concatenate:
		// 使用反射将s1和s2连接在一起
		result = reflect.AppendSlice(reflect.ValueOf(s1), reflect.ValueOf(s2)).Interface()
	case RemoveDuplicates:
		// 使用反射去重
		result = removeDuplicates(s1, s2)
	case Overwrite:
		result = s2
	case ReplaceByIndexPreferRight, ReplaceByIndexPreferRightRec:
		result = replaceByIndex(s1, s2, reflect.ValueOf(s2).Len())
	case ReplaceByIndexPreferLeft, ReplaceByIndexPreferLeftRec:
		result = replaceByIndex(s1, s2, reflect.ValueOf(s1).Len())
	case ReplaceByIndexPreferMax, ReplaceByIndexPreferMaxRec:
		maxLen := reflect.ValueOf(s1).Len()
		if reflect.ValueOf(s2).Len() > maxLen {
			maxLen = reflect.ValueOf(s2).Len()
		}
		result = replaceByIndex(s1, s2, maxLen)
	}

	return result
}

// replaceByIndex 按索引位替换Slice元素
func replaceByIndex(s1, s2 interface{}, length int) interface{} {
	// 创建一个类型为s1的slice
	sliceType := reflect.TypeOf(s1)
	result := reflect.MakeSlice(sliceType, length, length).Interface()

	// 复制s1和s2到result
	reflect.Copy(reflect.ValueOf(result), reflect.ValueOf(s1))
	reflect.Copy(reflect.ValueOf(result), reflect.ValueOf(s2))

	return result
}
