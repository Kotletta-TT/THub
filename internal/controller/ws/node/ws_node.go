package node

import (
	"github.com/Kotletta-TT/THub/internal/controller"
	"github.com/Kotletta-TT/THub/logger"

	"github.com/gorilla/websocket"
)

type WSNode struct {
	ws  *websocket.Conn
	l   logger.Logger
	tmp *controller.WindowSize
}

func NewWSNode(ws *websocket.Conn, l logger.Logger) *WSNode {
	return &WSNode{ws: ws, l: l, tmp: &controller.WindowSize{}}
}

func (wn *WSNode) Read(p []byte) (n int, err error) {
	wn.l.Info("Node Read")
	t, data, err := wn.ws.ReadMessage()
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

func (wn *WSNode) Write(p []byte) (n int, err error) {
	wn.l.Info("Node Write")
	msgType := websocket.TextMessage
	if p[0] == '{' || p[0] == '[' {
		if err = wn.tmp.UnmarshalJSON(p); err == nil {
			msgType = websocket.BinaryMessage
		}
	}
	err = wn.ws.WriteMessage(msgType, p)

	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (wn *WSNode) Close() error {
	panic("not implemented") // TODO: Implement
}
