Access Logging Middleware
=========================

The access logging middleware logs each request with the following information:

  - correlationID
  - client address
  - HTTP method
  - path

Example
-------

```go
finalHandler := midgard.StackMiddlewareHandler(
    []midgard.Middleware{
        access_log.New(),
    },
    http.HandlerFunc(HelloHandler),
)
```