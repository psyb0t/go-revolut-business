# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.1.0, 2026-10-16

- Added `revolut`, a REST client generated with oapi-codegen from Revolut's Business API OpenAPI document (spec version 1.0, OpenAPI 3.0.1, revolut-openapi commit 6f0da32 of 2026-09-29).
- `make generate` downloads the current spec from the revolut-openapi repository, patches the parts oapi-codegen cannot load, and regenerates the client. The spec is not committed.
