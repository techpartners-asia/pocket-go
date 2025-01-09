package pocket

import (
	"bytes"
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
	authObj, authErr := p.authPocket()
	if authErr != nil {
		err = authErr
		return
	}
	p.loginObject = &authObj

	var requestByte []byte
	var requestBody *bytes.Reader
	if body == nil {
		requestBody = bytes.NewReader(nil)
	} else {
		requestByte, _ = json.Marshal(body)
		requestBody = bytes.NewReader(requestByte)
	}

	baseUrl := "https://" + p.merchantHost + api.Url + urlExt

	req, _ := http.NewRequest(api.Method, baseUrl, requestBody)
	req.Header.Add("Content-Type", utils.HttpContent)
	req.Header.Add("Authorization", "Bearer "+p.loginObject.AccessToken)

	res, err := http.DefaultClient.Do(req)

	response, _ = io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		return nil, errors.New(string(response))
	}

	return
}

// AuthPocket [Login to pocket]
func (p *pocket) authPocket() (authRes PocketGetTokenResponse, err error) {
	if p.loginObject != nil && p.expire_in != nil {
		now := time.Now().Unix()
		if now < *p.expire_in {
			return *p.loginObject, nil
		}
	}

	baseUrl := "https://" + p.oauthHost + "/auth/realms/" + p.realm + "/protocol/openid-connect/token"

	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("client_id", p.client_id)
	form.Add("client_secret", p.client_secret)

	req, err := http.NewRequest(PocketGetToken.Method, baseUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return authRes, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return authRes, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return authRes, fmt.Errorf("auth request failed with status %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return authRes, fmt.Errorf("failed to read response body: %w", err)
	}

	if err = json.Unmarshal(body, &authRes); err != nil {
		return authRes, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	expireD := time.Now().Unix() + int64(authRes.ExpiresIn)
	p.expire_in = &expireD
	p.loginObject = &authRes

	return authRes, nil
}
