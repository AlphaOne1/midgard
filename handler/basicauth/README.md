<!-- SPDX-FileCopyrightText: 2026 The midgard contributors.
     SPDX-License-Identifier: MPL-2.0
-->

Basic Authentication Middleware
===============================

Basic Authentication is an authentication method where the client sends its
credentials unencrypted to the server for validation. This middleware implements
this process. A detailed description can also be found
[here](https://en.wikipedia.org/wiki/Basic_access_authentication).

The credentials can be stored in different forms. This module has an interface
for authenticators that provide the final verification, and concentrates on the
protocol between the client and the server. Further, a simple example
authenticator is provided that is to be configured with the allowed credentials.

Example
-------

```go
finalHandler := midgard.StackMiddlewareHandler(
    []midgard.Middleware{
        helper.Must(basicauth.New(
            basicauth.WithRealm("example realm"),
            basicauth.WithAuthenticator(helper.Must(
                mapauth.New(
                    mapauth.WithAuths(map[string]string{
                        "user0": "pass0",
                        "user1": "pass1",
                    }),
                ),
            )),
        )),
    },
    http.HandlerFunc(HelloHandler),
)
```

If no realm is specified using `WithRealm` the default `Restricted` is used.
Not providing an authenticator is an error condition.
