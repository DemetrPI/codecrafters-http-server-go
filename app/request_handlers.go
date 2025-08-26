package main

import (
	"bufio"
	"io"

	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func parseRequest(requestData *bufio.Reader) (*HttpRequest, error, bool) {

	//Read request
	requestLine, err := requestData.ReadString('\n')
	if err == io.EOF {
		return &HttpRequest{}, nil, true
	}

	if err != nil {
		return nil, err, true
	}

	requestParts := strings.Fields(requestLine)
	if len(requestParts) < 3 {
		return nil, fmt.Errorf("invalid request"), true
	}

	requestContent := &HttpRequest{
		Method:  HTTPMethod(requestParts[0]),
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
	//Read body
	contentLength, _ := strconv.Atoi(requestContent.Headers["content-length"])
	if contentLength > 0 {
		body := make([]byte, contentLength)
		_, err := requestData.Read(body)
		if err != nil {
			return nil, err, true
		}
		requestContent.Body = string(body)
	}
	return requestContent, nil, false
}

func sendResponce(conn net.Conn, responce *HttpResponse, request *HttpRequest) error {
	//status line

	reqResp := bufio.NewWriter(conn)
	fmt.Fprint(reqResp, responce.Version, " ", responce.Status, CRLF)

	//write headers
	responce.Headers["content-length"] = fmt.Sprintf("%d", len(responce.Body))
	if val, ok := request.Headers["accept-encoding"]; ok {
		encodings := strings.SplitSeq(val, ",")
		for encoding := range encodings {
			if strings.TrimSpace(encoding) == "gzip" {
				responce.Headers["content-encoding"] = "gzip"
				handleCompression(responce)
				break
			}
		}
	}

	for key, value := range responce.Headers {
		fmt.Fprintf(reqResp, "%s: %s%s", key, value, CRLF)
	}

	//body
	reqResp.WriteString(CRLF + responce.Body)
	reqResp.Flush()

	return nil
}

func (request *HttpRequest) routeRequest() *HttpResponse {
	path := strings.Trim(request.URL, "/") //Remove leading and trailing '/'
	parts := strings.Split(path, "/")
	if len(parts) < 1 {
		return &HttpResponse{
			Status:  StatusBadRequest,
			Version: Version,
			Headers: make(map[string]string),
		}
	}
	endpoint := parts[0]
	if request.Method == HTTPGet {
		switch {
		case request.URL == "/":
			return &HttpResponse{
				Status:  StatusOk,
				Version: Version,
				Headers: make(map[string]string),
			}
		case endpoint == "echo" && len(parts) == 2:
			return &HttpResponse{
				Status:  StatusOk,
				Version: Version,
				Headers: map[string]string{ContentTypeHeader: ContentTypeText.value},
				Body:    parts[1],
			}
		case endpoint == "user-agent" && len(parts) == 1:
			return &HttpResponse{
				Status:  StatusOk,
				Version: Version,
				Headers: map[string]string{ContentTypeHeader: ContentTypeText.value},
				Body:    request.Headers["user-agent"],
			}
		case endpoint == "files" && len(parts) == 2:
			return handleFilesRequestGet(parts[1])
		default:
			return &HttpResponse{
				Status:  StatusNotFound,
				Version: Version,
				Headers: make(map[string]string),
			}
		}
	}
	if request.Method == HTTPPost {
		switch {
		case endpoint == "files" && len(parts) == 2:
			return handleFilesRequestPost(parts[1], []byte(request.Body))
		default:
			return &HttpResponse{
				Status:  StatusNotFound,
				Version: Version,
				Headers: make(map[string]string),
			}
		}
	}
	// Default return
	return &HttpResponse{
		Status:  StatusNotFound,
		Version: Version,
		Headers: make(map[string]string),
	}
}

func handleFilesRequestGet(filename string) *HttpResponse {
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

func handleFilesRequestPost(filename string, body []byte) *HttpResponse {
	args := os.Args
	if len(args) < 2 {
		return &HttpResponse{
			Status:  StatusInternalServerError,
			Version: Version,
			Headers: make(map[string]string),
		}
	}
	directoryName := args[2]

	writeErr := os.WriteFile(directoryName+filename, body, 0644)
	if writeErr != nil {
		return &HttpResponse{
			Status:  StatusInternalServerError,
			Version: Version,
			Headers: make(map[string]string),
		}
	}
	return &HttpResponse{
		Status:  StatusCreated,
		Version: Version,
		Headers: make(map[string]string),
	}
}
