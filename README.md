<p align="center">
    <img src="midgard_logo.svg" width="25%" alt="Logo"><br>
    <a href="https://github.com/AlphaOne1/midgard/actions/workflows/test.yml">
        <img src="https://github.com/AlphaOne1/midgard/actions/workflows/test.yml/badge.svg"
             alt="Pipeline Result">
    </a>
    <a href="https://goreportcard.com/report/github.com/AlphaOne1/midgard">
        <img src="https://goreportcard.com/badge/github.com/AlphaOne1/midgard"
             alt="Go Report Card">
    </a>
    <a href="https://scorecard.dev/viewer/?uri=github.com/AlphaOne1/midgard">
        <img src="https://api.scorecard.dev/projects/github.com/AlphaOne1/midgard/badge"
             alt="OpenSSF Scorecard">
    </a>
</p>

midgard
=======

*midgard* is a collection of Golang http middlewares and helper functionality
to use them more elegantly.

Usage
-----

*midgard* defines a type `Middleware` that is just a convenience to not always
having to write the full definition of what is commonly known as http middlware.

```go
type Middleware func(http.Handler) http.Handler
```

To ease the pain of stacking different middlwares *midgard* offsers two functions
to facilitate it. `StackMiddlewareHandler` stacks the given slice of middlewares
on top of each other and finally calls the given handler. It generates a new handler
that has all the given middlewares prepended:

```go
finalHandler := midgard.StackMiddlewareHandler(
    []midgard.Middleware{
        correlation.New(),
        util.Must(access_log.New(
            access_log.WithLogLevel(slog.LevelDebug))),
        util.Must(cors.New(
            cors.WithHeaders(cors.MinimumAllowedHeaders()),
            cors.WithMethods([]string{http.MethodGet}),
            cors.WithOrigins([]string{"*"}))),
        util.Must(method_filter.New(
            method_filter.WithMethods([]string{http.MethodGet}))),
        },
    http.HandlerFunc(HelloHandler),
)
```

`StackMiddleware` does basically the same, but without having given a handler.
It generates a new middleware:

```go
newMiddleware:= midgard.StackMiddleware(
    []midgard.Middleware{
        correlation.New(),
        util.Must(access_log.New(
            access_log.WithLogLevel(slog.LevelDebug))),
        util.Must(cors.New(
            cors.WithHeaders(cors.MinimumAllowedHeaders()),
            cors.WithMethods([]string{http.MethodGet}),
            cors.WithOrigins([]string{"*"}))),
        util.Must(method_filter.New(
            method_filter.WithMethods([]string{http.MethodGet}))),
    })
```

The native solution for this would be to nest the calls to the middleware like this:

```go
finalHandler := correlation.New()(
                    util.Must(access_log.New(
                        access_log.WithLogLevel(slog.LevelDebug)))(
                        util.Must(cors.New(
                            cors.WithHeaders(cors.MinimumAllowedHeaders()),
                            cors.WithMethods([]string{http.MethodGet}),
                            cors.WithOrigins([]string{"*"})))(
                            util.Must(method_filter.New(
                                method_filter.WithMethods([]string{http.MethodGet}))))))
```

As you see, depending on the number of middlewares, that can be quite confusing.
Further one cannot _easily_ dynamially add or remove middlewares.
