package revolut_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	revolut "github.com/psyb0t/go-revolut-business"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testToken      = "test-access-token"
	testAccountID  = "a1b2c3d4-0000-4000-8000-000000000001"
	testDepartment = "d1e2f3a4-0000-4000-8000-000000000002"
)

type capturedRequest struct {
	method string
	path   string
	auth   string
	body   []byte
}

// newClient starts an httptest server that records the one request it gets
// and answers with status and body, and returns a client pointed at it.
func newClient(
	t *testing.T,
	status int,
	body string,
) (*revolut.ClientWithResponses, *capturedRequest) {
	t.Helper()

	got := &capturedRequest{}

	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			payload, err := io.ReadAll(r.Body)
			assert.NoError(t, err)

			got.method = r.Method
			got.path = r.URL.Path
			got.auth = r.Header.Get("Authorization")
			got.body = payload

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			_, err = io.WriteString(w, body)
			assert.NoError(t, err)
		},
	))
	t.Cleanup(srv.Close)

	auth := revolut.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+testToken)

			return nil
		},
	)

	client, err := revolut.NewClientWithResponses(srv.URL, auth)
	require.NoError(t, err)

	return client, got
}

func TestGetAccounts(t *testing.T) {
	t.Parallel()

	client, got := newClient(t, http.StatusOK, `[{
		"id": "`+testAccountID+`",
		"name": "Main",
		"balance": 125.5,
		"currency": "EUR",
		"state": "active",
		"public": false,
		"account_type": "current",
		"created_at": "2026-01-02T03:04:05Z",
		"updated_at": "2026-01-02T03:04:05Z"
	}]`)

	res, err := client.GetAccountsWithResponse(t.Context())
	require.NoError(t, err)

	assert.Equal(t, http.MethodGet, got.method)
	assert.Equal(t, "/accounts", got.path)
	assert.Equal(t, "Bearer "+testToken, got.auth)

	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)
	require.Len(t, *res.JSON200, 1)

	account := (*res.JSON200)[0]
	assert.Equal(t, testAccountID, account.Id.String())
	assert.InDelta(t, 125.5, account.Balance, 0)
	assert.Equal(t, revolut.AccountCurrency("EUR"), account.Currency)
	require.NotNil(t, account.Name)
	assert.Equal(t, "Main", *account.Name)
}

func TestGetAccountByID(t *testing.T) {
	t.Parallel()

	client, got := newClient(t, http.StatusOK, `{
		"id": "`+testAccountID+`",
		"balance": 0,
		"currency": "GBP",
		"state": "active",
		"public": true,
		"account_type": "savings",
		"created_at": "2026-01-02T03:04:05Z",
		"updated_at": "2026-01-02T03:04:05Z"
	}`)

	res, err := client.GetAccountWithResponse(
		t.Context(),
		uuid.MustParse(testAccountID),
	)
	require.NoError(t, err)

	assert.Equal(t, http.MethodGet, got.method)
	assert.Equal(t, "/accounts/"+testAccountID, got.path)
	assert.Equal(t, "Bearer "+testToken, got.auth)

	require.NotNil(t, res.JSON200)
	assert.True(t, res.JSON200.Public)
	assert.Equal(
		t,
		revolut.AccountAccountType("savings"),
		res.JSON200.AccountType,
	)
}

func TestCreateDepartment(t *testing.T) {
	t.Parallel()

	client, got := newClient(
		t,
		http.StatusCreated,
		`{"id": "`+testDepartment+`"}`,
	)

	res, err := client.CreateDepartmentWithResponse(
		t.Context(),
		revolut.CreateDepartmentJSONRequestBody{Name: "Finance"},
	)
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, got.method)
	assert.Equal(t, "/departments", got.path)
	assert.Equal(t, "Bearer "+testToken, got.auth)

	var sent map[string]string

	require.NoError(t, json.Unmarshal(got.body, &sent))
	assert.Equal(t, map[string]string{"name": "Finance"}, sent)

	require.NotNil(t, res.JSON201)
	assert.Equal(t, testDepartment, res.JSON201.Id.String())
}

func TestErrorResponseDecodes(t *testing.T) {
	t.Parallel()

	client, _ := newClient(
		t,
		http.StatusUnauthorized,
		`{"code": 3001, "message": "invalid token"}`,
	)

	res, err := client.GetAccountsWithResponse(t.Context())
	require.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode())
	assert.Nil(t, res.JSON200)
	require.NotNil(t, res.JSONDefault)
	assert.Equal(t, "invalid token", res.JSONDefault.Message)
}
