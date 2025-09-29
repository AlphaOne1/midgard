<!-- SPDX-FileCopyrightText: 2025 The midgard contributors.
     SPDX-License-Identifier: MPL-2.0
-->

<!-- markdownlint-disable MD013 MD033 MD041 -->
<p align="center">
    <img src="midgard_logo.svg" width="25%" alt="Logo"><br>
    <a href="https://github.com/AlphaOne1/midgard/actions/workflows/test.yml"
       rel="external noopener noreferrer"
       target="_blank">
        <img src="https://github.com/AlphaOne1/midgard/actions/workflows/test.yml/badge.svg"
             alt="Test Pipeline Result">
    </a>
    <a href="https://github.com/AlphaOne1/midgard/actions/workflows/codeql.yml"
       rel="external noopener noreferrer"
       target="_blank">
        <img src="https://github.com/AlphaOne1/midgard/actions/workflows/codeql.yml/badge.svg"
             alt="CodeQL Pipeline Result">
    </a>
    <a href="https://github.com/AlphaOne1/midgard/actions/workflows/security.yml"
       rel="external noopener noreferrer"
       target="_blank">
        <img src="https://github.com/AlphaOne1/midgard/actions/workflows/security.yml/badge.svg"
             alt="Security Pipeline Result">
    </a>
    <a href="https://goreportcard.com/report/github.com/AlphaOne1/midgard"
       rel="external noopener noreferrer"
       target="_blank">
        <img src="https://goreportcard.com/badge/github.com/AlphaOne1/midgard"
             alt="Go Report Card">
    </a>
    <a href="https://app.codecov.io/gh/AlphaOne1/midgard"
       rel="external noopener noreferrer"
       target="_blank">
        <img src="https://codecov.io/gh/AlphaOne1/midgard/graph/badge.svg"
             alt="Code Coverage">
    </a>
    <a href="https://coderabbit.ai"
       rel="external noopener noreferrer"
       target="_blank">
       <img src="https://img.shields.io/coderabbit/prs/github/AlphaOne1/midgard"
            alt="CodeRabbit Reviews">
    </a>
    <a href="https://www.bestpractices.dev/projects/9251"
       rel="external noopener noreferrer"
       target="_blank">
        <img src="https://www.bestpractices.dev/projects/9251/badge"
             alt="OpenSSF Best Practises">
    </a>
    <a href="https://scorecard.dev/viewer/?uri=github.com/AlphaOne1/midgard"
       rel="external noopener noreferrer"
       target="_blank">
        <img src="https://api.scorecard.dev/projects/github.com/AlphaOne1/midgard/badge"
             alt="OpenSSF Scorecard">
    </a>
    <a href="https://api.reuse.software/info/github.com/AlphaOne1/midgard"
       rel="external noopener noreferrer"
       target="_blank">
        <img src="https://api.reuse.software/badge/github.com/AlphaOne1/midgard"
            alt="REUSE compliance">
    </a>
    <a href="https://app.fossa.com/projects/git%2Bgithub.com%2FAlphaOne1%2Fmidgard?ref=badge_shield&issueType=license"
       rel="external noopener noreferrer"
       target="_blank">
        <img src="https://app.fossa.com/api/projects/git%2Bgithub.com%2FAlphaOne1%2Fmidgard.svg?type=shield&issueType=license"
            alt="FOSSA License Status">
    </a>
    <a href="https://app.fossa.com/projects/git%2Bgithub.com%2FAlphaOne1%2Fmidgard?ref=badge_shield&issueType=security"
       rel="external noopener noreferrer"
       target="_blank">
        <img src="https://app.fossa.com/api/projects/git%2Bgithub.com%2FAlphaOne1%2Fmidgard.svg?type=shield&issueType=security"
             alt="FOSSA Security Status">
    </a>
    <a href="https://pkg.go.dev/github.com/AlphaOne1/midgard"
       rel="external noopener noreferrer"
       target="_blank">
        <img src="https://pkg.go.dev/badge/github.com/AlphaOne1/midgard.svg"
             alt="Go Reference">
    </a>
</p>
<!-- markdownlint-enable MD013 MD033 MD041 -->

midgard
=======

*midgard* is a collection of Golang http middlewares and helper functionality
to use them more elegantly.

Usage
-----

*midgard* defines a type `Middleware` that is just a convenience to not always
having to write the full definition of what is commonly known as http middleware.

```go
type Middleware func(http.Handler) http.Handler
```

To ease the pain of stacking different middlewares, *midgard* offers two functions
to facilitate it. `StackMiddlewareHandler` stacks the given slice of middlewares
on top of each other and finally calls the given handler. It generates a new handler
that has all the given middlewares prepended:

```go
finalHandler := midgard.StackMiddlewareHandler(
    []midgard.Middleware{
        util.Must(correlation.New()),
        util.Must(access_log.New(
            access_log.WithLogLevel(slog.LevelDebug))),
        util.Must(cors.New(
            cors.WithHeaders(cors.MinimumAllowedHeaders()),
            cors.WithMethods([]string{http.MethodGet}),
            cors.WithOrigins([]string{"*"}))),
        util.Must(method_filter.New(
            method_filter.WithMethods([]string{http.MethodGet}))),
        },
    http.HandlerFunc(HelloHandler),
)
```

`StackMiddleware` does basically the same, but without having given a handler.
It generates a new middleware:

```go
newMiddleware:= midgard.StackMiddleware(
    []midgard.Middleware{
        util.Must(correlation.New()),
        util.Must(access_log.New(
            access_log.WithLogLevel(slog.LevelDebug))),
        util.Must(cors.New(
            cors.WithHeaders(cors.MinimumAllowedHeaders()),
            cors.WithMethods([]string{http.MethodGet}),
            cors.WithOrigins([]string{"*"}))),
        util.Must(method_filter.New(
            method_filter.WithMethods([]string{http.MethodGet}))),
    })
```

The native solution for this would be to nest the calls to the middleware like this:

```go
finalHandler := util.Must(correlation.New())(
                    util.Must(access_log.New(
                        access_log.WithLogLevel(slog.LevelDebug)))(
                        util.Must(cors.New(
                            cors.WithHeaders(cors.MinimumAllowedHeaders()),
                            cors.WithMethods([]string{http.MethodGet}),
                            cors.WithOrigins([]string{"*"})))(
                            util.Must(method_filter.New(
                                method_filter.WithMethods([]string{http.MethodGet}))))))
```

As you see, depending on the number of middlewares, that can be quite confusing.
Further one cannot *easily* dynamically add or remove middlewares.
