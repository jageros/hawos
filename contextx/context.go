/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    context
 * @Date:    2021/9/9 4:53 下午
 * @package: contextx
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package contextx

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Context interface {
	context.Context
	Go(func(ctx Context) error)
	Wait() error
	CancelSub()
	WithCancel() (Context, CancelFunc)
	WithTimeout(timeout time.Duration) (Context, CancelFunc)
	WithDeadline(d time.Time) (Context, CancelFunc)
	WithValue(key, val interface{}) (Context, CancelFunc)
}

type CancelFunc context.CancelFunc

func WithSignal(sig ...os.Signal) (Context, CancelFunc) {
	ctx, cancel := signal.NotifyContext(context.Background(), sig...)
	return newGroup(ctx, CancelFunc(cancel))
}

func Default() (Context, CancelFunc) {
	return WithSignal(syscall.SIGINT, syscall.SIGTERM)
}

func Background() Context {
	ctx, cancel := context.WithCancel(context.Background())
	c, _ := newGroup(ctx, CancelFunc(cancel))
	return c
}
func TODO() Context {
	ctx, cancel := context.WithCancel(context.TODO())
	c, _ := newGroup(ctx, CancelFunc(cancel))
	return c
}
