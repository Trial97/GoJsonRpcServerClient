package main

import (
	"io"

	"github.com/go-zeromq/zmq4"
)

// NewZmqMessage return a io.ReadWriteCloser for a message
func NewZmqMessage(socket zmq4.Socket, msg zmq4.Msg) (cln *ZmqMessage) {
	cln = &ZmqMessage{
		socket: socket,
		ID:     string(msg.Frames[0]),
		data:   make([]byte, 0),
	}
	for _, frame := range msg.Frames[2:] {
		cln.data = append(cln.data, frame...) // add the data ignoring the first 2 frames
	}
	return
}

// ZmqMessage implements io.ReadWriteCloser based on the 0MQ socket
type ZmqMessage struct {
	ID     string
	socket zmq4.Socket // used to send the response
	data   []byte
	indx   int64
}

// Close closes the socket
func (c *ZmqMessage) Close() (err error) {
	return
}

// Read read a mesg and try to populate the p using it's content
func (c *ZmqMessage) Read(p []byte) (n int, err error) {
	if c.indx >= int64(len(c.data)) {
		return 0, io.EOF
	}
	n = copy(p, c.data[c.indx:])
	c.indx += int64(n)
	return
}

// Write writes a new message
func (c *ZmqMessage) Write(p []byte) (n int, err error) {
	return len(p), c.socket.Send(zmq4.NewMsgFrom([]byte(c.ID), []byte{}, p)) // id, empty frame, message
}
