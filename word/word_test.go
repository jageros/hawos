/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    word_test
 * @Date:    2021/12/17 6:14 下午
 * @package: word
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package word

import (
	"fmt"
	"testing"
)

func Test_CheckWord(t *testing.T) {
	ds := CheckWord("戈")
	fmt.Println(ds[0].GetRadicals())
}
