/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    filter_test
 * @Date:    2022/2/21 17:50
 * @package: sensitive
 * @Version: x.x.x
 *
 * @Description: xxx
 *
 */

package sensitive

import (
	"fmt"
	"testing"
)

func TestFindAll(t *testing.T) {
	f := New()
	err := f.LoadWordDict("../config/dict.txt")
	if err != nil {
		t.Error(err)
		return
	}
	ss := f.FindAll("fuck what the fuck fuck gun fuck!")
	fmt.Println(ss)
}
