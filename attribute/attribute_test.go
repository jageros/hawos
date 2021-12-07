/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    attribute_test
 * @Date:    2021/12/7 2:35 下午
 * @package: attribute
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package attribute

import (
	"github.com/jageros/hawox/contextx"
	"testing"
)

var (
	ctx    contextx.Context
	cancel contextx.CancelFunc
)

func Test_Init(t *testing.T) {
	ctx, cancel = contextx.Default()
	Initialize(ctx, func(opt *Option) {
		opt.Addr = "127.0.0.1:27017"
		opt.DBName = "attribute"
	})
}

func Test_SaveValue(t *testing.T) {

	attr := NewAttrMgr("attr", "test")

	attrMap := NewMapAttr()
	attrMap.SetStr("name", "jackss")

	attr.SetMapAttr("info", attrMap)

	err := attr.Save(true)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("save successful!")
	}
}

func Test_GetValue(t *testing.T) {
	attr := NewAttrMgr("attr", "test")
	err := attr.Load(true)
	if err != nil {
		t.Error(err)
	}

	attrMap := attr.GetMapAttr("info")

	name := attrMap.GetStr("name")

	t.Logf("get name=%s", name)
}

func Test_Stop(t *testing.T) {
	cancel()
	err := ctx.Wait()
	t.Logf("stop with: %v", err)
}

