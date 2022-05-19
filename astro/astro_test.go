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
	"fmt"
	"testing"
)

func Test_Lunar(t *testing.T) {
	d := GetDate("1993-08-27")
	fmt.Println(d.NewEightWord(3).EWString())

	//d2 := GetDate("1996-09-01")
	//fmt.Println(d2.NewEightWord(15).EWString())
}
