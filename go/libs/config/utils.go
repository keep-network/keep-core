// Copyright 2018 The Keep Authors.  See LICENSE.md for details.
package config

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
	"time"
	"strings"
)

type APIError struct {
	ErrorMessage string `json:"error_message"`
	HTTPStatus   int    `json:"http_status"`
}

type HttpErrorHandler struct {
	Caller   string
	Response http.ResponseWriter
	Request  *http.Request
}

const (
	ErrorActionErr = iota
	ErrorActionWarn
	ErrorActionDebug
	ErrorActionInfo
)

func NewHttpErrorHandle(caller string, response http.ResponseWriter, request *http.Request) *HttpErrorHandler {
	return &HttpErrorHandler{caller, response, request}
}

// HandleError locally, according to the action passed to h.Handle, and then serialized
// in json and sent to the remote address via http, then returns true.
// Otherwise, if there is no error, h.Handle returns false
func (h *HttpErrorHandler) Handle(err error, httpStatus int, action int) bool {
	if err != nil {
		_, filepath, line, _ := runtime.Caller(1)
		_, file := path.Split(filepath)
		Error.Printf("HttpErrorHandler()->[file:%s line:%d]: %s", file, line, err.Error())
		apiErr := &APIError{
			ErrorMessage: err.Error(),
			HTTPStatus:   httpStatus,
		}
		serialErr, _ := json.Marshal(&apiErr)
		h.Response.Header().Set("Content-Type", "application/json")
		h.Response.WriteHeader(httpStatus)
		io.WriteString(h.Response, string(serialErr))
	}
	return err != nil
}

// HandlePanic _Never_ returns on error, instead it panics
func FromLineOfFile() string {
		_, filepath, line, _ := runtime.Caller(1)
		_, file := path.Split(filepath)
		return fmt.Sprintf("[file:%s line:%d]", file, line)
}

// HandlePanic _Never_ returns an error, instead it panics
func HandlePanic(err error) {
	if err != nil {
		_, filePath, lineNo, _ := runtime.Caller(1)
		_, fileName := path.Split(filePath)
		msg := fmt.Sprintf("[file:%s line:%d]: %s", fileName, lineNo, err.Error())
		panic(msg)
	}
}

func HandleError(err error, action int) bool {
	if err != nil {
		_, filepath, line, _ := runtime.Caller(1)
		_, file := path.Split(filepath)
		switch action {
		case ErrorActionErr:
			Error.Printf("[file:%s line:%d]: %s", file, line, err.Error())
			break
		case ErrorActionWarn:
			Error.Printf("[file:%s line:%d]: %s", file, line, err.Error())
			break
		case ErrorActionDebug:
			Error.Printf("[file:%s line:%d]: %s", file, line, err.Error())
			break
		case ErrorActionInfo:
			Error.Printf("[file:%s line:%d]: %s", file, line, err.Error())
			break
		}
	}
	return err != nil
}

func WriteFile(filename string, source io.Reader) error {
	writer, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer writer.Close()
	io.Copy(writer, source)
	return nil
}

// This is neat: https://coderwall.com/p/cp5fya/measuring-execution-time-in-go
func TimeTrack(start time.Time, name string) {
	if Config.LogTimeTrack == true {
		elapsed := time.Since(start)
		Info.Printf("%s took %s", name, elapsed)
	}
}

// pad str with padWith count times to right
func PadRight(str string, padWith string, length int) string {
	count := length - len(str)
	if count < 0 {
		count = 0
	}
	return str + strings.Repeat(padWith, count)
}

func InSlice(slice []string, searchFor string) (found bool) {
	for _, v := range slice {
		if searchFor == v {
			found = true
		}
	}
	return found
}
