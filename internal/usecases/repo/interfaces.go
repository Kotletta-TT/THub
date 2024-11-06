package repo

import "io"

type NodeInterface interface {
	GetId() string
	GetConn() io.ReadWriter
}
