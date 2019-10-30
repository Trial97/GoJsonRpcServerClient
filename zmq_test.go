package main

import (
	"fmt"
	"os"
	"testing"
)

func TestZMQStartServer(t *testing.T) {
	// cmd := exec.Command("./server", "&")
	// err := cmd.Start()
	// if err != nil {
	// 	t.Fatal("SeverError:", err)
	// }
	port := 8224
	go StartZMQServer(port)
}
func TestZMQCientSum(t *testing.T) {
	port := 8224
	// go StartServer(port)
	c := CreateNewZMQClient("localhost", port)
	if c == nil {
		t.Fatal("user not connected")
	}
	defer c.Close()
	A, B := 10, 15
	var r int
	CallSumFunc(c, A, B, &r)
	if r != A+B {
		t.Fatal("Sum is wrong")
	}
}
func TestZMQCientWrite(t *testing.T) {
	port := 8224
	// go StartServer(port)
	c := CreateNewZMQClient("localhost", port)
	if c == nil {
		t.Fatal("user not connected")
	}
	defer c.Close()
	A, F := 10, "t.txt"
	var r int
	var b bool
	CallWriteNumber(c, A, F, &b)
	fd, err := os.Open(F)
	if err != nil {
		t.Fatal("FileError:", err)
	}
	fmt.Fscanf(fd, "%d\n", &r)
	if r != A {
		t.Fatal("number not saved")
	}
}
func TestZMQCienRead(t *testing.T) {
	port := 8224
	// go StartServer(port)
	c := CreateNewZMQClient("localhost", port)
	if c == nil {
		t.Fatal("user not connected")
	}
	defer c.Close()
	A, F := 10, "t2.txt"

	fo, err := os.Create(F)
	if err != nil {
		t.Fatal("FileError:", err)
	}
	defer fo.Close()
	fmt.Fprintf(fo, "%d", A)

	var r int
	CallReadNumber(c, F, &r)

	if r != A {
		t.Fatal("number not saved")
	}
}
