package itm

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func Marshal(v interface{}) (map[string]string, error) {
	var result = map[string]string{}
	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	kd := val.Kind()
	if kd == reflect.Pointer {
		typ = typ.Elem()
		val = val.Elem()
		kd = val.Kind()
	}
	switch kd {
	case reflect.Struct:
		num := val.NumField()
		for i := 0; i < num; i++ {
			tagVal := typ.Field(i).Tag.Get("json")
			value := val.Field(i)
			if tagVal != "" {
				vv, err := MarshalValue(value)
				if err != nil {
					return nil, err
				}
				result[tagVal] = vv
			}
		}
	case reflect.Map:
		it := val.MapRange()
		for it.Next() {
			k := it.Key()
			vv, err := MarshalValue(it.Value())
			if err != nil {
				return nil, err
			}
			result[k.String()] = vv
		}
	}
	return result, nil
}

func MarshalValue(value reflect.Value) (string, error) {
	var result string
	switch value.Kind() {
	case reflect.String:
		return value.String(), nil
	case reflect.Struct, reflect.Map:
		result, err := json.Marshal(value.Interface())
		return string(result), err
	case reflect.Interface:
		return MarshalValue(reflect.ValueOf(value.Interface()))
	case reflect.Array, reflect.Slice:
		cnt := value.Len()
		result = "["
		for i := 0; i < cnt; i++ {
			vv := value.Index(i)
			v, err := MarshalValue(vv)
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
