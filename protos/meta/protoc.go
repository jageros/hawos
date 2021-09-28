/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    protoc
 * @Date:    2021/6/10 10:54 上午
 * @package: protoc
 * @Version: v1.0.0
 *
 * @Description: 定义协议meta接口
 *
 */

package meta

import "errors"

var metaData = make(map[int32]IMeta)

var NoMetaErr = errors.New("NoMetaErr")

type IMeta interface {
	GetMsgID() IMsgID
	EncodeArg(interface{}) ([]byte, error)
	DecodeArg([]byte) (interface{}, error)
	EncodeReply(interface{}) ([]byte, error)
	DecodeReply([]byte) (interface{}, error)
}

type IMsgID interface {
	ID() int32
	String() string
}

func RegisterMeta(meta IMeta) {
	metaData[meta.GetMsgID().ID()] = meta
}

func GetMeta(msgId int32) (IMeta, error) {
	if m, ok := metaData[msgId]; ok {
		return m, nil
	} else {
		return nil, NoMetaErr
	}
}
