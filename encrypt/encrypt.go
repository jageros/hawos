/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    utils
 * @Date:    2021/8/24 10:26 上午
 * @package: utils
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package encrypt

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func Md5(src string) string {
	data := []byte(src)
	has := md5.Sum(data)
	md5Str := hex.EncodeToString(has[:])
	return md5Str
}

func Md516(src string) string {
	data := []byte(src)
	has := md5.Sum(data)
	md5Str := hex.EncodeToString(has[:])
	return md5Str[8:24]
}

func Md516Upper(src string) string {
	data := []byte(src)
	has := md5.Sum(data)
	md5Str := fmt.Sprintf("%X", has)
	return md5Str[8:24]
}

func Sha1(src string) string {
	h := sha1.New()
	h.Write([]byte(src))
	sh := hex.EncodeToString(h.Sum(nil))
	return sh
}

func Sha256(src string) string {
	m := sha256.New()
	m.Write([]byte(src))
	res := hex.EncodeToString(m.Sum(nil))
	return res
}

func HmacSha256Base64(src, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(src))
	sign := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return sign
}