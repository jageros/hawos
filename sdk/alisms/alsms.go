/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    alsms
 * @Date:    2021/8/25 5:16 下午
 * @package: alisms
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package main

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"git.hawtech.cn/jager/hawox/errcode"
)

//AccessKey ID： LTAI5tBF3a38nXEhpXW4Up1A
//AccessKey Secret： YjKRC2XyMdRIkbofo2MCEHZkC40n83

var client *dysmsapi20170525.Client
var opt *Option

type Option struct {
	AccessKeyId     string
	AccessKeySecret string
	SmsSignName     string
	TemplateCode    string
	Domain          string
}

func Initialize(opfs ...func(opt *Option)) error {
	opt = &Option{
		AccessKeyId:     "LTAI5tBF3a38nXEhpXW4Up1A",
		AccessKeySecret: "YjKRC2XyMdRIkbofo2MCEHZkC40n83",
		SmsSignName:     "HawData",
		TemplateCode:    "xxxxx",
		Domain:          "dysmsapi.aliyuncs.com",
	}

	for _, opf := range opfs {
		opf(opt)
	}

	config := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: &opt.AccessKeyId,
		// 您的AccessKey Secret
		AccessKeySecret: &opt.AccessKeySecret,
	}
	// 访问的域名
	config.Endpoint = &opt.Domain
	var err error
	client, err = dysmsapi20170525.NewClient(config)
	return err
}

func SendVerifyCode(phone, code string) (err error) {
	if client == nil || opt == nil {
		return errcode.New(-1, "NotInitializeSms")
	}

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  &phone,
		TemplateParam: &code,
		SignName:      &opt.SmsSignName,
		TemplateCode:  &opt.TemplateCode,
	}
	_, err = client.SendSms(sendSmsRequest)
	return
}
