Correlation ID Middleware
=========================

The correlation ID middleware adds a `X-Correlation-ID` header to incoming
requests, if they do not contain one already.
The correlation ID can be used track the control flow in systems of
microservices.

Example
-------

```go
finalHandler := midgard.StackMiddlewareHandler(
    []midgard.Middleware{
        correlation.New(),
    },
    http.HandlerFunc(HelloHandler),
)
```
