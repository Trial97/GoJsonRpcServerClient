package main

import (
	"context"
	"fmt"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/go-zeromq/zmq4"
)

func NewZmqServer() *ZmqServer {
	return &ZmqServer{
		end:    make(chan struct{}, 1),
		socket: zmq4.NewRouter(context.Background(), zmq4.WithID(zmq4.SocketIdentity("router"))),
	}
}

type ZmqServer struct {
	socket zmq4.Socket
	end    chan struct{}
}

func (s *ZmqServer) LisenAndServeRPC(endpoint string) (err error) {
	err = s.socket.Listen("tcp://" + endpoint)
	if err != nil {
		return err
	}
	defer s.socket.Close()

	// create rpc Server
	arith := new(MyServer)
	server := rpc.NewServer()
	server.Register(arith)

	for {
		select {
		case <-s.end: // stop the server
			return nil
		default:
		}
		msg, err := s.socket.Recv() // this is something like accept but return only the message
		if err != nil {
			return err
		}
		if len(msg.Frames) < 3 { // the message should have at least 3 frames: one for socketID,one empty and the rest should be the message
			continue
		}
		go server.ServeCodec(jsonrpc.NewServerCodec(NewZmqMessage(s.socket, msg)))
	}
}

func (s *ZmqServer) Close() (err error) {
	s.end <- struct{}{}
	close(s.end)
	return s.socket.Close()
}

func StartZMQServer(port int) {
	sport := fmt.Sprintf(":%d", port)
	server := NewZmqServer()
	server.LisenAndServeRPC(sport)
}
