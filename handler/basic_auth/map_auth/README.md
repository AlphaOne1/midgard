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
		util.Must(New(
			WithAuthenticator(map_auth.New(map_auth.WithAuths(map[string]string{
                "user0": "pass0",
                "user1": "pass1",
            }))),
			WithRealm("testrealm"))),
	},
	http.HandlerFunc(util.DummyHandler),
)
```

Be aware that writing credentials inside of program code is _not_ advisable and
is just used here to illustrate the usage.