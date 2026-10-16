package revolut_test

import (
	"context"
	"net/http"
	"os"

	revolut "github.com/psyb0t/go-revolut-business"
)

// The README usage example, kept here so it has to compile. It talks to the
// live sandbox, so it has no Output block and go test does not run it.
//
//nolint:testableexamples // needs a real sandbox access token
func Example() {
	auth := revolut.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set(
				"Authorization",
				"Bearer "+os.Getenv("REVOLUT_ACCESS_TOKEN"),
			)

			return nil
		},
	)

	client, err := revolut.NewClientWithResponses(
		"https://sandbox-b2b.revolut.com/api/1.0",
		auth,
	)
	if err != nil {
		panic(err)
	}

	res, err := client.GetAccountsWithResponse(context.Background())
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
