<!-- SPDX-FileCopyrightText: 2025 The midgard contributors.
     SPDX-License-Identifier: MPL-2.0
-->

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
you're adding new features or bugfixes to *midgard*, please include tests.


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
go test -race ./...
```


### Implement your fix or feature

At this point, you're ready to make your changes. Feel free to ask for help.
Be sure to have run the `go fmt` tool to have a unified code style:

```bash
go fmt ./...
```


### Test your feature

After implementing your feature, add tests that cover all major code paths. A
test coverage of 100% is not always possible. We acknowledge that there are hard
to trigger conditions that you might check for, but are not producible in a test
suite, but aim for the best. At least every code path of normal, input-triggered
use should be covered.


### Document your feature

All the good intentions go to waste, if nobody can enjoy the fruits of this labor
due to non-existent (or bad, or wrong) documentation. Please take care that you
include:

- a concise description of your new feature
- generate new or update (in case) the existing examples
- update CHANGELOG.md

The CHANGELOG document contains the changes of the next release and all the
changes of the current major version since x.0.0. On a major release, the
CHANGELOG can be emptied as the older changes are still visible in the history
of the version control system.


### Create a Pull Request

At this point, if your changes look good and tests are passing, you are ready to
create a pull request.

GitHub Actions will run the test suite against the latest Go version. There are
tests that most likely did not run on the developers machine (CodeQL, Trivy).
These tests may produce warnings. Take those warnings seriously even if they
seem harmless. Too many harmless warnings could possibly overlay really serious
ones, so all warnings are to be resolved.


Merging a PR (maintainers only)
-------------------------------

A maintainer can only merge a Pull Request into master if:

- CI is passing,
- approved by another maintainer
- and is up to date with the default branch.

In addition to these automatic checks, the following conditions have to be met:

- Changelog: The changelog has the sections with the changes of the past releases, these
  are immutable less for corrections. There is at least one section explicating the next
  release. All visible changes have to be included in the changelog. Further security
  fixes have to be included here.

Any maintainer is allowed to merge a PR if all of these conditions have been met.


Developer Certificate of Origin (DCO)
-------------------------------------

This project enforces the [Developer Certificate of Origin (DCO)](https://developercertificate.org)
(a copy is included in [DCO.txt](DCO.txt)).
Every commit must be signed off, i.e., the commit message contains a line like:

Signed-off-by: Your Name <you@example.com>

Use `git commit -s` to add this automatically (this DCO “Signed-off-by” is a
legal attestation and is different from cryptographic commit signing such as
GPG/SSH). Pull requests and pushes without proper sign-offs will fail the
compliance checks. You can alias `git commit` command to always include the
flag as follows:

```shell
git config --global alias.ci 'commit -s'
```
