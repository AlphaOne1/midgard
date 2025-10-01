<!-- SPDX-FileCopyrightText: 2025 The midgard contributors.
     SPDX-License-Identifier: MPL-2.0
-->
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