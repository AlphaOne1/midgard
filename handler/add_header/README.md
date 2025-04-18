Header Adding Middleware
=========================

The header adding middleware is used to add custom headers to responses. Note that headers are generally not
required to be unique. So repetitions are allowed and this middleware will not check for duplicates.

Example
-------

The header adding should be configured with actual headers to add. It will work as a NoOp, if none are provided.

```go
finalHandler := midgard.StackMiddlewareHandler(
    []midgard.Middleware{
        add_header.New(
            add_header.WithHeaders([][2]string{
                "X-Test": "TestHeaderValue",
            }),
        ),
    },
    http.HandlerFunc(HelloHandler),
)
```
