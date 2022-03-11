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
)

type Context interface {
	context.Context
	Go(func(ctx context.Context) error)
	Wait() error
}

type CancelFunc context.CancelFunc

func Default() (Context, CancelFunc) {
	return WithSignal(syscall.SIGINT, syscall.SIGTERM)
}

func WithSignal(sig ...os.Signal) (Context, CancelFunc) {
	ctx, cancel := signal.NotifyContext(context.Background(), sig...)
	return newGroup(ctx, CancelFunc(cancel))
}

func WithCancel(parent context.Context) (Context, CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	return newGroup(ctx, CancelFunc(cancel))
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
