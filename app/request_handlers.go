package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

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
	requestResponse := fmt.Sprint(res.Version, " ", res.Status, CRLF)

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

func handleFilesRequest(filename string) *HttpResponse {
	args := os.Args
	if len(args) < 2 {
		return &HttpResponse{
			Status:  StatusInternalServerError,
			Version: Version,
			Headers: make(map[string]string),
		}
	}
	directoryName := args[2]

	content, err := os.ReadFile(directoryName + filename)
	if err != nil {
		return &HttpResponse{
			Status:  StatusNotFound,
			Version: Version,
			Headers: make(map[string]string),
		}
	}
	return &HttpResponse{
		Status:  StatusOk,
		Version: Version,
		Headers: map[string]string{ContentTypeHeader: ContentTypeOctetStream.value},
		Body:    string(content),
	}
}

func (request *HttpRequest) routeRequest() *HttpResponse {
	path := strings.Trim(request.URL, "/") //Remove leading and trailing '/'
	parts := strings.Split(path, "/")
	fmt.Println("Parts:==>", parts)

	switch {
	case request.URL == "/":
		return &HttpResponse{
			Status:  StatusOk,
			Version: Version,
			Headers: make(map[string]string),
		}
	case parts[0] == "echo" && len(parts) == 2:
		return &HttpResponse{
			Status:  StatusOk,
			Version: Version,
			Headers: map[string]string{ContentTypeHeader: ContentTypeText.value},
			Body:    parts[1],
		}
	case parts[0] == "user-agent" && len(parts) == 1:
		return &HttpResponse{
			Status:  StatusOk,
			Version: Version,
			Headers: map[string]string{ContentTypeHeader: ContentTypeText.value},
			Body:    request.Headers["user-agent"],
		}
	case parts[0] == "files" && len(parts) == 2:
		return handleFilesRequest(parts[1])
	default:
		return &HttpResponse{
			Status:  StatusNotFound,
			Version: Version,
			Headers: make(map[string]string),
		}
	}
}
