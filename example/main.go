package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"os"

	"github.com/AlphaOne1/midgard"
	"github.com/AlphaOne1/midgard/handler/access_log"
	"github.com/AlphaOne1/midgard/handler/correlation"
	"github.com/AlphaOne1/midgard/handler/cors"
	"github.com/AlphaOne1/midgard/handler/method_filter"
	"github.com/AlphaOne1/midgard/util"
)

//go:embed hello.html
var helloPage []byte

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(helloPage)
}

func main() {
	fmt.Println("Example for midgard usage")

	finalHandler := midgard.StackMiddlewareHandler(
		[]midgard.Middleware{
			correlation.New(),
			access_log.New(),
			cors.New([]string{http.MethodGet}, []string{"*"}),
			util.Must(method_filter.New(method_filter.WithMethods([]string{http.MethodGet}))),
		},
		http.HandlerFunc(HelloHandler),
	)

	http.Handle("/", finalHandler)

	if listenErr := http.ListenAndServe(":8080", nil); listenErr != nil {
		fmt.Println("got error listening:", listenErr)
		os.Exit(1)
	}

	fmt.Println("finished")
	os.Exit(0)
}
