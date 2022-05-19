package wsc

import (
	"context"
	"github.com/jageros/hawox/contextx"
	"github.com/jageros/hawox/logx"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Close codes defined in RFC 6455, section 11.7.
// Duplicate of codes from gorilla/websocket for convenience.
const (
	CloseNormalClosure           = 1000
	CloseGoingAway               = 1001
	CloseProtocolError           = 1002
	CloseUnsupportedData         = 1003
	CloseNoStatusReceived        = 1005
	CloseAbnormalClosure         = 1006
	CloseInvalidFramePayloadData = 1007
	ClosePolicyViolation         = 1008
	CloseMessageTooBig           = 1009
	CloseMandatoryExtension      = 1010
	CloseInternalServerErr       = 1011
	CloseServiceRestart          = 1012
	CloseTryAgainLater           = 1013
	CloseTLSHandshake            = 1015
)

// Duplicate of codes from gorilla/websocket for convenience.
var validReceivedCloseCodes = map[int]bool{
	// see http://www.iana.org/assignments/websocket/websocket.xhtml#close-code-number

	CloseNormalClosure:           true,
	CloseGoingAway:               true,
	CloseProtocolError:           true,
	CloseUnsupportedData:         true,
	CloseNoStatusReceived:        false,
	CloseAbnormalClosure:         false,
	CloseInvalidFramePayloadData: true,
	ClosePolicyViolation:         true,
	CloseMessageTooBig:           true,
	CloseMandatoryExtension:      true,
	CloseInternalServerErr:       true,
	CloseServiceRestart:          true,
	CloseTryAgainLater:           true,
	CloseTLSHandshake:            false,
}

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

// New creates a new melody instance with default Upgrader and Config.
func New(ctx contextx.Context) *Melody {

	//hub := newHub()

	//ctx.Go(func(ctx contextx.Context) error {
	//	hub.run()
	//	return nil
	//})

	return &Melody{
		Config:                   newConfig(),
		ctx:                      ctx,
		messageHandler:           func(*Session, []byte) {},
		messageHandlerBinary:     func(*Session, []byte) {},
		messageSentHandler:       func(*Session, []byte) {},
		messageSentHandlerBinary: func(*Session, []byte) {},
		errorHandler:             func(*Session, error) {},
		closeHandler:             nil,
		connectHandler:           func(*Session) {},
		disconnectHandler:        func(*Session) {},
		pingHandler:              func(*Session) {},
		//hub:                      hub,
	}
}

// HandleConnect fires fn when a session connects.
func (m *Melody) HandleConnect(fn func(*Session)) {
	m.connectHandler = fn
}

// HandleDisconnect fires fn when a session disconnects.
func (m *Melody) HandleDisconnect(fn func(*Session)) {
	m.disconnectHandler = fn
}

// HandlePing fires fn when a pong is received from a session.
func (m *Melody) HandlePing(fn func(*Session)) {
	m.pingHandler = fn
}

// HandleMessage fires fn when a text message comes in.
func (m *Melody) HandleMessage(fn func(*Session, []byte)) {
	m.messageHandler = fn
}

// HandleMessageBinary fires fn when a binary message comes in.
func (m *Melody) HandleMessageBinary(fn func(*Session, []byte)) {
	m.messageHandlerBinary = fn
}

// HandleSentMessage fires fn when a text message is successfully sent.
func (m *Melody) HandleSentMessage(fn func(*Session, []byte)) {
	m.messageSentHandler = fn
}

// HandleSentMessageBinary fires fn when a binary message is successfully sent.
func (m *Melody) HandleSentMessageBinary(fn func(*Session, []byte)) {
	m.messageSentHandlerBinary = fn
}

// HandleError fires fn when a session has an error.
func (m *Melody) HandleError(fn func(*Session, error)) {
	m.errorHandler = fn
}

// HandleClose sets the handler for close messages received from the session.
// The code argument to h is the received close code or CloseNoStatusReceived
// if the close message is empty. The default close handler sends a close frame
// back to the session.
//
// The application must read the connection to process close messages as
// described in the section on Control Frames above.
//
// The connection read methods return a CloseError when a close frame is
// received. Most applications should handle close messages as part of their
// normal error handling. Applications should only set a close handler when the
// application must perform some action before sending a close frame back to
// the session.
func (m *Melody) HandleClose(fn func(*Session, int, string) error) {
	if fn != nil {
		m.closeHandler = fn
	}
}

// HandleRequest upgrades http requests to websocket connections and dispatches them to be handled by the melody instance.
func (m *Melody) Connect(addr string, keys ...map[string]interface{}) (*Session, error) {
	return m.ConnectWithHeader(addr, nil, keys...)
}

// HandleRequestWithKeys does the same as HandleRequest but populates session.Keys with keys.
func (m *Melody) ConnectWithHeader(addr string, header http.Header, keys ...map[string]interface{}) (*Session, error) {
	//if m.hub.closed() {
	//	return errors.New("melody instance is closed")
	//}

	ctx, cancel := context.WithTimeout(m.ctx, time.Second*5)
	defer cancel()
	conn, r, err := websocket.DefaultDialer.DialContext(ctx, addr, header)

	if err != nil {
		return nil, err
	}

	kvs := map[string]interface{}{}

	for _, ks := range keys {
		for k, v := range ks {
			kvs[k] = v
		}
	}

	session := &Session{
		Response: r,
		Keys:     kvs,
		conn:     conn,
		output:   make(chan *envelope, m.Config.MessageBufferSize),
		melody:   m,
		open:     true,
		rwmutex:  &sync.RWMutex{},
	}

	//m.hub.register <- session

	m.connectHandler(session)

	ctx_, _ := contextx.WithCancel(m.ctx)
	ctx_.Go(func(ctx context.Context) error {
		return session.writePump(ctx)
	})

	ctx_.Go(func(ctx context.Context) error {
		return session.readPump(ctx)
	})

	ctx_.Go(func(ctx context.Context) error {
		<-ctx.Done()
		//if !m.hub.closed() {
		//	select {
		//	case m.hub.unregister <- session:
		//	default:
		//	}
		//}
		session.close()
		m.disconnectHandler(session)
		return ctx.Err()
	})
	m.ctx.Go(func(_ context.Context) error {
		err := ctx_.Wait()
		logx.Err(err).Msg("session ctx done")
		return nil
	})

	return session, nil
}

// Broadcast broadcasts a text message to all sessions.
//func (m *Melody) Broadcast(msg []byte) error {
//	if m.hub.closed() {
//		return errors.New("melody instance is closed")
//	}
//
//	message := &envelope{t: websocket.TextMessage, msg: msg}
//	m.hub.broadcast <- message
//
//	return nil
//}

// BroadcastFilter broadcasts a text message to all sessions that fn returns true for.
//func (m *Melody) BroadcastFilter(msg []byte, fn func(*Session) bool) error {
//	if m.hub.closed() {
//		return errors.New("melody instance is closed")
//	}
//
//	message := &envelope{t: websocket.TextMessage, msg: msg, filter: fn}
//	m.hub.broadcast <- message
//
//	return nil
//}

// BroadcastOthers broadcasts a text message to all sessions except session s.
//func (m *Melody) BroadcastOthers(msg []byte, s *Session) error {
//	return m.BroadcastFilter(msg, func(q *Session) bool {
//		return s != q
//	})
//}

// BroadcastMultiple broadcasts a text message to multiple sessions given in the sessions slice.
//func (m *Melody) BroadcastMultiple(msg []byte, sessions []*Session) error {
//	for _, sess := range sessions {
//		if writeErr := sess.Write(msg); writeErr != nil {
//			return writeErr
//		}
//	}
//	return nil
//}

// BroadcastBinary broadcasts a binary message to all sessions.
//func (m *Melody) BroadcastBinary(msg []byte) error {
//	if m.hub.closed() {
//		return errors.New("melody instance is closed")
//	}
//
//	message := &envelope{t: websocket.BinaryMessage, msg: msg}
//	m.hub.broadcast <- message
//
//	return nil
//}

// BroadcastBinaryFilter broadcasts a binary message to all sessions that fn returns true for.
//func (m *Melody) BroadcastBinaryFilter(msg []byte, fn func(*Session) bool) error {
//	if m.hub.closed() {
//		return errors.New("melody instance is closed")
//	}
//
//	message := &envelope{t: websocket.BinaryMessage, msg: msg, filter: fn}
//	m.hub.broadcast <- message
//
//	return nil
//}

// BroadcastBinaryOthers broadcasts a binary message to all sessions except session s.
//func (m *Melody) BroadcastBinaryOthers(msg []byte, s *Session) error {
//	return m.BroadcastBinaryFilter(msg, func(q *Session) bool {
//		return s != q
//	})
//}

// Close closes the melody instance and all connected sessions.
//func (m *Melody) Close() error {
//	if m.hub.closed() {
//		return errors.New("melody instance is already closed")
//	}
//
//	m.hub.exit <- &envelope{t: websocket.CloseMessage, msg: []byte{}}
//
//	return nil
//}

// CloseWithMsg closes the melody instance with the given close payload and all connected sessions.
// Use the FormatCloseMessage function to format a proper close message payload.
//func (m *Melody) CloseWithMsg(msg []byte) error {
//	if m.hub.closed() {
//		return errors.New("melody instance is already closed")
//	}
//
//	m.hub.exit <- &envelope{t: websocket.CloseMessage, msg: msg}
//
//	return nil
//}

// Len return the number of connected sessions.
//func (m *Melody) Len() int {
//	return m.hub.len()
//}

// IsClosed returns the status of the melody instance.
//func (m *Melody) IsClosed() bool {
//	return m.hub.closed()
//}

// FormatCloseMessage formats closeCode and text as a WebSocket close message.
//func FormatCloseMessage(closeCode int, text string) []byte {
//	return websocket.FormatCloseMessage(closeCode, text)
//}
