<!-- SPDX-FileCopyrightText: 2026 The midgard contributors.
     SPDX-License-Identifier: MPL-2.0
-->

Rate Limiter
============

_Rate Limiter_ is a request rate limiting middleware. It uses a _Limiter_ to do
the heavy lifting.

Example
-------

```go
rl := ratelimit.New(
    ratelimit.WithLimiter(locallimit.New(
        locallimit.WithTargetRate(10),
        locallimit.WithSleepInterval(100*time.Millisecond))))
```
