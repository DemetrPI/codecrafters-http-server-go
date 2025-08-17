package main

import "net"

type Server struct {
	listener net.Listener
}

type ContentType struct {
	value string
}

type HttpStatus int

type HttpRequest struct {
	Method  string
	URL     string
	Version string
	Headers map[string]string
	Body    string
}

type HttpResponse struct {
	Status  HttpStatus
	Version string
	Headers map[string]string
	Body    string
}

const (
	//Statuses
	StatusOk                  HttpStatus = 200
	StatusNotFound            HttpStatus = 404
	StatusInternalServerError HttpStatus = 500
	StatusBadRequest          HttpStatus = 400
	//Headers
	ContentTypeHeader string = "Content-Type"
	//Other
	CRLF      string = "\r\n"
	Port      string = ":4221"
	IpAddress string = "0.0.0.0"
	Protocol  string = "tcp"
	Version   string = "HTTP/1.1"
)

var (
	ContentTypeText        = ContentType{value: "text/plain"}
	ContentTypeHtml        = ContentType{value: "text/html"}
	ContentTypeJson        = ContentType{value: "application/json"}
	ContentTypeXml         = ContentType{value: "application/xml"}
	ContentTypeOctetStream = ContentType{value: "application/octet-stream"}
)
