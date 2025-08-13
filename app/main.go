package main

import (
	"fmt"
	"net"
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

	switch {
	case parts[1] == "/":
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	case strings.Contains(parts[1], "/echo/"):
		res, ok := strings.CutPrefix(parts[1], "/echo/")
		if !ok {
			conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
			break
		}
		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(res), res)))
	default:
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}
