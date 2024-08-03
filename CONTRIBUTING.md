Contributing
============

Thanks for your interest in contributing to *midgard*. Please take a moment to
review this document __before submitting a pull request__.

Pull requests
-------------

__Please ask first before starting work on any significant new features.__

It's never a fun experience to have your pull request declined after investing a
lot of time and effort into a new feature. To avoid this from happening, we
request that contributors create
[a feature request](https://github.com/AlphaOne1/midgard/discussions/new?category=ideas)
to first discuss any new ideas. Your ideas and suggestions are welcome!

Please ensure that the tests are passing when submitting a pull request. If
you're adding new features to ActiveAdmin, please include tests.

Where do I go from here?
------------------------

For any questions, support, or ideas, etc.
[please create a GitHub discussion](https://github.com/AlphaOne1/midgard/discussions/new).
If you've noticed a bug, [please submit an issue][new issue].

### Fork and create a branch

If this is something you think you can fix, then [fork midgard] and create a
branch with a descriptive name.

### Get the test suite running

Make sure you're using a recent Go version.

You can run the test suite from the base folder using the following command:

```bash
go test ./...
```

### Implement your fix or feature

At this point, you're ready to make your changes. Feel free to ask for help.
Be sure to have run the go fmt tool, to have a unified code style:

```bash
go fmt ./...
```

### Create a Pull Request

At this point, if your changes look good and tests are passing, you are ready to
create a pull request.

Github Actions will run the test suite against the latest Go version. There are
tests that most likey did not run in the developers machine (CodeQL, Trivy). These
tests may produce warnings. Take those warnings serious even if they seem harmless.
Too many harmless warnings could possibly overlay really serious ones, so all
warnings are to be resolved.

Merging a PR (maintainers only)
-------------------------------

A Pull Request can only be merged into master by a maintainer if:

* CI is passing,
* approved by another maintainer
* and is up to date with the default branch.

Any maintainer is allowed to merge a PR if all of these conditions ae met.
