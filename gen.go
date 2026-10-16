// Package revolut is a Go client for the Revolut Business REST API,
// generated with oapi-codegen from the OpenAPI document Revolut publishes.
package revolut

//go:generate bash scripts/fetch-spec.sh
//go:generate bash scripts/patch-spec.sh
//go:generate go tool oapi-codegen --config=config.yaml api/business.patched.yaml
