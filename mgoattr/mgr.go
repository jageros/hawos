package mgoattr

import (
	"errors"
	"github.com/jageros/hawox/contextx"
)

var (
	NotExistsErr = errors.New("NotExistsErr")
)

// =========== end ============

type AttrMgr struct {
	*MapAttr
	dbCli *dbClient
	name  string
	id    interface{}
}

func NewAttrMgr(name string, id interface{}) *AttrMgr {
	return &AttrMgr{
		name:    name,
		id:      id,
		MapAttr: NewMapAttr(),
		dbCli:   getOrNewDbClient(dbConfig),
	}
}

func (a *AttrMgr) Load() error {
	if data, err := a.dbCli.load(a.name, a.id); err != nil {
		return err
	} else {
		if data == nil {
			return NotExistsErr
		}
		a.AssignMap(data)
		a.SetDirty(false)
		return nil
	}
}

func (a *AttrMgr) Copy(id interface{}) error {
	if data, err := a.dbCli.load(a.name, id); err != nil {
		return err
	} else {
		if data == nil {
			return NotExistsErr
		}
		a.AssignMap(data)
		a.SetDirty(false)
		return nil
	}
}

func (a *AttrMgr) Save(needReply bool) error {
	if a.Dirty() {
		data := a.ToMap()
		a.SetDirty(false)
		return a.dbCli.save(a.name, a.id, data, needReply)
	} else {
		return nil
	}
}

func (a *AttrMgr) Insert() error {
	data := a.ToMap()
	return a.dbCli.insert(a.name, a.id, data)
}

func (a *AttrMgr) Delete(needReply bool) error {
	return a.dbCli.del(a.name, a.id, needReply)
}

func (a *AttrMgr) Exists() (bool, error) {
	return a.dbCli.exists(a.name, a.id)
}

func (a *AttrMgr) GetAttrID() interface{} {
	return a.id
}

func LoadAll(attrName string) ([]*AttrMgr, error) {
	datas, err := getOrNewDbClient(dbConfig).loadAll(attrName)
	if err != nil {
		return nil, err
	}

	var attrs []*AttrMgr
	for _, data := range datas {
		a := NewAttrMgr(attrName, data.GetAttrID())
		a.AssignMap(data.GetData())
		a.SetDirty(false)
		attrs = append(attrs, a)
	}
	return attrs, nil
}

func ForEach(attrName string, callback func(*AttrMgr)) {
	getOrNewDbClient(dbConfig).forEach(attrName, func(attrID interface{}, data map[string]interface{}) {
		a := NewAttrMgr(attrName, attrID)
		a.AssignMap(data)
		a.SetDirty(false)
		callback(a)
	})
}

func Initialize(ctx contextx.Context, opfs ...func(opt *Option)) {
	initDbs(ctx)
	initDBConfig(opfs...)
}
