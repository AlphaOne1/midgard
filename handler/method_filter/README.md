Methods Filter Middleware
=========================

The methods filter middleware can be used to limit the allowed methods for a
certain endpoint.

HTTP has several methods, e.g. GET, PUT, DELETE, that can be used in requests
to an endpoint. The method gives a coarse hint of what is intended with the
call. If requests to certain methods are dangerous this middleware could be used
to block those using a whitelist.

Example
-------

```go
	finalHandler := midgard.StackMiddlewareHandler(
		[]midgard.Middleware{
			util.Must(method_filter.New(method_filter.WithMethods([]string{"GET"}))),
		},
		http.HandlerFunc(HelloHandler),
	)
```
