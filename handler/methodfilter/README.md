<!-- SPDX-FileCopyrightText: 2026 The midgard contributors.
     SPDX-License-Identifier: MPL-2.0
-->

Methods Filter Middleware
=========================

The methods filter middleware can be used to limit the allowed methods for a
certain endpoint.

HTTP has several methods, e.g. GET, PUT and DELETE, that can be used in requests
to an endpoint. The method gives a coarse hint of the intention of that request.
If requests utilizing methods, for any reason are undesirable, dangerous or just
to be blocked in the specific use case, this middleware could be used to block
those requests using a whitelist.

Blocked requests receive an HTTP status of 405 - method not allowed and a
corresponding content as text/plain. Handlers down the handler stack are not called
in case of a blocked request.

Requests with methods contained in the white list, pass this filter without further
action.

Example
-------

```go
finalHandler := midgard.StackMiddlewareHandler(
    []midgard.Middleware{
        helper.Must(methodfilter.New(
            methodfilter.WithMethods([]string{http.MethodGet}))),
    },
    http.HandlerFunc(HelloHandler),
)
```
