Access Logging Middleware
=========================

The access logging middleware logs each request with the following information:

- correlationID
- client address
- HTTP method
- path

Example
-------

The access logging will use the default logger `slog.Default()` and INFO log level.

```go
finalHandler := midgard.StackMiddlewareHandler(
    []midgard.Middleware{
        access_log.New(),
    },
    http.HandlerFunc(HelloHandler),
)
```

It can also be configured with a custom logger and level.

```go
finalHandler := midgard.StackMiddlewareHandler(
    []midgard.Middleware{
        access_log.New(
            access_log.WithLogger(someOtherLogger),
            access_log.WithLogLevel(slog.LevelDebug),
        ),
    },
    http.HandlerFunc(HelloHandler),
)
```
