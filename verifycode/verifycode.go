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

package verifycode

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

var captcha *base64Captcha.Captcha

func Initialize(opfs ...func(option *Option)) {
	opt := defaultOption()

	for _, opf := range opfs {
		opf(opt)
	}

	if opt.EnableRedis {
		captcha = base64Captcha.NewCaptcha(base64Captcha.NewDriverDigit(opt.Height, opt.Width, opt.Length, opt.MaxSkew, opt.DotCount), &redisStore{expireTime: opt.ExpireTime})
	} else {
		captcha = base64Captcha.NewCaptcha(base64Captcha.NewDriverDigit(opt.Height, opt.Width, opt.Length, opt.MaxSkew, opt.DotCount), base64Captcha.DefaultMemStore)
	}
}

func getCaptcha() *base64Captcha.Captcha {
	if captcha == nil {
		Initialize()
	}
	return captcha
}

func GenCodeAuto() (id, value, b64s string, err error) {
	c := getCaptcha()
	id, b64s, err = c.Generate()
	if err != nil {
		return
	}
	value = c.Store.Get(id, false)
	return
}

func GenCode(id, value string) (b64s string, err error) {
	c := getCaptcha()
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

func Verify(id, answer string) bool {
	c := getCaptcha()
	return c.Verify(id, answer, true)
}
