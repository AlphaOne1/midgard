package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"os"

	"github.com/AlphaOne1/midgard"
	"github.com/AlphaOne1/midgard/handler"
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
			handler.Correlation,
			handler.NewMethodsFilter([]string{"GET"}),
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
