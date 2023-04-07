/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    verifycode
 * @Date:    2022/5/16 13:06
 * @package: verifycode
 * @Version: x.x.x
 *
 * @Description: xxx
 *
 */

package captcha

import (
	"github.com/mojocn/base64Captcha"
	"time"
)

const (
	// Height png height in pixel.
	_Height = 65
	// Width Captcha png width in pixel.
	_Width = 140
	// DefaultLen Default number of digits in captcha solution.
	_Length = 4
	// MaxSkew max absolute skew factor of a single digit.
	_MaxSkew = 0.8
	// DotCount Number of background circles.
	_DotCount = 30
)

type Option struct {
	// EnableRedis user redis as store.
	EnableRedis bool
	// Height png height in pixel.
	Height int
	// Width Captcha png width in pixel.
	Width int
	// DefaultLen Default number of digits in captcha solution.
	Length int
	// MaxSkew max absolute skew factor of a single digit.
	MaxSkew float64
	// DotCount Number of background circles.
	DotCount int
	// ExpireTime Redis expireTime.
	ExpireTime time.Duration
}

func defaultOption() *Option {
	return &Option{
		EnableRedis: false,
		ExpireTime:  time.Minute * 8,
		Height:      _Height,
		Width:       _Width,
		Length:      _Length,
		MaxSkew:     _MaxSkew,
		DotCount:    _DotCount,
	}
}

type Captcha struct {
	*base64Captcha.Captcha
}

func NewCaptcha(opfs ...func(option *Option)) *Captcha {
	opt := defaultOption()

	for _, opf := range opfs {
		opf(opt)
	}

	var captcha *base64Captcha.Captcha
	if opt.EnableRedis {
		captcha = base64Captcha.NewCaptcha(base64Captcha.NewDriverDigit(opt.Height, opt.Width, opt.Length, opt.MaxSkew, opt.DotCount), &redisStore{expireTime: opt.ExpireTime})
	} else {
		captcha = base64Captcha.NewCaptcha(base64Captcha.NewDriverDigit(opt.Height, opt.Width, opt.Length, opt.MaxSkew, opt.DotCount), base64Captcha.DefaultMemStore)
	}
	return &Captcha{Captcha: captcha}
}

func (c *Captcha) GenCodeAuto() (id, value, b64s string, err error) {
	id, b64s, err = c.Generate()
	if err != nil {
		return
	}
	value = c.Store.Get(id, false)
	return
}

func (c *Captcha) GenCode(id, value string) (b64s string, err error) {
	it, err := c.Driver.DrawCaptcha(value)
	if err != nil {
		return
	}
	err = c.Store.Set(id, value)
	if err != nil {
		return
	}
	b64s = it.EncodeB64string()
	return
}

func (c *Captcha) Validate(id, answer string) bool {
	ok := c.Verify(id, answer, false)
	if ok {
		c.Store.Get(id, true)
	}
	return ok
}
