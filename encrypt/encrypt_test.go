/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    encrypt_test
 * @Date:    2021/12/6 2:26 下午
 * @package: encrypt
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package encrypt

import (
	"fmt"
	"git.hawtech.cn/jager/hawox/logx"
	"git.hawtech.cn/jager/hawox/rsa"
	"testing"
)

func Test_Encrypt(t *testing.T) {
	var str = "fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!fuck!!"
	fmt.Printf("源数据：%v 长度=%d\n", []byte(str), len(str))
	bt, err := rsa.DefaultEncrypt([]byte(str))
	if err != nil {
		logx.Fatal().Err(err).Send()
	}
	fmt.Printf("加密数据：%v\n", bt)
	ss, err := rsa.DefaultDecrypt(bt)
	if err != nil {
		logx.Fatal().Err(err).Send()
	}
	fmt.Printf("加密解密后数据：%s\n", string(ss))
	fmt.Printf("源数据长度=%d 加密后长度=%d\n", len(ss), len(bt))
}
