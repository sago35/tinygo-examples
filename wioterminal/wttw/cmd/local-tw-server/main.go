package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/sago35/tinygo-examples/wioterminal/wttw/tweet"
)

var (
	index int
)

func handler(w http.ResponseWriter, r *http.Request) {
    t := tweet.S
	b, err := json.Marshal(t[index])
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	log.Printf(r.URL.String())
	fmt.Fprintf(w, string(b))
	index = (index + 1) % len(t)
}

func handlerU(w http.ResponseWriter, r *http.Request) {
    t := tweet.S
    t[index].UserName = strings.ToUpper(t[index].UserName)
	b, err := json.Marshal(t[index])
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	log.Printf(r.URL.String())
	fmt.Fprintf(w, string(b))
	index = (index + 1) % len(t)
}


func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/u", handlerU)

	fmt.Printf("start server\n")
	http.ListenAndServe(":8080", nil)
}
