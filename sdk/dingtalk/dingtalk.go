package dingtalk

import (
	"fmt"
	"github.com/jageros/hawox/encrypt"
	"github.com/jageros/hawox/errcode"
	"github.com/jageros/hawox/httpc"
	"time"
)

var (
	secret_   = "SECff53e03a11a5754e3773d7d090c6ac999802c9d559ea743ee93e1ba21daac89e"
	robotUrl_ = "https://oapi.dingtalk.com/robot/send?access_token=e0c7fff08fb219ab9d8c626c89cf9db6cad36df1c233bab3208561328e3259d3"
)

func UpdateConfig(secret, robotUrl string) {
	secret_ = secret
	robotUrl_ = robotUrl
}

type result struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func SendMsg(text string) error {
	timestamp := time.Now().UnixMilli()
	str := fmt.Sprintf("%d\n%s", timestamp, secret_)
	sign := encrypt.HmacSha256Base64(str, secret_)

	url := fmt.Sprintf("%s&timestamp=%d&sign=%s", robotUrl_, timestamp, sign)
	arg := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content": text,
		},
	}
	resp := new(result)
	err := httpc.Request2(httpc.POST, url, httpc.JSON, arg, nil, resp)
	if err != nil {
		return err
	}
	if resp.Errcode != 0 {
		return errcode.New(int32(resp.Errcode), resp.Errmsg)
	}
	return nil
}
