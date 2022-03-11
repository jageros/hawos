/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    contextx
 * @Date:    2021/9/9 4:36 下午
 * @package: contextx
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package contextx

import (
	"context"
	"sync"
	"time"
)

type groupCtx struct {
	ctx    context.Context
	cancel CancelFunc

	wg      sync.WaitGroup
	errOnce sync.Once
	err     error
}

func newGroup(ctx context.Context, cancel CancelFunc) (*groupCtx, CancelFunc) {
	return &groupCtx{
		ctx:    ctx,
		cancel: cancel,
	}, cancel
}

func (c *groupCtx) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *groupCtx) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *groupCtx) Err() error {
	return c.ctx.Err()
}

func (c *groupCtx) Value(key interface{}) interface{} {
	return c.ctx.Value(key)
}

func (c *groupCtx) Go(f func(ctx context.Context) error) {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		if err := f(c); err != nil {
			c.errOnce.Do(func() {
				c.err = err
				if c.cancel != nil {
					c.cancel()
				}
			})
		}
	}()
}

func (c *groupCtx) Wait() error {
	c.wg.Wait()
	if c.cancel != nil {
		c.cancel()
	}
	return c.err
}
