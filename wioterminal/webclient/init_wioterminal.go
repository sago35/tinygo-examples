//go:build wioterminal
// +build wioterminal

package main

import (
	"io"
	"time"

	"github.com/sago35/tinygo-examples/wioterminal/initialize"
	"tinygo.org/x/drivers/net/http"
)

var (
	ssid     string
	password string
)

func _init() error {
	message = "hello from TinyGo"
	_, err := initialize.Wifi(ssid, password, time.Second*10)
	return err
}

func post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	return http.Post(url, contentType, body)
}
