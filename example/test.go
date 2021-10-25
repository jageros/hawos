/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    test
 * @Date:    2021/10/25 3:34 下午
 * @package: example
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package main

import (
	"fmt"
	"github.com/jageros/hawox/rsa"
	"log"
)

func main() {
	var str = "xxxxxxxwhact the fuckwhat the fuckwhat the fuckwhat the fuckwhat the fuckwhat the fuckwhat the fuckwhat the fuckwhat the fuckwhat the fuckwhat the fuckwhat the fuckwhat the fuckwhat the fuckwhat the the fuckwhat the fuckwhat the fuckwhat the fuck"
	fmt.Printf("源数据：%v\n 长度=%d\n", []byte(str), len(str))
	bt, err := rsa.DefaultEncrypt([]byte(str))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("加密数据：%v\n", bt)
	ss, err := rsa.DefaultDecrypt(bt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("加密解密后数据：%v\n", ss)
	fmt.Printf("源数据长度=%d 加密后长度=%d\n", len(ss), len(bt))
}
