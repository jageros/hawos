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
	dateMap, err := ReadFromJsonFile("date.json")
	if err != nil {
		t.Error(err)
		return
	} else {
		fmt.Println("=== init successful ===")
		fmt.Printf("DayCount=%d\n", len(dateMap))
		fmt.Println(dateMap["1993-08-27"].EightWord(3))
	}

	//f, err := os.OpenFile("date1.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//defer f.Close()
	//
	//for y := 2021; y <= 2021; y++ {
	//	fmt.Printf("==== %d ====\n", y)
	//	for m := 1; m <= 1; m++ {
	//		ds, err := Crawling(y, m)
	//		if err != nil {
	//			t.Error(err)
	//			return
	//		}
	//		bty, err := json.Marshal(ds)
	//		if err != nil {
	//			t.Error(err)
	//			return
	//		}
	//		_, err = f.WriteString(string(bty) + ",\n")
	//		//_, err = fmt.Fprintln(buf, string(bty)+",")
	//		if err != nil {
	//			t.Error(err)
	//			return
	//		}
	//	}
	//	//time.Sleep(time.Second * 10)
	//}

	//ds, err := Crawling(1904, 12)
	//if err != nil {
	//	t.Error(err)
	//}
	//for _, d := range ds {
	//	bty, err := json.Marshal(d)
	//	if err != nil {
	//		t.Error(err)
	//	}
	//	fmt.Println(string(bty) + ",")
	//}
}
