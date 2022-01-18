/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    udpx
 * @Date:    2022/1/18 5:06 下午
 * @package: udpx
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package udpx

import (
	"fmt"
	"github.com/jageros/hawox/contextx"
	"net"
	"os"
	"strings"
)

func Init(ctx contextx.Context, ops ...func(opt *Config)) error {
	s := defaultServer()
	for _, op := range ops {
		op(s)
	}
	addr, err := net.ResolveUDPAddr("udp", s.addr())
	if err != nil {
		return err
	}
	addr.String()
	s.conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	return nil
}

func ssss() {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer conn.Close()

	for {
		// Here must use make and give the lenth of buffer
		data := make([]byte, config.SERVER_RECV_LEN)
		_, rAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Println(err)
			continue
		}

		strData := string(data)
		fmt.Println("Received:", strData)

		upper := strings.ToUpper(strData)
		_, err = conn.WriteToUDP([]byte(upper), rAddr)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("Send:", upper)
	}
}
