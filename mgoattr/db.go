package mgoattr

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jageros/hawox/contextx"
	"github.com/xiaonanln/go-xnsyncutil/xnsyncutil"
)

var (
	clients  []*dbClient
	initOnce sync.Once
)

type dbClient struct {
	opt                  *Option
	dbEngine             *mongoEngine
	operationQueue       *xnsyncutil.SyncQueue
	recentWarnedQueueLen int
	shutdownOnce         sync.Once
	shutdownNotify       chan struct{}
}

func initDbs(ctx contextx.Context) {
	initOnce.Do(func() {
		ctx.Go(func(ctx context.Context) error {
			<-ctx.Done()
			for _, c := range clients {
				c.shutdown()
			}
			clients = []*dbClient{}
			return nil
		})
	})
}

func getOrNewDbClient(opt *Option) *dbClient {
	for _, cli := range clients {
		if cli.opt.Addr == opt.Addr && cli.opt.DB == opt.DB {
			return cli
		}
	}

	cli := &dbClient{
		opt:            opt,
		operationQueue: xnsyncutil.NewSyncQueue(),
		shutdownNotify: make(chan struct{}),
		shutdownOnce:   sync.Once{},
	}

	err := cli.assureDBEngineReady()
	if err != nil {
		log.Fatalf("db engine %s is not ready: %s", opt, err)
	}

	clients = append(clients, cli)
	go cli.dbRoutine()
	return cli
}

func (c *dbClient) assureDBEngineReady() (err error) {
	if c.dbEngine != nil {
		return
	}
	c.dbEngine, err = openMongoDB(c.opt.Addr, c.opt.DB, c.opt.User, c.opt.Password)
	return
}

func (c *dbClient) insert(attrName string, attrID interface{}, data map[string]interface{}) error {
	req := &insertRequest{
		attrName: attrName,
		attrID:   attrID,
		data:     data,
		c:        make(chan error, 1),
	}

	c.operationQueue.Push(req)
	c.checkOperationQueueLen()

	err := <-req.c
	return err
}

func (c *dbClient) save(attrName string, attrID interface{}, data map[string]interface{}, needReply bool) error {
	req := &saveRequest{
		attrName: attrName,
		attrID:   attrID,
		data:     data,
	}
	if needReply {
		req.c = make(chan error, 1)
	}

	c.operationQueue.Push(req)
	c.checkOperationQueueLen()

	if needReply {
		return <-req.c
	} else {
		return nil
	}
}

func (c *dbClient) del(attrName string, attrID interface{}, needReply bool) error {
	req := &delRequest{
		attrName: attrName,
		attrID:   attrID,
	}
	if needReply {
		req.c = make(chan error, 1)
	}

	c.operationQueue.Push(req)
	c.checkOperationQueueLen()

	if needReply {
		return <-req.c
	} else {
		return nil
	}
}

func (c *dbClient) load(attrName string, attrID interface{}) (map[string]interface{}, error) {
	req := &loadRequest{
		attrName: attrName,
		attrID:   attrID,
		c:        make(chan *loadResult, 1),
	}

	c.operationQueue.Push(req)
	c.checkOperationQueueLen()

	result := <-req.c
	return result.data, result.err
}

func (c *dbClient) exists(attrName string, attrID interface{}) (bool, error) {
	req := &existsRequest{
		attrName: attrName,
		attrID:   attrID,
		c:        make(chan *existsResult, 1),
	}

	c.operationQueue.Push(req)
	c.checkOperationQueueLen()

	result := <-req.c
	return result.exists, result.err
}

func (c *dbClient) loadAll(attrName string) ([]interface {
	GetAttrID() interface{}
	GetData() map[string]interface{}
}, error) {

	req := &loadAllRequest{
		attrName: attrName,
		c:        make(chan *loadAllResult, 1),
	}

	c.operationQueue.Push(req)
	c.checkOperationQueueLen()

	result := <-req.c
	return result.datas, result.err
}

func (c *dbClient) forEach(attrName string, callback func(attrID interface{}, data map[string]interface{})) {
	req := &forEachRequest{attrName: attrName, iter: nil, c: make(chan *forEachResult, 1)}

	for true {
		c.operationQueue.Push(req)
		c.checkOperationQueueLen()

		result := <-req.c
		if result.err != nil {
			return
		}

		if !result.hasMore {
			return
		}
		callback(result.attrID, result.data)
	}
}

func (c *dbClient) checkOperationQueueLen() {
	qlen := c.operationQueue.Len()
	if qlen > 100 && qlen%100 == 0 && c.recentWarnedQueueLen != qlen {
		log.Printf("db %s operation queue length = %d", c.opt, qlen)
		c.recentWarnedQueueLen = qlen
	}
}

func (c *dbClient) shutdown() {
	c.shutdownOnce.Do(func() {
		var waitTime time.Duration
		for c.operationQueue.Len() > 0 {
			if waitTime > 10*time.Second {
				log.Printf("db %s Shutdown timeout, left op %d", c.opt.format(), c.operationQueue.Len())
				break
			}
			t := 100 * time.Millisecond
			waitTime += t
			time.Sleep(t)
		}

		c.operationQueue.Close()
		<-c.shutdownNotify
	})
}

