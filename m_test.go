package main

import (
	"fmt"
	"os"
	"testing"
)

func TestStartServer(t *testing.T) {
	// cmd := exec.Command("./server", "&")
	// err := cmd.Start()
	// if err != nil {
	// 	t.Fatal("SeverError:", err)
	// }
	port := 8223
	go StartServer(port)
}
func TestCientSum(t *testing.T) {
	port := 8223
	// go StartServer(port)
	c := CreateNewClient("localhost", port)
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
func TestCientWrite(t *testing.T) {
	port := 8223
	// go StartServer(port)
	c := CreateNewClient("localhost", port)
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
func TestCienRead(t *testing.T) {
	port := 8223
	// go StartServer(port)
	c := CreateNewClient("localhost", port)
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
