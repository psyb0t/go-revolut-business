# go-revolut-business

[![Go Reference](https://pkg.go.dev/badge/github.com/psyb0t/go-revolut-business.svg)](https://pkg.go.dev/github.com/psyb0t/go-revolut-business)
[![CI](https://github.com/psyb0t/go-revolut-business/actions/workflows/pipeline.yml/badge.svg?branch=main)](https://github.com/psyb0t/go-revolut-business/actions/workflows/pipeline.yml)
[![version](https://raw.githubusercontent.com/psyb0t/go-revolut-business/badges/version.svg)](https://github.com/psyb0t/go-revolut-business/tags)
[![license](https://raw.githubusercontent.com/psyb0t/go-revolut-business/badges/license.svg)](LICENSE)
[![imported by](https://raw.githubusercontent.com/psyb0t/go-revolut-business/badges/importers.svg)](https://github.com/psyb0t/go-revolut-business/blob/badges/importers.md)

Unofficial Go client for the [Revolut Business API](https://developer.revolut.com/docs/api/business), generated with [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) from the [OpenAPI document](https://github.com/revolut-engineering/revolut-openapi) Revolut publishes.

| Package | API | Generator |
|---|---|---|
| `revolut` | REST, `https://b2b.revolut.com/api/1.0` | [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) |

`client.gen.go` is generated. Don't edit it by hand, regenerate it.

## Install

```sh
go get github.com/psyb0t/go-revolut-business
```

## Usage

Every endpoint in the spec has a typed method on `ClientWithResponses`.

```go
package main

import (
	"context"
	"net/http"
	"os"

	revolut "github.com/psyb0t/go-revolut-business"
)

func main() {
	auth := revolut.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+os.Getenv("REVOLUT_ACCESS_TOKEN"))

		return nil
	})

	client, err := revolut.NewClientWithResponses("https://sandbox-b2b.revolut.com/api/1.0", auth)
	if err != nil {
		panic(err)
	}

	res, err := client.GetAccountsWithResponse(context.Background())
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
```

The same code lives in `example_test.go`, so it has to compile.

The full API is in the [Go reference](https://pkg.go.dev/github.com/psyb0t/go-revolut-business); from inside your module, `go doc github.com/psyb0t/go-revolut-business` prints the same.

## Auth

Every request carries `Authorization: Bearer <access_token>`. Revolut issues the access token through OAuth2: you sign a JWT client assertion with the private key whose certificate you uploaded in the Business API settings, and exchange it together with an authorisation code or refresh token at `/auth/token`. Access tokens expire after 40 minutes. See [Revolut's Business API guides](https://developer.revolut.com/docs/guides/manage-accounts/introduction) for the details.

This module does not implement that flow. Get a token however you like and set it with a `RequestEditorFn`, as above.

| Environment | Base URL |
|---|---|
| Production | `https://b2b.revolut.com/api/1.0` |
| Sandbox | `https://sandbox-b2b.revolut.com/api/1.0` |

Sandbox uses test data and needs its own sandbox account and certificate.

## Agent integrations

The [skill](.agents/skills/go-revolut-business) teaches a coding agent how to use go-revolut-business when it writes Go. It works in any agent that reads `.agents/skills/`, and installs natively in the clients below.

### Claude Code

```bash
claude plugin marketplace add psyb0t/agents
claude plugin install go-revolut-business@psyb0t
```

### Codex

```bash
codex plugin marketplace add psyb0t/agents
codex plugin add go-revolut-business@psyb0t
```

Codex also picks the skill up with no install in any checkout of this repo.

### OpenClaw

The skill is published to ClawHub on every release:

```bash
openclaw skills install @psyb0t/go-revolut-business
```

## Regenerating

Revolut's spec is not committed. `make generate` downloads the current OpenAPI document into the gitignored `api/` directory and regenerates the client:

```sh
make generate
```

`make generate` is just `go generate ./...`, and `gen.go` declares the steps:

- `scripts/fetch-spec.sh` downloads `business.yaml` from the revolut-openapi repository.
- `scripts/patch-spec.sh` writes `api/business.patched.yaml`, which oapi-codegen can load. It drops numeric keywords that hold prose instead of numbers (`maximum: P7D`), rewrites camelCase path placeholders to the snake_case parameter names the operations declare, and gives inline property schemas that collide with a component schema of the same Go name an `x-go-type-name` ending in `Value` (`ExchangeReasonCodeValue`, `TransferReasonCodeValue`).
- oapi-codegen runs with `config.yaml`.

`.github/workflows/regen.yml` does the same every day in CI. When the new spec changes the client, it runs `make lint`, `make test` and `make sec`. If a check fails, Claude edits the code and the checks run again, up to three times, and a review step reads the result. A passing change is committed with a CHANGELOG entry and tagged as the next release: a patch version for fixes and additions, a minor version for breaking changes.

## Development

Every command runs in the `Dockerfile.dev` image. The host needs only Docker and Make.

```sh
make help      # list targets
make dep       # go mod tidy + vendor
make lint      # shellcheck, shfmt, go fix, golangci-lint
make test      # go test -race
make sec       # gitleaks, govulncheck, semgrep
make generate  # download the spec and regenerate
```

## License

[MIT](LICENSE). The generated code comes from Revolut's published API specification. The specification itself is not redistributed here. This project is not affiliated with Revolut.
