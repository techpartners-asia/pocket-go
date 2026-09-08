package pocket

import (
	"errors"
	"time"
)

// Token is a Pocket access token together with everything the caller needs to
// decide when to replace it.
//
// The SDK does not cache tokens: [Pocket.Login] performs exactly one HTTP call
// and returns what the SSO said, and [Pocket.SetToken] installs the token the
// next request will carry. Whoever owns the client owns the token's lifetime —
// expiry tracking, renewal scheduling and deduplication of concurrent logins
// are all the caller's, because only the caller knows whether the token is
// shared beyond this process.
type Token struct {
	// AccessToken is the bearer sent as Authorization on every API call.
	AccessToken string

	// RefreshToken is returned by the SSO. Pocket exposes no refresh call
	// here, so it is carried for completeness only: renew with [Pocket.Login].
	RefreshToken string

	// TokenType is the SSO's token_type, normally "Bearer".
	TokenType string

	// ExpiresAt is when AccessToken stops being accepted. It is the ZERO time
	// when the SSO reported no expires_in at all, and a caller must then treat
	// the token as good for this one use only.
	ExpiresAt time.Time

	// RefreshExpiresAt is when RefreshToken stops being accepted, under the
	// same zero-time caveat as ExpiresAt.
	RefreshExpiresAt time.Time

	// Scope is returned verbatim for diagnostics.
	Scope string
}

// IsZero reports whether the token carries no credential at all.
func (t Token) IsZero() bool { return t.AccessToken == "" }

var (
	// ErrNoToken is returned when an API call is attempted before a token has
	// been installed with [Pocket.SetToken]. The SDK never authenticates on
	// its own, so this is a wiring mistake in the caller, not a gateway
	// failure.
	ErrNoToken = errors.New("pocket: no access token installed: call Login then SetToken")

	// ErrUnauthorized is returned when Pocket rejects the installed token with
	// 401 or 403. The caller should discard the token, log in again and retry —
	// the rejected request was never processed, so retrying cannot
	// double-create an invoice.
	ErrUnauthorized = errors.New("pocket: access token rejected")
)

// tokenFrom converts an SSO token response into a Token.
//
// Keycloak reports expires_in as a duration in seconds, so the expiry is
// anchored to when we received it — which is why this conversion happens at
// the point of the response and not later.
func tokenFrom(res PocketGetTokenResponse, now time.Time) Token {
	return Token{
		AccessToken:      res.AccessToken,
		RefreshToken:     res.RefreshToken,
		TokenType:        res.TokenType,
		ExpiresAt:        anchoredExpiry(now, res.ExpiresIn),
		RefreshExpiresAt: anchoredExpiry(now, res.RefreshExpiresIn),
		Scope:            res.Scope,
	}
}

func anchoredExpiry(now time.Time, seconds int) time.Time {
	if seconds <= 0 {
		return time.Time{}
	}
	return now.Add(time.Duration(seconds) * time.Second)
}
