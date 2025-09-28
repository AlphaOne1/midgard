<!-- SPDX-FileCopyrightText: 2025 The midgard contributors.
     SPDX-License-Identifier: MPL-2.0
-->

CORS Middleware
===============

CORS (Cross Origin Resource Sharing) is a technique that allows a microservice
in combination with a client to control the way a microservice can be called.
For this, the microservice provides information from which websites calls to the
microservice can be sent. The client on the other side checks this information
and prevents websites from calling the microservice, that are not allowed to do
so.
A more in depth description can be found
[here](https://en.wikipedia.org/wiki/Cross-origin_resource_sharing)

This middleware intercepts the `OPTIONS` method to provide the CORS information
to clients. It gets a list of allowed methods and headers and will prevent
requests to pass through, that do not fulfill the requirements.

Example
-------

```go
finalHandler := midgard.StackMiddlewareHandler(
    []midgard.Middleware{
        util.Must(cors.New(
            cors.WithHeaders(cors.MinimumAllowHeaders()),
            cors.WithMethods([]string{http.MethodGet}),
            cors.WithOrigins([]string{"*"}))),
    },
    http.HandlerFunc(HelloHandler),
)
```

If no headers are specified, all headers are allowed. A minimal set of headers
is provided via `cors.MinimumAllowHeaders`.
Similar, if no methods are specified, all methods are allowed.
If at least one of the allowed origins is `*` or nothing is specified, the
allowed origins are set to just contain `*`.
Thus, a CORS middleware that is not parametrized, will allow all requests to
pass and not filter anything. It just intercepts the `OPTIONS` method.
