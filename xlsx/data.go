/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    data
 * @Date:    2021/9/2 3:50 下午
 * @package: xlsx
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package xlsx

import (
	"context"
	"github.com/jageros/hawox/contextx"
	"github.com/jageros/hawox/logx"
	attribute "github.com/jageros/hawox/mgoattr"
	"github.com/jageros/hawox/redis"
)

var (
	allData     = make(map[string]IData)
	allDataList []IData
	OnReload    func()
)

type IData interface {
	load() error
	reload() error
	onReload()
	Init([]byte) error
	Name() string
	AddReloadCallback(func(data IData))
}

type BaseData struct {
	I               IData
	reloadCallbacks []func(data IData)
}

func (g *BaseData) AddReloadCallback(f func(data IData)) {
	g.reloadCallbacks = append(g.reloadCallbacks, f)
}

func (g *BaseData) load() error {
	attr := attribute.NewAttrMgr("xlsxdata", g.I.Name())
	err := attr.Load()
	if err != nil {
		logx.Err(err).Str("TableName", g.I.Name()).Msg("ReadFile")
		return err
	}

	if err := g.I.Init([]byte(attr.GetStr("data"))); err != nil {
		logx.Err(err).Str("TableName", g.I.Name()).Msg("load init")
		return err
	}

	return nil
}

func (g *BaseData) reload() error {
	if err := g.load(); err != nil {
		return err
	}
	return nil
}

func (g *BaseData) onReload() {
	for _, f := range g.reloadCallbacks {
		f(g.I)
	}
}

func GetData(name string) IData {
	if d, ok := allData[name]; ok {
		return d
	} else {
		return nil
	}
}

func addData(gdata IData) {
	name := gdata.Name()
	if _, ok := allData[name]; !ok {
		allData[name] = gdata
		allDataList = append(allDataList, gdata)
	}
}

func Load(ctx contextx.Context, gdatas ...IData) {

	for _, gdata := range gdatas {
		addData(gdata)
	}

	for i := 0; i < len(allDataList); i++ {
		data := allDataList[i]
		err := data.load()
		logx.Err(err).Str("TableName", data.Name()).Msg("xlsxdata load")
	}

	ctx.Go(func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				err := redis.Subscribe(func(channel, msg string) {

					if msg != "reload" {
						return
					}

					var reloadDatas []IData
					for i := 0; i < len(allDataList); i++ {
						data := allDataList[i]
						if err := data.reload(); err != nil {
							logx.Err(err).Str("TableName", data.Name()).Msg("xlsxdata reload")
						} else {
							reloadDatas = append(reloadDatas, data)
							logx.Infof("xlsxdata %s reload ok", data.Name())
						}
					}

					for _, data := range reloadDatas {
						data.onReload()
					}

					if OnReload != nil {
						OnReload()
					}
				}, "jsondata")
				if err != nil {
					return err
				}
			}
		}
	})
}
