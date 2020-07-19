package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/sago35/tinygo-examples/wioterminal/wttw/tweet"
)

var (
	index int
)

var (
	twitter *tw
	prev    string
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
	if false {
		handlerU_user(w, r)
	} else {
		handlerU_search(w, r)
	}
}

func handlerU_user(w http.ResponseWriter, r *http.Request) {
	log.Printf("handlerU\n")
	//searchWord := r.URL.Query().Get("keyword")
	searchWord := "sago35tk"

	var err error
	since := int64(0)
	max := int64(0)

	sinceStr := r.URL.Query().Get("since")
	if sinceStr != "" {
		fmt.Printf("sinceStr: %q\n", sinceStr)
		since, err = strconv.ParseInt(sinceStr, 10, 64)
		if err != nil {
			log.Printf("since error: %s", err.Error())
		}
	}

	maxStr := r.URL.Query().Get("max")
	if maxStr != "" {
		fmt.Printf("maxStr: %q\n", maxStr)
		max, err = strconv.ParseInt(maxStr, 10, 64)
		if err != nil {
			log.Printf("max error: %s", err.Error())
		}
		max--
	}

	tweets, err := twitter.UserTimeline(searchWord, since, max)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	if len(tweets) == 0 {
		fmt.Fprintf(w, prev)
		return
	}

	idx := 0
	if sinceStr != "" {
		idx = len(tweets) - 1
	}

	t := twitter.convert(tweets[idx])
	b, err := json.Marshal(t)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	log.Printf("%s - %d %s\n", r.URL.String(), t.Id, t.CreatedAt)
	fmt.Fprintf(w, string(b))
	prev = string(b)
}

func handlerU_search(w http.ResponseWriter, r *http.Request) {
	log.Printf("handlerU_search\n")
	//searchWord := r.URL.Query().Get("keyword")
	searchWord := "#umedago"
	//searchWord := "#tinygo"

	var err error
	since := int64(0)
	max := int64(0)

	sinceStr := r.URL.Query().Get("since")
	if sinceStr != "" {
		fmt.Printf("sinceStr: %q\n", sinceStr)
		since, err = strconv.ParseInt(sinceStr, 10, 64)
		if err != nil {
			log.Printf("since error: %s", err.Error())
		}
	}

	maxStr := r.URL.Query().Get("max")
	if maxStr != "" {
		fmt.Printf("maxStr: %q\n", maxStr)
		max, err = strconv.ParseInt(maxStr, 10, 64)
		if err != nil {
			log.Printf("max error: %s", err.Error())
		}
		max--
	}

	tweets, err := twitter.GetSearch(searchWord, since, max)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	if len(tweets) == 0 {
		fmt.Fprintf(w, prev)
		return
	}

	idx := 0
	if sinceStr != "" {
		idx = len(tweets) - 1
	}

	t := twitter.convert(tweets[idx])
	b, err := json.Marshal(t)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	log.Printf("%s - %d %s\n", r.URL.String(), t.Id, t.CreatedAt)
	fmt.Fprintf(w, string(b))
	prev = string(b)
}

func handlerE(w http.ResponseWriter, r *http.Request) {
	// black 1 x 240 px
	img := make([]byte, 240*2*2)
	for i := range img {
		img[i] = 0x7F
	}
	w.Write(img)

	//w.Write([]byte{0, 1, 2, 3, 4, 5})
	log.Printf(r.URL.String())
	//fmt.Fprintf(w, "handlerE %s", r.URL.Query().Get("url"))
}

// for connection check
func handlerC(w http.ResponseWriter, r *http.Request) {
	log.Printf(r.URL.String())
	fmt.Fprintf(w, "ok")
}

func main() {
	var err error
	twitter, err = newTw()
	if err != nil {
		log.Fatal(err.Error())
	}

	http.HandleFunc("/", handler)
	//http.HandleFunc("/u", handler)
	http.HandleFunc("/u/", handlerU)
	http.HandleFunc("/e", handlerE)
	http.HandleFunc("/c", handlerC)

	fmt.Printf("start server\n")
	http.ListenAndServe(":8081", nil)
}
