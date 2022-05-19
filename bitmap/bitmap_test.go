/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    bitmap
 * @Date:    2021/11/24 4:48 下午
 * @package: bitmap
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package bitmap_test

import (
	"fmt"
	"git.hawtech.cn/jager/hawox/bitmap"
	"testing"
)

func Test_bitmap(t *testing.T) {
	bm := bitmap.New()

	for i := 23214235; i < 23214235+10000000000; i++ {
		bm.Add(i)
	}

	fmt.Printf("bitmap总长度：%d\n", bm.Len())

	if bm.Has(3847594929) {
		fmt.Println("3847594929 has find")
	}
}
