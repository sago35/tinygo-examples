//go:build !wioterminal
// +build !wioterminal

package main

import (
	"io"
	"net/http"
)

func _init() error {
	message = "hello from Go"
	return nil
}

func post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	return http.Post(url, contentType, body)
}
