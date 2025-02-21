package common

import (
	"context"
	"net"
)

type Dialer func(ctx context.Context, addr string) (net.Conn, error)

type ListenerFactory func(ctx context.Context, addr string) (net.Listener, error)

func DialTCP(ctx context.Context, addr string) (net.Conn, error) {
	return net.Dial("tcp", addr)
}

func ListenTCP(ctx context.Context, addr string) (net.Listener, error) {
	return net.Listen("tcp", addr)
}
