/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    default
 * @Date:    2021/8/27 11:39 上午
 * @package: sensitive
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package sensitive

import (
	"io"
)

var filter_ *Filter

func Initialize(path string) error {
	filter_ = New()
	return filter_.LoadWordDict(path)
}

// UpdateNoisePattern 更新去噪模式
func UpdateNoisePattern(pattern string) {
	filter_.UpdateNoisePattern(pattern)
}

// LoadWordDict 加载敏感词字典
func LoadWordDict(path string) error {
	return filter_.LoadWordDict(path)
}

// LoadNetWordDict 加载网络敏感词字典
func LoadNetWordDict(url string) error {
	return filter_.LoadNetWordDict(url)
}

// Load common method to add words
func Load(rd io.Reader) error {
	return filter_.Load(rd)
}

// AddWord 添加敏感词
func AddWord(words ...string) {
	filter_.trie.Add(words...)
}

// filter_ 过滤敏感词
func FilterTxt(text string) string {
	return filter_.Filter(text)
}

// Replace 和谐敏感词
func Replace(text string, repl rune) string {
	return filter_.Replace(text, repl)
}

// FindIn 检测敏感词
func FindIn(text string) (bool, string) {
	return filter_.FindIn(text)
}

// FindAll 找到所有匹配词
func FindAll(text string) []string {
	return filter_.FindAll(text)
}

// Validate 检测字符串是否合法
func Validate(text string) (bool, string) {
	return filter_.Validate(text)
}

// RemoveNoise 去除空格等噪音
func RemoveNoise(text string) string {
	return filter_.RemoveNoise(text)
}
