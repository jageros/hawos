package wsc

import (
	"errors"
	"github.com/jageros/hawox/contextx"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Session wrapper around websocket connections.
type Session struct {
	Response *http.Response
	Keys     map[string]interface{}
	conn     *websocket.Conn
	output   chan *envelope
	melody   *Melody
	open     bool
	rwmutex  *sync.RWMutex
}

func (s *Session) Send(bytes []byte) {
	msg := &envelope{
		t:   websocket.TextMessage,
		msg: bytes,
	}
	s.writeMessage(msg)
}

func (s *Session) writeMessage(message *envelope) {
	if s.closed() {
		s.melody.errorHandler(s, errors.New("tried to write to closed a session"))
		return
	}

	select {
	case s.output <- message:
	default:
		s.melody.errorHandler(s, errors.New("session message buffer is full"))
	}
}

func (s *Session) writeRaw(message *envelope) error {
	if s.closed() {
		return errors.New("tried to write to a closed session")
	}

	s.conn.SetWriteDeadline(time.Now().Add(s.melody.Config.WriteWait))
	err := s.conn.WriteMessage(message.t, message.msg)

	if err != nil {
		return err
	}

	return nil
}

func (s *Session) closed() bool {
	s.rwmutex.RLock()
	defer s.rwmutex.RUnlock()

	return !s.open
}

func (s *Session) close() {
	if !s.closed() {
		s.rwmutex.Lock()
		s.open = false
		s.conn.Close()
		close(s.output)
		s.rwmutex.Unlock()
	}
}

func (s *Session) pong() {
	s.writeRaw(&envelope{t: websocket.PongMessage, msg: []byte{}})
}

func (s *Session) writePump(ctx contextx.Context) error {
loop:
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-s.output:
			if !ok {
				break loop
			}

			err := s.writeRaw(msg)

			if err != nil {
				s.melody.errorHandler(s, err)
				return err
			}

			if msg.t == websocket.CloseMessage {
				break loop
			}

			if msg.t == websocket.TextMessage {
				s.melody.messageSentHandler(s, msg.msg)
			}

			if msg.t == websocket.BinaryMessage {
				s.melody.messageSentHandlerBinary(s, msg.msg)
			}
		}
	}
	return nil
}

func (s *Session) readPump(ctx contextx.Context) error {
	//s.conn.SetReadLimit(s.melody.Config.MaxMessageSize)
	err := s.conn.SetReadDeadline(time.Now().Add(s.melody.Config.PongWait))
	if err != nil {
		return err
	}

	//s.conn.SetPongHandler(func(string) error {
	//	s.conn.SetReadDeadline(time.Now().Add(s.melody.Config.PongWait))
	//	s.melody.pingHandler(s)
	//	return nil
	//})

	s.conn.SetPingHandler(func(appData string) error {
		err = s.conn.SetReadDeadline(time.Now().Add(s.melody.Config.PongWait))
		s.pong()
		s.melody.pingHandler(s)
		return err
	})

	if s.melody.closeHandler != nil {
		s.conn.SetCloseHandler(func(code int, text string) error {
			return s.melody.closeHandler(s, code, text)
		})
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			t, message, err := s.conn.ReadMessage()

			if err != nil {
				s.melody.errorHandler(s, err)
				return err
			}

			if t == websocket.TextMessage {
				s.melody.messageHandler(s, message)
			}

			if t == websocket.BinaryMessage {
				s.melody.messageHandlerBinary(s, message)
			}
		}
	}
}

// Write writes message to session.
func (s *Session) Write(msg []byte) error {
	if s.closed() {
		return errors.New("session is closed")
	}

	s.writeMessage(&envelope{t: websocket.TextMessage, msg: msg})

	return nil
}

// WriteBinary writes a binary message to session.
func (s *Session) WriteBinary(msg []byte) error {
	if s.closed() {
		return errors.New("session is closed")
	}

	s.writeMessage(&envelope{t: websocket.BinaryMessage, msg: msg})

	return nil
}

// Close closes session.
func (s *Session) Close() error {
	if s.closed() {
		return errors.New("session is already closed")
	}

	s.writeMessage(&envelope{t: websocket.CloseMessage, msg: []byte{}})

	return nil
}

// CloseWithMsg closes the session with the provided payload.
// Use the FormatCloseMessage function to format a proper close message payload.
func (s *Session) CloseWithMsg(msg []byte) error {
	if s.closed() {
		return errors.New("session is already closed")
	}

	s.writeMessage(&envelope{t: websocket.CloseMessage, msg: msg})

	return nil
}

// IsClosed returns the status of the connection.
func (s *Session) IsClosed() bool {
	return s.closed()
}

// Set is used to store a new key/value pair exclusivelly for this session.
// It also lazy initializes s.Keys if it was not used previously.
func (s *Session) Set(key string, value interface{}) {
	if s.Keys == nil {
		s.Keys = make(map[string]interface{})
	}

	s.Keys[key] = value
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (s *Session) Get(key string) (value interface{}, exists bool) {
	if s.Keys != nil {
		value, exists = s.Keys[key]
	}

	return
}

// MustGet returns the value for the given key if it exists, otherwise it panics.
func (s *Session) MustGet(key string) interface{} {
	if value, exists := s.Get(key); exists {
		return value
	}

	panic("Key \"" + key + "\" does not exist")
}
