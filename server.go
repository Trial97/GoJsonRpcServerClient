package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"time"
)

//MyServer used to sum the numbers
type MyServer struct {
	A int
}

//Sum adds numbers
func (b *MyServer) Sum(a *Args, r *int) error {
	*r = a.A + a.B
	if b.A != 0 {
		fmt.Println(b.A)
	}
	return nil
}

//WriteNumber write number to file;
func (b *MyServer) WriteNumber(a Argsw, reply *bool) error {
	fo, err := os.Create(a.F)
	if err != nil {
		log.Println("WriteNumberError:", err)
		return err
	}
	defer fo.Close()
	fmt.Fprintf(fo, "%d", a.A)
	*reply = true
	return nil
}

//ReadNumber reads number from file
func (b *MyServer) ReadNumber(f string, r *int) error {
	fd, err := os.Open(f)
	if err != nil {
		log.Println("ReadNumberError:", err)
		return err
	}
	fmt.Fscanf(fd, "%d\n", r)
	return nil
}

var (
	server *rpc.Server
	conns  []net.Conn
	wait   chan struct{}
)

//StartServer for server export
func StartServer(port int) {
	wait = make(chan struct{}, 1)
	wait <- struct{}{}
	arith := new(MyServer)

	server = rpc.NewServer()
	server.Register(arith)

	// server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

	sport := fmt.Sprintf(":%d", port)
	l, e := net.Listen("tcp", sport)
	if e != nil {
		log.Fatal("ListenError:", e)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("AcceptError", err)
		}

		// if reload wait until we switch the server

		fmt.Println("WaitForLock in for")
		b := <-wait
		wait <- b
		fmt.Println("DoneLock in for")
		conns = append(conns, conn)
		go func(conn net.Conn) {
			server.ServeCodec(jsonrpc.NewServerCodec(conn))
			conn.Close()
			// remove from slice/map
		}(conn)
	}
}

func switchServer(A int) {
	arith := new(MyServer)
	arith.A = A

	// create a new server with all the services
	newServer := rpc.NewServer()
	newServer.Register(arith)

	// aquiere the lock
	fmt.Println("WaitForLock", len(wait))
	lock := <-wait
	fmt.Println("AquieredLock")
	server = newServer
	oldConns := conns
	conns = nil
	wait <- lock

	// close or set the deadLine for the old connections
	for _, conn := range oldConns {
		// conn.Close()
		conn.SetReadDeadline(time.Now())
		// conn.SetDeadline(time.Now()) // the clients will receive `connection is shut down` that we should handle in rpcclient library to reconnect
		// rpc.ErrShutdown and maybe io.ErrUnexpectedEOF
	}
}

// func main() {
// 	go StartServer(8223)
// }
