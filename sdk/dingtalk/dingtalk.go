package dingtalk

import (
	"fmt"
	"github.com/jageros/hawox/encrypt"
	"github.com/jageros/hawox/errcode"
	"github.com/jageros/hawox/httpc"
	"time"
)

type Option struct {
	Secret   string
	RobotUrl string
}

type DingTalk struct {
	opt *Option
}

func NewDingTalk(opfs ...func(opt *Option)) *DingTalk {
	opt_ := &Option{}
	for _, opf := range opfs {
		opf(opt_)
	}
	return &DingTalk{opt: opt_}
}

type result struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func (d *DingTalk) SendMsg(text string) error {
	timestamp := time.Now().UnixMilli()
	str := fmt.Sprintf("%d\n%s", timestamp, d.opt.Secret)
	sign := encrypt.HmacSha256Base64(str, d.opt.Secret)

	url := fmt.Sprintf("%s&timestamp=%d&sign=%s", d.opt.RobotUrl, timestamp, sign)
	arg := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content": text,
		},
	}
	resp := new(result)
	err := httpc.RequestWithInterface(httpc.POST, url, httpc.JSON, arg, nil, resp)
	if err != nil {
		return err
	}
	if resp.Errcode != 0 {
		return errcode.New(int32(resp.Errcode), resp.Errmsg)
	}
	return nil
}
