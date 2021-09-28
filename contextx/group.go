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

	subCancel []CancelFunc

	wg      sync.WaitGroup
	errOnce sync.Once
	err     error
	parent  *groupCtx
}

func newGroup(ctx context.Context, cancel CancelFunc) (*groupCtx, CancelFunc) {
	return &groupCtx{
		ctx:    ctx,
		cancel: cancel,
	}, cancel
}

func (gCtx *groupCtx) subGroup(ctx context.Context, cancel CancelFunc) (Context, CancelFunc) {
	gCtx.subCancel = append(gCtx.subCancel, cancel)
	return &groupCtx{
		ctx:    ctx,
		cancel: cancel,
		parent: gCtx,
	}, cancel
}

func (gCtx *groupCtx) add(delta int) {
	if gCtx.parent != nil {
		gCtx.parent.add(delta)
	}
	gCtx.wg.Add(delta)
}

func (gCtx *groupCtx) done() {
	if gCtx.parent != nil {
		gCtx.parent.done()
	}
	gCtx.wg.Done()
}

func (gCtx *groupCtx) Deadline() (deadline time.Time, ok bool) {
	return gCtx.ctx.Deadline()
}

func (gCtx *groupCtx) Done() <-chan struct{} {
	return gCtx.ctx.Done()
}

func (gCtx *groupCtx) Err() error {
	return gCtx.ctx.Err()
}

func (gCtx *groupCtx) Value(key interface{}) interface{} {
	return gCtx.ctx.Value(key)
}

func (gCtx *groupCtx) Go(f func(ctx Context) error) {
	gCtx.add(1)

	go func() {
		defer gCtx.done()

		if err := f(gCtx); err != nil {
			gCtx.errOnce.Do(func() {
				gCtx.err = err
				if gCtx.cancel != nil {
					gCtx.cancel()
				}
			})
		}
	}()
}

func (gCtx *groupCtx) Wait() error {
	gCtx.wg.Wait()
	if gCtx.cancel != nil {
		gCtx.cancel()
	}
	return gCtx.err
}

func (gCtx *groupCtx) CancelSub() {
	for _, cancel := range gCtx.subCancel {
		cancel()
	}
}

// ================ sub ctx ================

func (gCtx *groupCtx) WithCancel() (Context, CancelFunc) {
	ctx, cancel := context.WithCancel(gCtx)
	return gCtx.subGroup(ctx, CancelFunc(cancel))
}
func (gCtx *groupCtx) WithTimeout(timeout time.Duration) (Context, CancelFunc) {
	ctx, cancel := context.WithTimeout(gCtx, timeout)
	return gCtx.subGroup(ctx, CancelFunc(cancel))
}
func (gCtx *groupCtx) WithDeadline(d time.Time) (Context, CancelFunc) {
	ctx, cancel := context.WithDeadline(gCtx, d)
	return gCtx.subGroup(ctx, CancelFunc(cancel))
}
func (gCtx *groupCtx) WithValue(key, val interface{}) (Context, CancelFunc) {
	ctx, cancel := context.WithCancel(gCtx)
	ctx = context.WithValue(ctx, key, val)
	return gCtx.subGroup(ctx, CancelFunc(cancel))
}
