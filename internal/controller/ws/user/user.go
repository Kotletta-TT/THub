package user

import (
	"encoding/json"
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

type UserRoutes struct {
	tuc  *uc.TransferUseCase
	lnuc *uc.ListNodesUseCase
	l    logger.Logger
}

func NewUserRoutes(tuc *uc.TransferUseCase, lnuc *uc.ListNodesUseCase, l logger.Logger) *UserRoutes {
	return &UserRoutes{
		tuc:  tuc,
		lnuc: lnuc,
		l:    l,
	}
}

func (ur *UserRoutes) TransferToNode(w http.ResponseWriter, r *http.Request) {
	nodeId, err := getNodeId(r.URL)
	if err != nil {
		ur.l.Error(err.Error()) // fix me to standard
		return
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ur.l.Error("Upgrade user ws error!") // fix me to standard
		return
	}
	wsUsr := NewWSUser(c, ur.l)

	ur.tuc.Transfer(wsUsr, nodeId)
}

func (ur *UserRoutes) ListNodes(w http.ResponseWriter, r *http.Request) {
	nodes := ur.lnuc.ListNodes()
	jNodes, err := json.Marshal(nodes)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		ur.l.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jNodes)
	w.WriteHeader(http.StatusOK)
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
