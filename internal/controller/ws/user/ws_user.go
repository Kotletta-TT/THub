package user

import (
	"github.com/Kotletta-TT/THub/internal/controller"
	"github.com/Kotletta-TT/THub/logger"

	"github.com/gorilla/websocket"
)

type WSUser struct {
	ws  *websocket.Conn
	l   logger.Logger
	tmp *controller.WindowSize
}

func NewWSUser(ws *websocket.Conn, l logger.Logger) *WSUser {
	return &WSUser{ws: ws, l: l, tmp: &controller.WindowSize{}}
}

func (wu *WSUser) Read(p []byte) (n int, err error) {
	wu.l.Info("User Read")
	t, data, err := wu.ws.ReadMessage()
	switch t {
	case websocket.CloseMessage:
		panic("not implemented") // TODO: Implement
	case websocket.TextMessage:
		n = copy(p, data)
	case websocket.BinaryMessage:
		n = copy(p, data)
	}
	return n, err
}

func (wu *WSUser) Write(p []byte) (n int, err error) {
	wu.l.Info("User Write")
	msgType := websocket.TextMessage
	if p[0] == '{' || p[0] == '[' {
		if err = wu.tmp.UnmarshalJSON(p); err == nil {
			msgType = websocket.BinaryMessage
		}
	}
	err = wu.ws.WriteMessage(msgType, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (wu *WSUser) Close() error {
	panic("not implemented") // TODO: Implement
}
