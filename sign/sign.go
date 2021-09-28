/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    arg
 * @Date:    2021/8/26 3:02 下午
 * @package: sign
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package sign

import (
	"fmt"
	"github.com/jageros/hawox/encrypt"
	"sort"
	"strings"
)

// 1、参与签名参数：所有一级值为字符串和整型类型的参数

// 2、字符串拼接格式： 对key按字典序升序排列，中间用‘&’隔开， 形如：key1=value1&key2=value2&key3=value3
// ps: 所谓字典序升序排列，直观上就如同在字典中排列单词一样排序，按照字母表或数字表里递增顺序的排列次序，即先考虑第一个 “字母”，在相同的情况下考虑第二个 “字母”，依此类推。

// 3、加密算法： 用给定的secret进行HmacSHA256加密生成签名串

// 4、签名结果：将生成的签名串使用 Base64 进行编码

const (
	signArgKey = "signature"
	secret     = "cd803706cdc822e372fd7c73c0f109b9"
)

type Arg struct {
	key   string
	value string
}

type Args struct {
	argList    []Arg
	Signature  string
	SignArgKey string
	Secret     string
}

func defaultArgs() *Args {
	return &Args{
		SignArgKey: signArgKey,
		Secret:     secret,
	}
}

// ================== sort interface ====================

func (a *Args) Len() int {
	return len(a.argList)
}

func (a *Args) Less(i, j int) bool {
	return a.argList[i].key < a.argList[j].key
}

func (a *Args) Swap(i, j int) {
	a.argList[i], a.argList[j] = a.argList[j], a.argList[i]
}

// ========================================================

func NewArgs(opfs ...func(args *Args)) *Args {
	args := defaultArgs()
	for _, opf := range opfs {
		opf(args)
	}
	return args
}

func NewArgsWithMap(args map[string]interface{}, opfs ...func(ags *Args)) *Args {
	a := defaultArgs()
	for _, opf := range opfs {
		opf(a)
	}

	for key, val := range args {
		if key == a.SignArgKey {
			if v, ok := val.(string); ok {
				a.Signature = v
			}
			continue
		}
		var v string
		switch val.(type) {
		case uint8, uint16, uint32, uint64, int8, int16, int32, int64, int:
			v = fmt.Sprintf("%d", val)
		case string:
			v = val.(string)
		default:
			continue
		}
		a.argList = append(a.argList, Arg{
			key:   key,
			value: v,
		})
	}

	return a
}

func (a *Args) Add(key, value string) {
	a.argList = append(a.argList, Arg{
		key:   key,
		value: value,
	})
}

func (a *Args) Adds(args ...Arg) {
	a.argList = append(a.argList, args...)
}

func (a *Args) AddM(args map[string]string) {
	for key, val := range args {
		if key == a.SignArgKey {
			a.Signature = val
		} else {
			a.argList = append(a.argList, Arg{
				key:   key,
				value: val,
			})
		}
	}
}

func (a *Args) AddMi(args map[string]interface{}) {
	for key, val := range args {
		if key == a.SignArgKey {
			if v, ok := val.(string); ok {
				a.Signature = v
			}
			continue
		}
		var v string
		switch val.(type) {
		case uint8, uint16, uint32, uint64, int8, int16, int32, int64, int:
			v = fmt.Sprintf("%v", val)
		case float32, float64:
			v = fmt.Sprintf("%v", val)
		case string:
			v = val.(string)
		}
		a.argList = append(a.argList, Arg{
			key:   key,
			value: v,
		})
	}
}

func (a *Args) GenSignature() string {
	sort.Sort(a)
	var args []string
	for _, arg := range a.argList {
		args = append(args, fmt.Sprintf("%s=%s", arg.key, arg.value))
	}
	src := strings.Join(args, "&")
	fmt.Println(src)
	return encrypt.HmacSha256Base64(src, a.Secret)
}

func (a *Args) VerifySign() bool {
	return a.Signature == a.GenSignature()
}
