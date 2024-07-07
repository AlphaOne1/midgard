Rate Limiter
============

_Rate Limiter_ is a request rate limiting middleware. It uses a *Limiter* to do
the heavy lifting.

Example
-------

```go
rl := New(WithLimiter(local_limit.New(
	local_limit.WithTargetRate(10))),
	local_limit.WithSleepInterval(100*time.Millisecond),
)
```