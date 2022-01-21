/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    utils
 * @Date:    2022/1/21 5:47 下午
 * @package: utils
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func ToInt(i interface{}) int {
	switch i.(type) {
	case int:
		return i.(int)
	default:
		return int(ToInt64(i))
	}
}

func ToIntSlice(i interface{}) []int {
	switch i.(type) {
	case []int:
		return i.([]int)
	}

	vs := ToInt64Slice(i)
	var vss []int
	for _, v := range vs {
		vss = append(vss, int(v))
	}
	return vss
}

func ToInt64(i interface{}) int64 {
	switch i.(type) {
	case uint:
		return int64(i.(uint))
	case int:
		return int64(i.(int))
	case uint8:
		return int64(i.(uint8))
	case int8:
		return int64(i.(int8))
	case uint16:
		return int64(i.(uint16))
	case int16:
		return int64(i.(int16))
	case uint32:
		return int64(i.(uint32))
	case int32:
		return int64(i.(int32))
	case uint64:
		return int64(i.(uint64))
	case int64:
		return i.(int64)
	}

	ii, _ := strconv.ParseInt(fmt.Sprintf("%v", i), 10, 64)
	return ii
}

func ToInt64Slice(i interface{}) []int64 {
	var vs []int64
	switch i.(type) {
	case []uint:
		for _, v := range i.([]uint) {
			vs = append(vs, int64(v))
		}
		return vs
	case []int:
		for _, v := range i.([]int) {
			vs = append(vs, int64(v))
		}
		return vs
	case []uint8:
		for _, v := range i.([]uint8) {
			vs = append(vs, int64(v))
		}
		return vs
	case []int8:
		for _, v := range i.([]int8) {
			vs = append(vs, int64(v))
		}
		return vs
	case []uint16:
		for _, v := range i.([]uint16) {
			vs = append(vs, int64(v))
		}
		return vs
	case []int16:
		for _, v := range i.([]int16) {
			vs = append(vs, int64(v))
		}
		return vs
	case []uint32:
		for _, v := range i.([]uint32) {
			vs = append(vs, int64(v))
		}
		return vs
	case []int32:
		for _, v := range i.([]int32) {
			vs = append(vs, int64(v))
		}
		return vs
	case []uint64:
		for _, v := range i.([]uint64) {
			vs = append(vs, int64(v))
		}
		return vs
	case []int64:
		return i.([]int64)
	case []string:
		strs := i.([]string)
		for _, str := range strs {
			ii, err := strconv.ParseInt(str, 10, 64)
			if err == nil {
				vs = append(vs, ii)
			}
		}
		return vs
	}

	strs := strings.Split(fmt.Sprintf("%v", i), " ")
	for _, str := range strs {
		ii, err := strconv.ParseInt(str, 10, 64)
		if err == nil {
			vs = append(vs, ii)
		}
	}

	return vs
}

func ToFloat64(i interface{}) float64 {
	switch i.(type) {
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
		return float64(ToInt64(i))
	case float32:
		return float64(i.(float32))
	case float64:
		return i.(float64)
	}
	v, _ := strconv.ParseFloat(fmt.Sprintf("%v", i), 64)
	return v
}

func ToFloat32(i interface{}) float32 {
	switch i.(type) {
	case float32:
		return i.(float32)
	}
	return float32(ToFloat64(i))
}

func ToFloat64Slice(i interface{}) []float64 {
	switch i.(type) {
	case []float64:
		return i.([]float64)
	case []float32:
		var vs []float64
		for _, v := range i.([]float32) {
			vs = append(vs, float64(v))
		}
		return vs
	case []string:
		var vs []float64
		strs := i.([]string)
		for _, str := range strs {
			ii, err := strconv.ParseFloat(str, 64)
			if err == nil {
				vs = append(vs, ii)
			}
		}
	case []int8, []uint8, []int16, []uint16, []int, []uint, []int32, []uint32, []int64, []uint64:
		vss := ToInt64Slice(i)
		var vs []float64
		for _, v := range vss {
			vs = append(vs, float64(v))
		}
		return vs
	}

	var vs []float64
	strs := strings.Split(fmt.Sprintf("%v", i), " ")
	for _, str := range strs {
		ii, err := strconv.ParseFloat(str, 64)
		if err == nil {
			vs = append(vs, ii)
		}
	}

	return vs
}

func ToFloat32Slice(i interface{}) []float32 {
	switch i.(type) {
	case []float32:
		return i.([]float32)
	}

	vs := ToFloat64Slice(i)
	var vss []float32
	for _, v := range vs {
		vss = append(vss, float32(v))
	}
	return vss
}