func (c *dbClient) dbRoutine() {
	defer func() {
		err := recover()
		if err != nil {
			log.Printf("db %s routine paniced: %s", c.opt.format(), err)
		} else {
			c.dbEngine.close()
			c.dbEngine = nil
			close(c.shutdownNotify)
		}
	}()

	for {
		err := c.assureDBEngineReady()
		if err != nil {
			log.Printf("db %s engine is not ready: %s", c.opt.format(), err)
			time.Sleep(time.Second)
			continue
		}

		if c.dbEngine == nil {
			log.Fatalf("db %s engine is nil", c.opt.format())
		}

		req := c.operationQueue.Pop()
		if req == nil {
			break
		}

		req2, ok := req.(iDbRequest)
		if !ok {
			log.Printf("db: unknown operation: %v", req)
			continue
		}

		op := startOperation(fmt.Sprintf("db:%s", req2.name()))

		err = req2.execute(c.dbEngine)
		if err != nil {
			log.Printf("db: %s %s failed: %s", c.opt.format(), req2.name(), err)

			if err != nil && c.dbEngine.isEOF(err) {
				c.dbEngine.close()
				c.dbEngine = nil
			}
		}

		op.finish(100 * time.Millisecond)
	}
}

type iDbRequest interface {
	name() string
	execute(engine *mongoEngine) error
}

type saveRequest struct {
	attrName string
	attrID   interface{}
	data     map[string]interface{}
	c        chan error
}

func (r *saveRequest) name() string {
	return "save"
}

func (r *saveRequest) execute(engine *mongoEngine) error {
	err := engine.write(r.attrName, r.attrID, r.data)
	if r.c != nil {
		r.c <- err
	}
	return err
}

type delRequest struct {
	attrName string
	attrID   interface{}
	c        chan error
}

func (r *delRequest) name() string {
	return "del"
}

func (r *delRequest) execute(engine *mongoEngine) error {
	err := engine.del(r.attrName, r.attrID)
	if r.c != nil {
		r.c <- err
	}
	return err
}

type loadRequest struct {
	attrName string
	attrID   interface{}
	c        chan *loadResult
}

type loadResult struct {
	data map[string]interface{}
	err  error
}

func (r *loadRequest) name() string {
	return "load"
}

func (r *loadRequest) execute(engine *mongoEngine) error {
	data, err := engine.read(r.attrName, r.attrID)
	if err != nil {
		data = nil
	}

	if r.c != nil {
		r.c <- &loadResult{
			data: data,
			err:  err,
		}
	}
	return err
}

type existsRequest struct {
	attrName string
	attrID   interface{}
	c        chan *existsResult
}

type existsResult struct {
	exists bool
	err    error
}

func (r *existsRequest) name() string {
	return "exists"
}

func (r *existsRequest) execute(engine *mongoEngine) error {
	exists, err := engine.exists(r.attrName, r.attrID)
	if r.c != nil {
		r.c <- &existsResult{
			exists: exists,
			err:    err,
		}
	}
	return err
}

type loadAllRequest struct {
	attrName string
	c        chan *loadAllResult
}

type loadAllResult struct {
	datas []interface {
		GetAttrID() interface{}
		GetData() map[string]interface{}
	}
	err error
}

func (r *loadAllRequest) name() string {
	return "loadAll"
}

func (r *loadAllRequest) execute(engine *mongoEngine) error {
	datas, err := engine.readAll(r.attrName)
	if err != nil {
		datas = nil
	}

	if r.c != nil {
		r.c <- &loadAllResult{
			datas: datas,
			err:   err,
		}
	}
	return err
}

type forEachRequest struct {
	attrName string
	iter     func() (attrID interface{}, data map[string]interface{}, hasMore bool)
	c        chan *forEachResult
}

type forEachResult struct {
	attrID  interface{}
	data    map[string]interface{}
	hasMore bool
	err     error
}

func (r *forEachRequest) name() string {
	return "forEach"
}

func (r *forEachRequest) execute(engine *mongoEngine) error {
	var err error
	var attrID interface{}
	var data map[string]interface{}
	var hasMore bool
	if r.iter != nil {
		attrID, data, hasMore = r.iter()
	} else {
		r.iter, err = engine.query(r.attrName)
		if err == nil {
			attrID, data, hasMore = r.iter()
		}
	}

	r.c <- &forEachResult{
		attrID:  attrID,
		data:    data,
		hasMore: hasMore,
		err:     err,
	}
	return err
}

type insertRequest struct {
	attrName string
	attrID   interface{}
	data     map[string]interface{}
	c        chan error
}

func (r *insertRequest) name() string {
	return "insert"
}

func (r *insertRequest) execute(engine *mongoEngine) error {
	err := engine.insert(r.attrName, r.attrID, r.data)
	if r.c != nil {
		r.c <- err
	}
	return err
}

// ======================

type operation struct {
	name      string
	startTime time.Time
}

var operationAllocPool = sync.Pool{
	New: func() interface{} {
		return &operation{}
	},
}

func startOperation(operationName string) *operation {
	op := operationAllocPool.Get().(*operation)
	op.name = operationName
	op.startTime = time.Now()
	return op
}

func (op *operation) finish(warnThreshold time.Duration) {
	takeTime := time.Now().Sub(op.startTime)
	if warnThreshold > 0 && takeTime >= warnThreshold {
		log.Printf("opmon: operation %s takes %s > %s", op.name, takeTime, warnThreshold)
	}
	operationAllocPool.Put(op)
}
