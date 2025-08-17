package main

import "net"

type HttpRequest struct {
	Method  string
	URL     string
	Version string
	Headers map[string]string
	Body    string
}

type HttpStatus int

type HttpResponse struct {
	Status  HttpStatus
	Version string
	Headers map[string]string
	Body    string
}

const (
	StatusOk                  HttpStatus = 200
	StatusNotFound            HttpStatus = 404
	StatusInternalServerError HttpStatus = 500
	StatusBadRequest          HttpStatus = 400
	CRLF                      string     = "\r\n"
	Port                      string     = ":4221"
	IpAddress                 string     = "0.0.0.0"
	Protocol                  string     = "tcp"
	Version                   string     = "HTTP/1.1"
	ContentType               string     = "content-type"
)

type Server struct {
	listener net.Listener
}
