Local Limit
===========

_Local Limit_ is a rate limiter that works on instance level. It is not intended
to limit the global request rate to numerous instances of a service. This can be
achieved using a somehow synchronized rate limiting algorithm, e.g. via a redis.

_Local Limit_ allows a certain number of requests (called drops here) per second
to pass. Inside, it uses a dripping like algorithm to add possible requests to a
queue. Drops are added in certain intervals. Drop rates less than one per second
are possible, as fractional drops are stored for the next iteration.

Every incoming request takes one (whole) drop out of the queue of drops, until
it is empty. If there a no drops available, a request waits for an arbitrary
time for a new drop. If none arrives, the request is rejected.

Drops can accumulate to a specified maximum. So services will not be overwhelmed
if after a longer period without requests, the requests start again.

As there is an internal go routing caring to add the drops, a Stop() function is
provided to gracefully shut the limiter down. This shutdown is asynchronous and
will occur in the next iteration for the drops.

Example
-------

```go
limiter := util.Must(New(
    WithTargetRate(v.TargetRate),
    WithDefaultSleepInterval(v.SleepTime)))

if limiter.Limit() {
    doWork()
}
```