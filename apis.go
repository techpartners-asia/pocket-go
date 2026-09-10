package pocket

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/techpartners-asia/pocket-go/utils"
)

// Pocket
var (
	PocketGetToken = utils.API{
		Url:    "/auth/realms",
		Method: http.MethodPost,
	}
	PocketCreateInvoice = utils.API{
		Url:    "/v2/invoicing/generate-invoice",
		Method: http.MethodPost,
	}
	PocketGetInvoiceByOrderNumber = utils.API{
		Url:    "/v2/invoicing/invoices/order-number",
		Method: http.MethodPost,
	}
	PocketGetInvoiceByInvoiceID = utils.API{
		Url:    "/v2/invoicing/invoices/invoice-id",
		Method: http.MethodPost,
	}
)

func (p *pocket) httpRequest(body interface{}, api utils.API, urlExt string) (response []byte, err error) {
	token := p.Token()
	if token.IsZero() {
		// The SDK no longer authenticates behind the caller's back; an absent
		// token is a wiring mistake, not a gateway failure.
		return nil, ErrNoToken
	}

	var requestBody io.Reader = bytes.NewReader(nil)
	if body != nil {
		requestByte, marshalErr := json.Marshal(body)
		if marshalErr != nil {
			return nil, fmt.Errorf("pocket: encode request: %w", marshalErr)
		}
		requestBody = bytes.NewReader(requestByte)
	}

	baseUrl := "https://" + p.merchantHost + api.Url + urlExt

	req, err := http.NewRequest(api.Method, baseUrl, requestBody)
	if err != nil {
		return nil, fmt.Errorf("pocket: build request: %w", err)
	}
	req.Header.Add("Content-Type", utils.HttpContent)
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("pocket: request failed: %w", err)
	}
	defer res.Body.Close()

	response, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("pocket: read response: %w", err)
	}

	// A rejected token is its own error so the caller can tell "replace the
	// token and retry" apart from "Pocket refused this request".
	if res.StatusCode == http.StatusUnauthorized || res.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("%w (Status: %d): %s", ErrUnauthorized, res.StatusCode, string(response))
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(string(response))
	}

	return response, nil
}

// Login [Pocket SSO-гоос Access Token авах]
//
// Login performs exactly one request and caches nothing: the returned token is
// the caller's to hold, store and install with [Pocket.SetToken]. Concurrent
// callers each issue their own request, so deduplicating them is the caller's
// job too.
func (p *pocket) Login(ctx context.Context) (Token, error) {
	baseUrl := "https://" + p.oauthHost + "/auth/realms/" + p.realm + "/protocol/openid-connect/token"

	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("client_id", p.client_id)
	form.Add("client_secret", p.client_secret)

	req, err := http.NewRequestWithContext(ctx, PocketGetToken.Method, baseUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return Token{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Token{}, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer res.Body.Close()

	// Rejected credentials are not a provider outage: the caller can tell the
	// two apart and avoid counting a configuration mistake against Pocket.
	if res.StatusCode == http.StatusUnauthorized || res.StatusCode == http.StatusForbidden {
		return Token{}, fmt.Errorf("%w (Status: %d)", ErrUnauthorized, res.StatusCode)
	}
	if res.StatusCode != http.StatusOK {
		return Token{}, fmt.Errorf("auth request failed with status %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Token{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var authRes PocketGetTokenResponse
	if err = json.Unmarshal(body, &authRes); err != nil {
		return Token{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if authRes.AccessToken == "" {
		// Returning a tokenless response as success would send every later
		// call out with an empty bearer.
		return Token{}, errors.New("pocket: auth response contained no access token")
	}

	return tokenFrom(authRes, time.Now()), nil
}

// SetToken installs the token subsequent calls will carry. Passing the zero
// Token clears it, which makes the next call fail with [ErrNoToken] rather
// than reach Pocket unauthenticated.
func (p *pocket) SetToken(token Token) {
	p.mu.Lock()
	p.token = token
	p.mu.Unlock()
}

// Token returns the installed token.
func (p *pocket) Token() Token {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.token
}
