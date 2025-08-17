package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func (s *Server) Listen() {
	l, err := net.Listen(Protocol, IpAddress+Port)
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

func (s HttpStatus) String() string {
	switch s {
	case StatusOk:
		return "200 OK"
	case StatusBadRequest:
		return "400 Bad Request"
	case StatusNotFound:
		return "404 Not Found"
	case StatusInternalServerError:
		return "500 Internal Server Error"
	default:
		return fmt.Sprintf("Status %d", int(s))
	}
}

func parseRequest(conn net.Conn) (*HttpRequest, error) {
	requestData := bufio.NewReader(conn)

	//Read request
	requestLine, err := requestData.ReadString('\n')
	if err != nil {
		return nil, err
	}

	requestParts := strings.Fields(requestLine)
	if len(requestParts) < 3 {
		return nil, fmt.Errorf("invalid request")
	}

	requestContent := &HttpRequest{
		Method:  requestParts[0],
		URL:     requestParts[1],
		Version: requestParts[2],
		Headers: make(map[string]string),
	}

	//Read headers
	for {
		rawLine, err := requestData.ReadString('\n')
		strippedLine := strings.TrimSpace(rawLine)
		if err != nil || strippedLine == "" {
			break
		}
		headerParts := strings.Split(strippedLine, ":")
		if len(headerParts) == 2 {
			key := strings.ToLower(headerParts[0])
			value := strings.TrimSpace(headerParts[1])
			requestContent.Headers[key] = value
		}
	}
	fmt.Println("RequestContent:==>", requestContent)
	return requestContent, nil
}

func sendResponce(conn net.Conn, res *HttpResponse) {
	//status line
	requestResponse := fmt.Sprintf("HTTP/1.1 %s%s", res.Status, CRLF)

	//headers
	res.Headers["content-length"] = fmt.Sprintf("%d", len(res.Body))
	for key, value := range res.Headers {
		requestResponse += fmt.Sprintf("%s: %s%s", key, value, CRLF)
	}

	//body
	requestResponse += CRLF + res.Body
	fmt.Println("Responce:==>", requestResponse)
	conn.Write([]byte(requestResponse))
}

func (request *HttpRequest) routeRequest() *HttpResponse {
	path := strings.Trim(request.URL, "/") //Remove leading and trailing '/'
	parts := strings.Split(path, "/")
	fmt.Println("Parts:==>", parts)

	switch {
	case request.URL == "/":
		return &HttpResponse{
			Status:  StatusOk,
			Version: request.Version,
			Headers: make(map[string]string),
		}
	case parts[0] == "echo" && len(parts) == 2:
		return &HttpResponse{
			Status:  StatusOk,
			Version: request.Version,
			Headers: map[string]string{"content-type": "text/plain"},
			Body:    parts[1],
		}
	case parts[0] == "user-agent" && len(parts) == 1:
		return &HttpResponse{
			Status:  StatusOk,
			Version: request.Version,
			Headers: map[string]string{"content-type": "text/plain"},
			Body:    request.Headers["user-agent"],
		}
	default:
		return &HttpResponse{
			Status:  StatusNotFound,
			Headers: make(map[string]string),
		}
	}
}
