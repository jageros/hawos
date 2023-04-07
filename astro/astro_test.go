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
	d1 := GetDate("1993-08-27")
	fmt.Println(d1.NewEightWord(3).EWString())
	d := GetDate("1997-12-01")
	for i := 0; i <= 23; i += 2 {
		fmt.Println(d.NewEightWord(i).EWString())
	}

	//fmt.Println(d.Suitable, d.Avoid)

	//d2 := GetDate("1996-09-01")
	//fmt.Println(d2.NewEightWord(15).EWString())
}
