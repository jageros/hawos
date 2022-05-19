/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    lnode_test
 * @Date:    2022/3/30 16:53
 * @package: lnode
 * @Version: x.x.x
 *
 * @Description: xxx
 *
 */

package lnode

import (
	"fmt"
	"testing"
)

func TestLNode_Reverse(t *testing.T) {
	l := &Node{}
	for i := 0; i < 10; i++ {
		l.Push(i)
		l.Push(i)
	}
	fmt.Println(l.ToArray())
	l.RemoveDup()
	fmt.Println(l.ToArray())
	l.Reverse()
	fmt.Println(l.ToArray())
}
