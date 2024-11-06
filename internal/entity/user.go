package entity

import "net"

type User struct {
	Id int
	Conn net.Conn
}
