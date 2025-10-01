<!-- SPDX-FileCopyrightText: 2025 The midgard contributors.
     SPDX-License-Identifier: MPL-2.0
-->

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
        accesslog.New(),
    },
    http.HandlerFunc(HelloHandler),
)
```

It can also be configured with a custom logger and level.

```go
finalHandler := midgard.StackMiddlewareHandler(
    []midgard.Middleware{
        accesslog.New(
            accesslog.WithLogger(someOtherLogger),
            accesslog.WithLogLevel(slog.LevelDebug),
        ),
    },
    http.HandlerFunc(HelloHandler),
)
```
