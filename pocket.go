package pocket

import "encoding/json"

type pocket struct {
	merchant      string
	client_id     string
	client_secret string
	terminal_id   int64
	environment   string // "sandbox" or "production" | default is "production"
	loginObject   *PocketGetTokenResponse
	expire_in     *int64
	oauthHost     string
	merchantHost  string
	realm         string
}

type Pocket interface {
	CreateInvoice(input PocketCreateInvoiceInput) (PocketCreateInvoiceResponse, error)
	GetInvoiceByInvoiceID(invoiceID string) (PocketInvoiceDetailResponse, error)
	GetInvoiceByOrderNumber(orderNumber string) (PocketInvoiceDetailResponse, error)
	Close()
}

const (
	realm               = "invescore"
	oauthHostProd       = "sso.invescore.mn"
	merchantHostProd    = "service.invescore.mn/merchant"
	oauthHostSandbox    = "sso-staging.invescore.mn"
	merchantHostSandbox = "service-staging.invescore.mn/merchant"
)

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
		loginObject:   nil,
		expire_in:     nil,
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

func (s *pocket) Close() {
	s.loginObject = nil
	s.expire_in = nil
}
