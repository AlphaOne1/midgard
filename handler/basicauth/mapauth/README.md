<!-- SPDX-FileCopyrightText: 2026 The midgard contributors.
     SPDX-License-Identifier: MPL-2.0
-->

Map Authorizer
==============

The map authorizer is a simple user-password matcher. It is configured inside
the program. It is, next to "always true", one of the most simple authorizers
possible.

Example
-------

```go
handler := midgard.StackMiddlewareHandler(
    []defs.Middleware{
        helper.Must(basicauth.New(
            basicauth.WithAuthenticator(helper.Must(
                mapauth.New(mapauth.WithAuths(map[string]string{
                    "user0": "pass0",
                    "user1": "pass1",
                })))),
            basicauth.WithRealm("testrealm"))),
    },
    http.HandlerFunc(helper.DummyHandler),
)
```

Be aware that writing credentials inside program code is _not_ advisable and is
just used here to illustrate the usage.
