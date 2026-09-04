<!-- SPDX-FileCopyrightText: 2026 The midgard contributors.
     SPDX-License-Identifier: MPL-2.0
-->

Release Next
============

- changed to self-maintained GoReportCard
- dependency updates

Release 0.3.0
=============

- updated the minimum required Go version to 1.27
- dependency updates and used now included uuid package instead of an external one
- harden GitHub Actions against injections via {{...}}
- reworked CORS middleware. The CORS now does _not_ filter requests for allowed
  methods and headers. The methods filtering can be achieved with the methodfilter package.
  Further, the `cors.MinimumAllowHeaders` method is removed, as it is no longer necessary
  to allow for commonly used headers, the enforcing clients/browsers have them integrated.

Release 0.2.1
=============

- dependency updates
- updated copyright year

Release 0.2.0
=============

- introduced golangci-lint and adapted code accordingly
  - removed underscore of all package names
  - renamed package util to helper
  - introduced error variables for common errors
- added race check to tests
- introduced REUSE compliance check
- added release workflow with provenance generation

Release 0.1.2
=============

- update crypto dependency

Release 0.1.1
=============

- fix add_headers to allow for duplicate headers

Release 0.1.0
=============

Initial release

- added access log middleware
- added header adding middleware
- added basic auth middleware
- added correlation ID middleware
- added CORS middleware
- added method filter middleware
- added rate limiting middleware