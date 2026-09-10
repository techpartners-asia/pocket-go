package pocket

import (
	"context"
	"encoding/json"
	"sync"
)

type pocket struct {
	merchant      string
	client_id     string
	client_secret string
	terminal_id   int64
	environment   string // "sandbox" or "production" | default is "production"
	oauthHost     string
	merchantHost  string
	realm         string

	// token is the credential installed by SetToken. The SDK reads it and
	// never populates it on its own; see Token.
	mu    sync.RWMutex
	token Token
}

// Pocket [Pocket SDK Interface]
//
// # Authentication
//
// This SDK does not manage tokens. Obtain one with [Pocket.Login], install it
// with [Pocket.SetToken], and every call below carries it. A call made with no
// token installed fails with [ErrNoToken]; a call whose token Pocket rejects
// fails with [ErrUnauthorized], which is the signal to log in again and retry.
type Pocket interface {
	// Login [Access Token авах] — one request, no caching.
	Login(ctx context.Context) (Token, error)

	// SetToken installs the token subsequent calls carry.
	SetToken(token Token)

	// Token returns the installed token.
	Token() Token

	CreateInvoice(input PocketCreateInvoiceInput) (PocketCreateInvoiceResponse, error)
	GetInvoiceByInvoiceID(invoiceID string) (PocketInvoiceDetailResponse, error)
	GetInvoiceByOrderNumber(orderNumber string) (PocketInvoiceDetailResponse, error)

	// Close clears the installed token.
	Close()
}

const (
	realm               = "invescore"
	oauthHostProd       = "sso.invescore.mn"
	merchantHostProd    = "service.invescore.mn/merchant"
	oauthHostSandbox    = "sso-staging.invescore.mn"
	merchantHostSandbox = "service-staging.invescore.mn/merchant"
)

// New performs no network I/O. The returned client has no token until one is
// installed with [Pocket.SetToken]; see [Pocket] on authentication.
func New(merchant, client_id, client_secret, environment string, terminal_id int64) Pocket {

	var oauthHost, merchantHost string

	switch environment {
	case "production":
		oauthHost = oauthHostProd
		merchantHost = merchantHostProd
	case "sandbox":
		oauthHost = oauthHostSandbox
		merchantHost = merchantHostSandbox
	default:
		oauthHost = oauthHostProd
		merchantHost = merchantHostProd
	}

	return &pocket{
		merchant:      merchant,
		client_id:     client_id,
		client_secret: client_secret,
		terminal_id:   terminal_id,
		environment:   environment,
		oauthHost:     oauthHost,
		realm:         realm,
		merchantHost:  merchantHost,
	}
}

func (p *pocket) CreateInvoice(input PocketCreateInvoiceInput) (PocketCreateInvoiceResponse, error) {

	body := PocketCreateInvoiceRequest{
		TerminalID:  p.terminal_id,
		Amount:      input.Amount,
		Info:        input.Info,
		OrderNumber: input.OrderNumber,
		InvoiceType: input.InvoiceType,
		Channel:     input.Channel,
	}

	res, err := p.httpRequest(body, PocketCreateInvoice, "")
	if err != nil {
		return PocketCreateInvoiceResponse{}, err
	}

	var response PocketCreateInvoiceResponse
	json.Unmarshal(res, &response)
	return response, nil
}

func (p *pocket) GetInvoiceByInvoiceID(invoiceID string) (PocketInvoiceDetailResponse, error) {
	request := PocketInvoiceDetailByInvoiceIDInput{
		TerminalID: p.terminal_id,
		InvoiceID:  invoiceID,
	}

	res, err := p.httpRequest(request, PocketGetInvoiceByInvoiceID, "")
	if err != nil {
		return PocketInvoiceDetailResponse{}, err
	}

	var response PocketInvoiceDetailResponse
	json.Unmarshal(res, &response)
	return response, nil
}

func (p *pocket) GetInvoiceByOrderNumber(orderNumber string) (PocketInvoiceDetailResponse, error) {

	request := PocketInvoiceDetailByOrderNumberInput{
		TerminalID:  p.terminal_id,
		OrderNumber: orderNumber,
	}

	res, err := p.httpRequest(request, PocketGetInvoiceByOrderNumber, "")
	if err != nil {
		return PocketInvoiceDetailResponse{}, err
	}

	var response PocketInvoiceDetailResponse
	json.Unmarshal(res, &response)
	return response, nil
}

// Close clears the installed token. It does not reach Pocket: there is
// nothing to revoke, and the token may still be in use elsewhere by whoever
// owns it.
func (s *pocket) Close() {
	s.SetToken(Token{})
}
