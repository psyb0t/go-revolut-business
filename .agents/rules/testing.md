# Testing

Refines `~/.claude/rule-details/go/writing-tests.md`.

- **TST1** Tests are black-box, in package `revolut_test`, and use only the exported API.
- **TST2** Tests never call Revolut. They point the client at an `httptest` server through `newClient` in `client_test.go`, which records the request and answers with a canned status and body.
- **TST3** Every test calls `t.Parallel()`. Assert the request method, path, auth header and body, then the decoded `JSONnnn` or `JSONDefault` field.
- **TST4** `example_test.go` holds the README usage example so it has to compile. Keep the two in sync.
- **TST5** `make test` runs `go test -race` and must pass.
