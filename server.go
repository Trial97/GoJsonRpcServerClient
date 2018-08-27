package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
)

//MyServer used to sum the numbers
type MyServer struct{}

//Sum adds numbers
func (b *MyServer) Sum(a *Args, r *int) error {
	*r = a.A + a.B
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

//StartServer for server export
func StartServer(port int) {
	arith := new(MyServer)

	server := rpc.NewServer()
	server.Register(arith)

	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

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

		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
		defer conn.Close()
	}
}

// func main() {
// 	go StartServer(8223)
// }
