package main

import (
	"bufio"
	"net"
)

func main() {

	s := Server{}
	s.Listen()
	defer s.Close()

	for {
		conn := s.Accept()
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	requestData := bufio.NewReader(conn)
	defer conn.Close()
	for {

		request, err, mustClose := parseRequest(requestData)
		if err != nil {
			responce := &HttpResponse{
				Status:  StatusBadRequest,
				Headers: make(map[string]string),
			}
			sendResponce(conn, responce, request)
			return
		}
		responce := request.routeRequest()
		sendResponce(conn, responce, request)
		if mustClose || request.Headers["connection"] == "close" {
			break
		}
	}
}
