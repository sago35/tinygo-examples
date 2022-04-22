package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")

	body := r.Body
	defer body.Close()

	b, err := ioutil.ReadAll(body)
	if err != nil {
		panic(err)
	}
	log.Printf("%s\n", string(b))
}

func main() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
	}
}
