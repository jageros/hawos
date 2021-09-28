/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    group
 * @Date:    2021/8/18 5:35 下午
 * @package: group
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package group

import (
	"context"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var g_ *Group

type Group struct {
	eg     *errgroup.Group
	ctx    context.Context
	cancel context.CancelFunc
}

func New() *Group {
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx2 := errgroup.WithContext(ctx)
	gp := &Group{
		eg:     eg,
		ctx:    ctx2,
		cancel: cancel,
	}
	return gp
}

func NewWithContext(ctx context.Context) *Group {
	ctx_, cancel := context.WithCancel(ctx)
	eg, ctx2 := errgroup.WithContext(ctx_)
	gp := &Group{
		eg:     eg,
		ctx:    ctx2,
		cancel: cancel,
	}
	return gp
}

func NewWithSignal(sig ...os.Signal) *Group {
	ctx, cancel := signal.NotifyContext(context.Background(), sig...)
	eg, ctx2 := errgroup.WithContext(ctx)
	gp := &Group{
		eg:     eg,
		ctx:    ctx2,
		cancel: cancel,
	}
	return gp
}

func Default() *Group {
	return NewWithSignal(syscall.SIGINT, syscall.SIGTERM)
}

func (g *Group) SubGroup() *Group {
	ctx, cancel := g.WithCancel()
	eg, ctx2 := errgroup.WithContext(ctx)
	gp := &Group{
		eg:     eg,
		ctx:    ctx2,
		cancel: cancel,
	}
	return gp
}

func (g *Group) Go(f func(ctx context.Context) error) {
	g.eg.Go(func() error {
		return f(g.ctx)
	})
}

func (g *Group) Wait() error {
	return g.eg.Wait()
}

func (g *Group) Cancel() {
	g.cancel()
}

func (g *Group) Context() context.Context {
	return g.ctx
}

func (g *Group) WithCancel() (context.Context, context.CancelFunc) {
	return context.WithCancel(g.ctx)
}

func (g *Group) WithTimeout(t time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(g.ctx, t)
}

func (g *Group) WithDeadline(t time.Time) (context.Context, context.CancelFunc) {
	return context.WithDeadline(g.ctx, t)
}

func (g *Group) WithValue(key, val interface{}) context.Context {
	return context.WithValue(g.ctx, key, val)
}

// ====== context interface ======

func (g *Group) Deadline() (deadline time.Time, ok bool) {
	return
}

func (g *Group) Done() <-chan struct{} {
	return g.ctx.Done()
}

func (g *Group) Err() error {
	return g.ctx.Err()
}

func (g *Group) Value(key interface{}) interface{} {
	return g.ctx.Value(key)
}

// ============ Global API ============

func init() {
	g_ = Default()
}

func Go(f func(ctx context.Context) error) {
	g_.Go(f)
}

func Wait() error {
	return g_.eg.Wait()
}

func Cancel() {
	g_.cancel()
}
