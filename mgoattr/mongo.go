package mgoattr

import (
	"context"
	"github.com/jageros/hawox/contextx"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"log"
	"time"
)

type attrData struct {
	attrID interface{}
	data   map[string]interface{}
}

func (ad *attrData) GetAttrID() interface{} {
	return ad.attrID
}

func (ad *attrData) GetData() map[string]interface{} {
	return ad.data
}

type mongoEngine struct {
	session  *mgo.Session
	database string
	ctx      contextx.Context
	cancel   contextx.CancelFunc
}

func openMongoDB(addr, dbname, user, passowrd string) (*mongoEngine, error) {
	log.Printf("Connecting MongoDB %s ...", addr)
	session, err := mgo.Dial("mongodb://" + addr + "/")
	if err != nil {
		return nil, err
	}

	db := session.DB(dbname)
	if user != "" {
		if err = db.Login(user, passowrd); err != nil {
			return nil, err
		}
	}

	session.SetMode(mgo.Strong, true)

	ctx, cancel := contextx.WithCancel(context.Background())

	ctx.Go(func(ctx context.Context) error {
		tk := time.NewTicker(time.Second * 10)
		var errNum int
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-tk.C:
				err := session.Ping()
				if err != nil {
					errNum++
					if errNum >= 10 {
						return err
					}
				} else {
					errNum = 0
				}
			}
		}
	})

	return &mongoEngine{
		session:  session,
		database: dbname,
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

func (e *mongoEngine) write(attrName string, attrID interface{}, data map[string]interface{}) error {
	col := e.getCollection(attrName)
	_, err := col.UpsertId(attrID, bson.M{
		"data": data,
	})
	col.Insert()
	col.Database.Session.Close()
	return err
}

func (e *mongoEngine) insert(attrName string, attrID interface{}, data map[string]interface{}) error {
	col := e.getCollection(attrName)
	err := col.Insert(bson.M{"_id": attrID, "data": data})
	col.Database.Session.Close()
	return err
}

func (e *mongoEngine) query(attrName string) (func() (attrID interface{}, data map[string]interface{}, hasMore bool), error) {
	col := e.session.DB(e.database).C(attrName)
	iter := col.Find(nil).Iter()
	return func() (attrID interface{}, data map[string]interface{}, hasMore bool) {
		var doc bson.M
		if iter.Next(&doc) {
			attrID = doc["_id"]
			data = e.convertM2Map(doc["data"].(bson.M))
			hasMore = true
		}
		return
	}, nil
}

func (e *mongoEngine) read(attrName string, attrID interface{}) (map[string]interface{}, error) {
	col := e.getCollection(attrName)
	q := col.FindId(attrID)
	var doc bson.M
	err := q.One(&doc)
	col.Database.Session.Close()

	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return e.convertM2Map(doc["data"].(bson.M)), nil
}

func (e *mongoEngine) convertM2Map(m bson.M) map[string]interface{} {
	ma := map[string]interface{}(m)
	e.convertM2MapInMap(ma)
	return ma
}

func (e *mongoEngine) convertM2MapInMap(m map[string]interface{}) {
	for k, v := range m {
		switch im := v.(type) {
		case bson.M:
			m[k] = e.convertM2Map(im)
		case map[string]interface{}:
			e.convertM2MapInMap(im)
		case []interface{}:
			e.convertM2MapInList(im)
		}
	}
}

func (e *mongoEngine) convertM2MapInList(l []interface{}) {
	for i, v := range l {
		switch im := v.(type) {
		case bson.M:
			l[i] = e.convertM2Map(im)
		case map[string]interface{}:
			e.convertM2MapInMap(im)
		case []interface{}:
			e.convertM2MapInList(im)
		}
	}
}

func (e *mongoEngine) getCollection(attrName string) *mgo.Collection {
	ses := e.session.Copy()
	return ses.DB(e.database).C(attrName)
}

func (e *mongoEngine) exists(attrName string, attrID interface{}) (bool, error) {
	col := e.getCollection(attrName)
	query := col.FindId(attrID)
	var doc bson.M
	err := query.One(&doc)
	col.Database.Session.Close()

	if err == nil {
		// doc found
		return true, nil
	} else if err == mgo.ErrNotFound {
		return false, nil
	} else {
		return false, err
	}
}

func (e *mongoEngine) close() {
	e.cancel()
	e.session.Close()
	e.ctx.Wait()
}

func (e *mongoEngine) isEOF(err error) bool {
	return err == io.EOF || err == io.ErrUnexpectedEOF
}

func (e *mongoEngine) del(attrName string, attrID interface{}) error {
	col := e.getCollection(attrName)
	err := col.RemoveId(attrID)
	col.Database.Session.Close()
	return err
}

func (e *mongoEngine) readAll(attrName string) ([]interface {
	GetAttrID() interface{}
	GetData() map[string]interface{}
}, error) {

	col := e.getCollection(attrName)
	q := col.Find(bson.M{})
	var docs []bson.M
	err := q.All(&docs)
	col.Database.Session.Close()

	if err != nil {
		if err == mgo.ErrNotFound {
			return []interface {
				GetAttrID() interface{}
				GetData() map[string]interface{}
			}{}, nil
		}
		return nil, err
	}

	var datas []interface {
		GetAttrID() interface{}
		GetData() map[string]interface{}
	}
	for _, doc := range docs {
		datas = append(datas, &attrData{
			attrID: doc["_id"],
			data:   e.convertM2Map(doc["data"].(bson.M)),
		})
	}

	return datas, nil
}
