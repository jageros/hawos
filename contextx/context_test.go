/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    context_test
 * @Date:    2022/3/11 16:00
 * @package: contextx
 * @Version: x.x.x
 *
 * @Description: xxx
 *
 */

package contextx

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestGo(t *testing.T) {
	ctx, cancel := Default()
	defer cancel()

	tm := ctx.Err()
	fmt.Println(tm)

	ctx.Go(func(ctx context.Context) error {
		tk := time.NewTimer(time.Second * 10)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tk.C:
			return errors.New("TimeoutErr")
		}
	})

	fmt.Println(ctx.Wait())
}
