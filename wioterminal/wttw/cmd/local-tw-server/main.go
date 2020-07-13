package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	b, err := json.Marshal(t[0])
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	log.Printf(r.URL.String())
	fmt.Fprintf(w, string(b))
}

func handlerE(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte{0, 1, 2, 3, 4, 5})
	log.Printf(r.URL.String())
	fmt.Fprintf(w, "handlerE %s", r.URL.Query().Get("url"))
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/u", handlerU)
	http.HandleFunc("/e", handlerE)

	fmt.Printf("start server\n")
	http.ListenAndServe(":8080", nil)
}
