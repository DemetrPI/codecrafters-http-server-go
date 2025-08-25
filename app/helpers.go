package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net"
	"os"
)

func (s *Server) errHandler(message string, err error) {
	if err != nil {
		fmt.Println(message, err.Error())
		os.Exit(1)
	}
}

func (s *Server) Listen() {
	l, err := net.Listen(Protocol, IpAddress+Port)
	s.errHandler("Error listening:", err)
	s.listener = l
}

func (s *Server) Accept() net.Conn {
	conn, err := s.listener.Accept()
	s.errHandler("Error accepting connection:", err)
	return conn
}

func (s *Server) Close() {
	err := s.listener.Close()
	s.errHandler("Error closing listener:", err)
}

func (s HttpStatus) String() string {
	switch s {
	case StatusOk:
		return "200 OK"
	case StatusBadRequest:
		return "400 Bad Request"
	case StatusCreated:
		return "201 Created"
	case StatusNotFound:
		return "404 Not Found"
	case StatusInternalServerError:
		return "500 Internal Server Error"
	default:
		return fmt.Sprintf("Status %d", int(s))
	}
}

func handleCompression(responce *HttpResponse) error {
	var buffer bytes.Buffer
	gzipWriter := gzip.NewWriter(&buffer)
	if _, err := gzipWriter.Write([]byte(responce.Body)); err != nil {
		return fmt.Errorf("compression failed, %s", err.Error())
	}
	if err := gzipWriter.Flush(); err != nil {
		return fmt.Errorf("flush failed, %s", err.Error())
	}
	if err := gzipWriter.Close(); err != nil {
		return fmt.Errorf("close failed, %s", err.Error())
	}
	responce.Headers["content-length"] = fmt.Sprintf("%d", buffer.Len())
	responce.Body = buffer.String()
	return nil
}
