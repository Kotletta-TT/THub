package node

import (
	"fmt"
	"net/http"
	"net/url"

	uc "github.com/Kotletta-TT/THub/internal/usecases"
	"github.com/Kotletta-TT/THub/logger"

	"github.com/gorilla/websocket"
)

const maxMessageSize = 512

var upgrader = websocket.Upgrader{
	ReadBufferSize:  maxMessageSize,
	WriteBufferSize: maxMessageSize,
}

type NodeRoutes struct {
	uc *uc.ConnectUseCase
	ud *uc.DisconnectUseCase
	l  logger.Logger
}

func NewNodeRoutes(uc *uc.ConnectUseCase, ud *uc.DisconnectUseCase, l logger.Logger) *NodeRoutes {
	return &NodeRoutes{
		uc: uc,
		ud: ud,
		l:  l,
	}
}

func (n *NodeRoutes) Connect(w http.ResponseWriter, r *http.Request) {
	nodeId, err := getNodeId(r.URL)
	if err != nil {
		return
		//TODO prepare error (close ws)
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	node := NewWSNode(conn, n.l)
	ctx, err := n.uc.NodeConnect(nodeId, node)
	if err != nil {
		return
		//TODO prepare error (close ws)
	}
	// TODO log Disconnect error
	defer n.ud.Disconnect(nodeId)
	<-ctx.Done()
}

func getNodeId(r *url.URL) (string, error) {
	params, err := url.ParseQuery(r.RawQuery)
	if err != nil {
		return "", err
	}
	if params.Has("nodeId") {
		return params.Get("nodeId"), nil
	}
	return "", fmt.Errorf("nodeId not found")
}
