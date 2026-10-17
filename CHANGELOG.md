# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.1.2, 2026-10-17

- Documentation only. The README and `AGENTS.md` describe the daily regeneration workflow: what it checks, how it fixes failing checks, how it picks the next version and which paths a fix may not touch.

## v0.1.1, 2026-10-16

- Added the `go-revolut-business` agent skill. It installs with `claude plugin install go-revolut-business@psyb0t` in Claude Code, `codex plugin add go-revolut-business@psyb0t` in Codex and `openclaw skills install @psyb0t/go-revolut-business` from ClawHub.

## v0.1.0, 2026-10-16

- Added `revolut`, a REST client generated with oapi-codegen from Revolut's Business API OpenAPI document (spec version 1.0, OpenAPI 3.0.1, revolut-openapi commit 6f0da32 of 2026-09-29).
- `make generate` downloads the current spec from the revolut-openapi repository, patches the parts oapi-codegen cannot load, and regenerates the client. The spec is not committed.
