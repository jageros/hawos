/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    email
 * @Date:    2021/8/23 2:30 下午
 * @package: email
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package email

import (
	"context"
	"fmt"
	"github.com/jageros/hawox/contextx"
	"github.com/jageros/hawox/logx"
	"github.com/jordan-wright/email"
	"mime"
	"net/smtp"
	"time"
)

var opt_ *Option
var pool *email.Pool

//var mailCh chan *email.Email

type Option struct {
	g           contextx.Context
	Identify    string
	Address     string
	ConnMax     int
	Username    string
	Password    string
	Host        string
	SendTimeout time.Duration
}

func defaultOption() *Option {
	return &Option{
		Address:     "smtp.163.com:25",
		ConnMax:     2,
		Username:    "lhj168os@163.com",
		Password:    "NHMAIBZAOMGXYPHB",
		Host:        "smtp.163.com",
		SendTimeout: time.Second * 10,
	}
}

func Initialize(ctx contextx.Context, opfs ...func(opt *Option)) error {
	opt := defaultOption()
	for _, opf := range opfs {
		opf(opt)
	}
	opt.g = ctx
	opt_ = opt
	var err error
	pool, err = email.NewPool(opt.Address, opt.ConnMax, smtp.PlainAuth(opt.Identify, opt.Username, opt.Password, opt.Host))
	if err != nil {
		return err
	}

	ctx.Go(func(ctx context.Context) error {
		<-ctx.Done()
		pool.Close()
		return ctx.Err()
	})

	return nil
}

type Email struct {
	To         []string
	From       string
	Title      string
	Content    []byte
	Attachment string
}

func (e *Email) mail() *email.Email {
	m := &email.Email{
		To:      e.To,
		From:    fmt.Sprintf("%s<%s>", mime.QEncoding.Encode("UTF-8", e.From), opt_.Username),
		Subject: e.Title,
		HTML:    e.Content,
	}
	if e.Attachment != "" {
		m.AttachFile(e.Attachment)
	}
	return m
}

func (e *Email) Send() {
	opt_.g.Go(func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := pool.Send(e.mail(), opt_.SendTimeout)
			if err != nil {
				logx.Err(err).Msg("Send Email")
			}
		}
		return nil
	})
}
