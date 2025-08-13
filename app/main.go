package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const CRLF = "\r\n"

type Server struct {
	listener net.Listener
}

func main() {

	s := Server{}
	s.Listen()
	defer s.Close()

	for {
		conn := s.Accept()
		go s.HandleConnection(conn)
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Client connected from:", conn.RemoteAddr())

	var buf = make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading: ", err.Error())
		return
	}

	req := string(buf)
	lines := strings.Split(req, CRLF)

	requestLine := lines[0]
	parts := strings.Split(requestLine, " ")
	if len(parts) < 2 {
		fmt.Println("Invalid request")
		return
	}
	path := parts[1]

	fmt.Println("Request path: ", path)

	var response string
	if path == "/" {
		response = "HTTP/1.1 200 OK\r\n\r\n"
	} else {
		response = "HTTP/1.1 404 Not Found\r\n\r\n"
	}
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing: ", err.Error())
	}
}

func (s *Server) Listen() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Error listening: ", err.Error())
		os.Exit(1)
	}
	s.listener = l
}

func (s *Server) Accept() net.Conn {
	conn, err := s.listener.Accept()
	if err != nil {
		fmt.Println("Error accepting: ", err.Error())
		os.Exit(1)
	}
	return conn
}

func (s *Server) Close() {
	err := s.listener.Close()
	if err != nil {
		fmt.Println("Error closing listener: ", err.Error())
		os.Exit(1)
	}
}
