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
	d := GetDate("1996-09-01")
	for i := 0; i < 24; i++ {
		fmt.Println(d.NewEightWord(i).EWString())
	}

	//d2 := GetDate("1996-12-25")
	//for i := 0; i < 24; i += 2 {
	//	fmt.Println(d2.LunarMonth+"月"+d2.LunarDay, d2.Animal, d2.Constellation(), d2.EightWords(i))
	//}
}
