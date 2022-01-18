/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    melody
 * @Date:    2022/1/18 5:38 下午
 * @package: udpx
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package udpx

import "github.com/jageros/hawox/contextx"

type handleMessageFunc func(*Session, []byte)
type handleErrorFunc func(*Session, error)
type handleCloseFunc func(*Session, int, string) error
type handleSessionFunc func(*Session)
type filterFunc func(*Session) bool


// Melody implements a websocket manager.
type Melody struct {
	Config                   *Config
	ctx                      contextx.Context
	messageHandler           handleMessageFunc
	messageHandlerBinary     handleMessageFunc
	messageSentHandler       handleMessageFunc
	messageSentHandlerBinary handleMessageFunc
	errorHandler             handleErrorFunc
	closeHandler             handleCloseFunc
	connectHandler           handleSessionFunc
	disconnectHandler        handleSessionFunc
	pingHandler              handleSessionFunc
	//hub                      *hub
}