package main

import "net"

type Server struct {
	listener net.Listener
}

type ContentType struct {
	value string
}

type HttpStatus int

type HTTPMethod string

type HttpRequest struct {
	Method  HTTPMethod
	URL     string
	Version string
	Headers map[string]string
	Cookies map[string]string
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
	StatusCreated             HttpStatus = 201
	StatusMovedPermanently    HttpStatus = 301
	StatusNotFound            HttpStatus = 404
	StatusInternalServerError HttpStatus = 500
	StatusBadRequest          HttpStatus = 400
	//Headers
	ContentTypeHeader string = "Content-Type"
	//HTTPMethods
	HTTPGet    HTTPMethod = "GET"
	HTTPPost   HTTPMethod = "POST"
	HTTPPut    HTTPMethod = "PUT"
	HTTPDelete HTTPMethod = "DELETE"
	//Others
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
