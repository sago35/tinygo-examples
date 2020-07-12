package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
    [
        {
            UserName: "Hiroshi SAKURAI",
            ScreenName: "anolivetree",
            CreatedAt: "Sat Jul 11 12:38:52 +0000 2020",
            FullText: "TinyGoを調べてみたのだけれど、簡単な組み込み機器になら十分つかえそう。あとはこれに優先度順のスケジューラーを載せればTinyGoで良いというケースはかなりあると思う。",
            FavoriteCount: 2,
            RetweetCount: 0,
            IsRetweet: false,
        }
    ]
`)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
