package entity

import "io"

type Node struct {
	id   string
	conn io.ReadWriter
}

func NewNode(id string, conn io.ReadWriter) *Node {
	return &Node{
		id:   id,
		conn: conn,
	}
}

func (n *Node) GetId() string {
	return n.id
}

func (n *Node) GetConn() io.ReadWriter {
	return n.conn
}
