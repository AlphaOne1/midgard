<!-- SPDX-FileCopyrightText: 2026 The midgard contributors.
     SPDX-License-Identifier: MPL-2.0
-->

HTPassWD Authorizer
==============

The htpasswd authorizer is a simple user-password matcher. It reads its
configuration from a htpasswd formatted file.

Example
-------

```go
handler := midgard.StackMiddlewareHandler(
    []defs.Middleware{
        helper.Must(basicauth.New(
            basicauth.WithAuthenticator(helper.Must(
                htpasswdauth.New(htpasswdauth.WithAuthFile('./testwd')))),
            basicauth.WithRealm("testrealm"))),
    },
    http.HandlerFunc(helper.DummyHandler),
)
```

Be aware that having password hashes accessible to the program potentially
exposes them to attackers. Use strong hashing like bcrypt to counteract brute
force attacks on the hashes or rainbow tables.
