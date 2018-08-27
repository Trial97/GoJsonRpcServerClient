package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

//CreateNewClient creates new client
func CreateNewClient(ip string, port int) *rpc.Client {
	is := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.Dial("tcp", is)

	if err != nil {
		log.Println("CreateNewClientError:", err)
		return nil
	}
	// defer conn.Close()
	return jsonrpc.NewClient(conn)
}

//CallSumFunc cals sum from server
func CallSumFunc(c *rpc.Client, A, B int, reply *int) error {
	args := &Args{A, B}
	return c.Call("MyServer.Sum", args, reply)
}

//CallWriteNumber cals WriteNumber from server
func CallWriteNumber(c *rpc.Client, A int, F string, reply *bool) error {
	args := &Argsw{A, F}
	return c.Call("MyServer.WriteNumber", args, reply)
}

//CallReadNumber cals ReadNumber from server
func CallReadNumber(c *rpc.Client, F string, reply *int) error {
	return c.Call("MyServer.ReadNumber", F, reply)
}

func main() {
	go StartServer(8223)
	c := CreateNewClient("localhost", 8223)
	defer c.Close()
	var R int
	A := 10
	B := 12
	e := CallSumFunc(c, A, B, &R)
	if e != nil {
		fmt.Print(e)
	}
	fmt.Printf("%d+%d=%d\n", A, B, R)
	var b bool
	e = CallWriteNumber(c, R, "f.txt", &b)
	if e != nil {
		fmt.Print(e)
	}
	e = CallReadNumber(c, "f.txt", &R)
	if e != nil {
		fmt.Print(e)
	}
	fmt.Println(R)
}
