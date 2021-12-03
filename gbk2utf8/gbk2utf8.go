/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    gbk2utf8
 * @Date:    2021/12/1 10:57 上午
 * @package: gbk2utf8
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package gbk2utf8

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
)

func Decode(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}
