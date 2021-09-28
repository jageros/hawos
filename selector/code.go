/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    decode
 * @Date:    2021/7/9 3:18 下午
 * @package: dscv
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package selector

import (
	"encoding/json"
)

type metaData struct {
	ID     string  `json:"Id"`
	Name   string  `json:"Name"`
	MsgIds []int32 `json:"msg_ids"`
}

func marshal(md *metaData) (string, error) {
	data, err := json.Marshal(md)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func unmarshal(data []byte) (md *metaData, err error) {
	err = json.Unmarshal(data, &md)
	return
}
