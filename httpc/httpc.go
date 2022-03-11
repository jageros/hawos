/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    post
 * @Date:    2021/6/23 11:33 上午
 * @package: http
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package httpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type METHOD string

func (m METHOD) String() string {
	return string(m)
}

const (
	JSON = "application/json"
	FORM = "application/x-www-form-urlencoded"

	GET  METHOD = "GET"
	POST METHOD = "POST"
	PUT  METHOD = "PUT"
)

var cli *http.Client

func init() {
	cli = http.DefaultClient
}

func SetHeader(req *http.Request, arg map[string]string) {
	for key, val := range arg {
		req.Header.Set(key, val)
	}
}

func Request(method METHOD, url string, contentType string, arg map[string]interface{}, header map[string]string) (result []byte, err error) {
	var data []byte

	switch contentType {
	case JSON:
		if arg != nil {
			data, err = json.Marshal(arg)
			if err != nil {
				return
			}
		}

	case FORM:
		data = marshal(arg)
	}

	req, err := http.NewRequest(method.String(), url, bytes.NewBuffer(data))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", contentType)

	SetHeader(req, header)

	resp, err := cli.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("ServerReturnErr: " + resp.Status)
	}
	return body, nil
}

func RequestWithInterface(method METHOD, url string, contentType string, arg map[string]interface{}, header map[string]string, result interface{}) error {
	body, err := Request(method, url, contentType, arg, header)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, result)
}

func RequestReturnMap(method METHOD, url string, contentType string, arg map[string]interface{}, header map[string]string) (result map[string]interface{}, err error) {
	err = RequestWithInterface(method, url, contentType, arg, header, &result)
	return
}

func marshal(arg map[string]interface{}) []byte {
	data := &url.Values{}
	for key, val := range arg {
		switch val.(type) {
		case int:
			v := strconv.Itoa(val.(int))
			data.Set(key, v)
		case string:
			data.Set(key, val.(string))
		case []string:
			vals := val.([]string)
			for _, v := range vals {
				data.Add(key, v)
			}
		}
	}
	return []byte(data.Encode())
}
