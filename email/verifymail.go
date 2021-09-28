/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    verifymail
 * @Date:    2021/8/23 4:10 下午
 * @package: email
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package email

import (
	"fmt"
	"github.com/jageros/hawox/contextx"
	"github.com/jordan-wright/email"
	"mime"
	"time"
)

var content = `<p>尊敬的%s, 您好！非常感谢您注册%s会员!</br>

请点击下面链接验证您的邮箱地址，有效时间%s。 </br>

<a href=%s>%s</a></br>

如果链接无法点击，请尝试拷贝该地址到您的浏览器地址栏进行验证。</br>

© 2021 %s </p>
`

type VerifyMail struct {
	Name      string
	VerifyUrl string
	ValidTime time.Duration
	Email
}

func (v *VerifyMail) content() []byte {
	if len(v.Content) > 0 {
		return v.Content
	}
	txt := fmt.Sprintf(content, v.Name, v.From, v.ValidTime.String(), v.VerifyUrl, v.VerifyUrl, v.From)
	v.Content = []byte(txt)
	return v.Content
}

func (v *VerifyMail) mail() *email.Email {
	m := &email.Email{
		To:      v.To,
		From:    fmt.Sprintf("%s<%s>", mime.QEncoding.Encode("UTF-8", v.From), opt_.Username),
		Subject: v.Title,
		HTML:    v.content(),
	}
	return m
}

func (v *VerifyMail) Send() {
	opt_.g.Go(func(ctx contextx.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			pool.Send(v.mail(), opt_.SendTimeout)
		}
		return nil
	})
}
