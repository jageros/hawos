/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    tianapi_test
 * @Date:    2021/11/30 15:11
 * @package: tianapi
 * @Version: x.x.x
 *
 * @Description: xxx
 *
 */

package tianapi

import (
	"fmt"
	"testing"
)

func Test_GetDateType(t *testing.T) {
	ty, err := CheckDateType("2022-02-05")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(ty.String(), ty.Info())
	}
}
