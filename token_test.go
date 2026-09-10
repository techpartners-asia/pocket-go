package pocket

import (
	"errors"
	"testing"
	"time"
)

func newTestPocket() Pocket {
	return New("merchant", "client-id", "client-secret", "sandbox", 1)
}

// A call with no token must fail before any request is built, rather than
// going out with an empty bearer and coming back as a confusing 401.
func TestCallWithoutTokenIsErrNoToken(t *testing.T) {
	p := newTestPocket()

	if _, err := p.CreateInvoice(PocketCreateInvoiceInput{Amount: 100}); !errors.Is(err, ErrNoToken) {
		t.Fatalf("CreateInvoice: expected ErrNoToken, got %v", err)
	}
	if _, err := p.GetInvoiceByInvoiceID("inv-1"); !errors.Is(err, ErrNoToken) {
		t.Fatalf("GetInvoiceByInvoiceID: expected ErrNoToken, got %v", err)
	}
	if _, err := p.GetInvoiceByOrderNumber("order-1"); !errors.Is(err, ErrNoToken) {
		t.Fatalf("GetInvoiceByOrderNumber: expected ErrNoToken, got %v", err)
	}
}

func TestSetTokenAndClose(t *testing.T) {
	p := newTestPocket()

	p.SetToken(Token{AccessToken: "tok", ExpiresAt: time.Now().Add(time.Hour)})
	if got := p.Token().AccessToken; got != "tok" {
		t.Fatalf("expected the installed token, got %q", got)
	}

	p.Close()
	if !p.Token().IsZero() {
		t.Fatal("Close did not clear the installed token")
	}
	if _, err := p.CreateInvoice(PocketCreateInvoiceInput{Amount: 100}); !errors.Is(err, ErrNoToken) {
		t.Fatalf("expected ErrNoToken after Close, got %v", err)
	}
}

// Keycloak reports expires_in as a duration, so the expiry has to be anchored
// to when the response arrived: a bare number leaves the caller guessing.
func TestTokenFromAnchorsDurations(t *testing.T) {
	now := time.Date(2026, 8, 26, 12, 0, 0, 0, time.UTC)

	tok := tokenFrom(PocketGetTokenResponse{
		AccessToken:      "tok",
		RefreshToken:     "rtok",
		ExpiresIn:        300,
		RefreshExpiresIn: 1800,
		TokenType:        "Bearer",
	}, now)

	if want := now.Add(5 * time.Minute); !tok.ExpiresAt.Equal(want) {
		t.Fatalf("ExpiresAt = %v, want %v", tok.ExpiresAt, want)
	}
	if want := now.Add(30 * time.Minute); !tok.RefreshExpiresAt.Equal(want) {
		t.Fatalf("RefreshExpiresAt = %v, want %v", tok.RefreshExpiresAt, want)
	}

	// A missing expiry must stay the zero time rather than become "now".
	tok = tokenFrom(PocketGetTokenResponse{AccessToken: "tok"}, now)
	if !tok.ExpiresAt.IsZero() {
		t.Fatalf("expected a zero ExpiresAt for a missing expires_in, got %v", tok.ExpiresAt)
	}
}
