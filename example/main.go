package main

import (
	_ "embed"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/AlphaOne1/midgard"
	"github.com/AlphaOne1/midgard/defs"
	"github.com/AlphaOne1/midgard/handler/access_log"
	"github.com/AlphaOne1/midgard/handler/correlation"
	"github.com/AlphaOne1/midgard/handler/cors"
	"github.com/AlphaOne1/midgard/handler/method_filter"
	"github.com/AlphaOne1/midgard/util"
)

//go:embed hello.html
var helloPage []byte

// HelloHandler is an intentionally simple http handler.
func HelloHandler(w http.ResponseWriter, _ *http.Request) {
	if _, err := w.Write(helloPage); err != nil {
		slog.Error("could not write hello page", slog.String("error", err.Error()))
	}
}

func main() {
	fmt.Println("Example for midgard usage")

	// generate a handler that is prepended with the given middlewares
	finalHandler := midgard.StackMiddlewareHandler(
		[]defs.Middleware{
			correlation.New(),
			util.Must(access_log.New(
				access_log.WithLogLevel(slog.LevelDebug))),
			util.Must(cors.New(
				cors.WithHeaders(cors.MinimumAllowHeaders()),
				cors.WithMethods([]string{http.MethodGet}),
				cors.WithOrigins([]string{"*"}))),
			util.Must(method_filter.New(
				method_filter.WithMethods([]string{http.MethodGet}))),
		},
		http.HandlerFunc(HelloHandler),
	)

	// register the newly generated handler for the / endpoint
	http.Handle("/", finalHandler)

	// start the server
	if listenErr := http.ListenAndServe(":8080", nil); listenErr != nil {
		fmt.Println("got error listening:", listenErr)
		os.Exit(1)
	}

	// normally the execution flow cannot reach these lines
	fmt.Println("finished")
	os.Exit(0)
}
