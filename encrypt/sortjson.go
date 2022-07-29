/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    sortjson
 * @Date:    2022/7/17 19:09
 * @package: encrypt
 * @Version: x.x.x
 *
 * @Description: xxx
 *
 */

package encrypt

import (
	"fmt"
	"reflect"
	"sort"
)

func MarshalStruct2JsonSortByKey(v interface{}, unlessKeys map[string]bool) (string, error) {
	var args []*struct {
		key   string
		value string
	}
	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	kd := val.Kind()
	if kd != reflect.Struct {
		return "", fmt.Errorf("InvalidParamType=%v", kd)
	}
	num := val.NumField()
	for i := 0; i < num; i++ {
		tagVal := typ.Field(i).Tag.Get("json")
		value := val.Field(i)
		if tagVal != "" && !unlessKeys[tagVal] {
			arg := &struct {
				key   string
				value string
			}{key: tagVal}
			var err error
			arg.value, err = MarshalValue(value, unlessKeys)
			if err != nil {
				return "", err
			}
			args = append(args, arg)
		}
	}
	sort.Slice(args, func(i, j int) bool {
		return args[i].key < args[j].key
	})

	result := "{"
	for i, arg := range args {
		if i == 0 {
			result += fmt.Sprintf("\"%s\":%v", arg.key, arg.value)
		} else {
			result += fmt.Sprintf(",\"%s\":%v", arg.key, arg.value)
		}
	}
	result += "}"
	return result, nil
}

func MarshalValue(value reflect.Value, unlessKeys map[string]bool) (string, error) {
	var result string
	switch value.Kind() {
	case reflect.String:
		result = fmt.Sprintf("\"%s\"", value.String())
	case reflect.Struct:
		obj, err := MarshalStruct2JsonSortByKey(value.Interface(), unlessKeys)
		if err != nil {
			return "", err
		}
		result = fmt.Sprintf("%v", obj)
	case reflect.Array, reflect.Slice:
		cnt := value.Len()
		result = "["
		for i := 0; i < cnt; i++ {
			vv := value.Index(i)
			v, err := MarshalValue(vv, unlessKeys)
			if err != nil {
				return "", err
			}
			if i == 0 {
				result += v
			} else {
				result += "," + v
			}
		}
		result += "]"
	default:
		result = fmt.Sprintf("%v", value.Interface())
	}
	return result, nil
}
