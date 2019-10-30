package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/go-zeromq/zmq4"
)

// NewZmqClient3 return a io.ReadWriteCloser for a REQ socket type
// http://api.zeromq.org/2-1:zmq-socket
func NewZmqClient3(endpoint string) (cln io.ReadWriteCloser, err error) {
	dealer := zmq4.NewReq(context.Background(), zmq4.WithID(zmq4.SocketIdentity("uniq []byte")))

	if err = dealer.Dial("tcp://" + endpoint); err != nil {
		return
	}
	cln = &ZmqClient3{
		socket: dealer,
	}
	return
}

// ZmqClient3 implements io.ReadWriteCloser based on the 0MQ socket
type ZmqClient3 struct {
	socket zmq4.Socket

	// used for reading the message
	prevData []byte
	indx     int64
}

// Close closes the socket
func (c *ZmqClient3) Close() error {
	return c.socket.Close()
}

// Read read a mesg and try to populate the p using it's content
func (c *ZmqClient3) Read(p []byte) (n int, err error) {
	if c.indx >= int64(len(c.prevData)) {
		c.indx = 0
		var msg zmq4.Msg
		msg, err = c.socket.Recv()
		if err != nil {
			return 0, err
		}
		c.prevData = msg.Bytes() // handle all frames?
	}
	n = copy(p, c.prevData[c.indx:])
	c.indx += int64(n)
	return
}

// Write writes a new message
func (c *ZmqClient3) Write(p []byte) (n int, err error) {
	return len(p), c.socket.Send(zmq4.NewMsgFrom(p))
}

func CreateNewZMQClient(ip string, port int) *rpc.Client {
	is := fmt.Sprintf("%s:%d", ip, port)
	conn, err := NewZmqClient3(is)

	if err != nil {
		log.Println("CreateNewClientError:", err)
		return nil
	}
	// defer conn.Close()
	return jsonrpc.NewClient(conn)
}
