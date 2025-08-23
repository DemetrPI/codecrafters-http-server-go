package main

import "net"


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
	defer conn.Close()

	request, err := parseRequest(conn)
	if err != nil {
		responce := &HttpResponse{
			Status:  StatusBadRequest,
			Headers: make(map[string]string),
		}
		sendResponce(conn, responce)
		return
		}
	responce := request.routeRequest()
	sendResponce(conn, responce)
}


