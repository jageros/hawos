/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    astro_test
 * @Date:    2021/11/30 15:50
 * @package: astro
 * @Version: x.x.x
 *
 * @Description: xxx
 *
 */

package astro

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_Lunar(t *testing.T) {
	ds, err := Crawling(2022, 1)
	if err != nil {
		t.Error(err)
	}
	for _, d := range ds {
		bty, err := json.Marshal(d)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(string(bty))
	}
}
