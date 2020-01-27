package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/clear-street/backend-screening-parthingle/src/handler"
)

func port() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	return ":" + port
}

func echo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hit echo")
	message := r.URL.Query()["message"][0]

	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprintf(w, message)
}
func main() {
	fmt.Println("Listening on " + port())
	http.HandleFunc("/v1/echo", echo)
	http.HandleFunc("/v1/trades", handler.TradesHandlerFunc)
	http.HandleFunc("/v1/trades/", handler.TradeHandlerFunc)
	http.ListenAndServe(port(), nil)

}
