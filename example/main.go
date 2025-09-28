// SPDX-FileCopyrightText: 2025 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/AlphaOne1/midgard/handler/correlation"

	"github.com/AlphaOne1/midgard"
	"github.com/AlphaOne1/midgard/defs"
	"github.com/AlphaOne1/midgard/handler/access_log"
	"github.com/AlphaOne1/midgard/handler/cors"
	"github.com/AlphaOne1/midgard/handler/method_filter"
	"github.com/AlphaOne1/midgard/util"
)

//go:embed hello.html
var helloPage []byte

// HelloHandler is an intentionally simple http handler.
func HelloHandler(w http.ResponseWriter, _ /* r */ *http.Request) {
	if _, err := w.Write(helloPage); err != nil {
		slog.Error("could not write hello page", slog.String("error", err.Error()))
	}
}

func main() {
	fmt.Println("Example for midgard usage")

	// generate a handler that is prepended with the given middlewares
	finalHandler := midgard.StackMiddlewareHandler(
		[]defs.Middleware{
			util.Must(correlation.New()),
			util.Must(access_log.New(
				access_log.WithLogLevel(slog.LevelDebug))),
			util.Must(cors.New(
				cors.WithHeaders(append(cors.MinimumAllowHeaders(), "X-Correlation-ID")),
				cors.WithMethods([]string{http.MethodGet}),
				cors.WithOrigins([]string{"*"}))),
			util.Must(method_filter.New(
				method_filter.WithMethods([]string{http.MethodGet}))),
		},
		http.HandlerFunc(HelloHandler),
	)

	// register the newly generated handler for the / endpoint
	http.Handle("/", finalHandler)

	server := &http.Server{Addr: "localhost:8080", Handler: nil}

	go func() {
		time.Sleep(1 * time.Second)
		_ = server.Shutdown(context.Background())
	}()

	// start the server
	if listenErr := server.ListenAndServe(); listenErr != nil && !errors.Is(listenErr, http.ErrServerClosed) {
		fmt.Println("got error listening:", listenErr)
		os.Exit(1)
	}

	fmt.Println("finished")
}
