package pocket

type (
	PocketGetTokenResponse struct {
		AccessToken      string `json:"access_token"`
		RefreshToken     string `json:"refresh_token"`
		ExpiresIn        int    `json:"expires_in"`
		RefreshExpiresIn int    `json:"refresh_expires_in"`
		TokenType        string `json:"token_type"`
		Scope            string `json:"scope"`
		SessionState     string `json:"session_state"`
		NotBeforePolicy  int    `json:"not-before-policy"`
	}

	PocketListBranchResponse struct {
		Content          []PocketBranch `json:"content"`
		Pageable         Pageable       `json:"pageable"`
		TotalElements    int            `json:"totalElements"`
		Last             bool           `json:"last"`
		TotalPages       int            `json:"totalPages"`
		Number           int            `json:"number"`
		Sort             Sort           `json:"sort"`
		Size             int            `json:"size"`
		First            bool           `json:"first"`
		NumberOfElements int            `json:"numberOfElements"`
		Empty            bool           `json:"empty"`
	}

	PocketBranch struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		Address       string `json:"address"`
		AccountNumber string `json:"accountNumber"`
		MerchantID    string `json:"merchantId"`
	}

	PocketCreateInvoiceInput struct {
		Amount      float64 `json:"amount"`      // example: 10000.0
		Info        string  `json:"info"`        // example: POCKET-ZERO-R100010002001
		OrderNumber string  `json:"orderNumber"` // example: R100010002021
		InvoiceType string  `json:"invoiceType"` // зээлээр бол "ZERO", хэвэвчтээс бол "PURCHASE"
		Channel     string  `json:"channel"`     // "ecommerce" or "pos" default is "merchant"
	}

	PocketCreateInvoiceRequest struct {
		TerminalID  int64   `json:"terminalId"`  // example: 74686381183671
		Amount      float64 `json:"amount"`      // example: 10000.0
		Info        string  `json:"info"`        // example: POCKET-ZERO-R100010002001
		OrderNumber string  `json:"orderNumber"` // example: R100010002021
		InvoiceType string  `json:"invoiceType"` // зээлээр бол "ZERO", хэвэвчтээс бол "PURCHASE"
		Channel     string  `json:"channel"`     // "ecommerce" or "pos" default is "merchant"
	}

	PocketCreateInvoiceResponse struct {
		ID          uint   `json:"id"`
		Qr          string `json:"qr"`
		OrderNumber string `json:"orderNumber"`
		DeepLink    string `json:"deeplink"`
	}

	PocketInvoiceDetailByOrderNumberInput struct {
		TerminalID  int64  `json:"terminalId"`
		OrderNumber string `json:"orderNumber"`
	}

	PocketInvoiceDetailByInvoiceIDInput struct {
		TerminalID int64  `json:"terminalId"`
		InvoiceID  string `json:"invoiceId"`
	}

	PocketInvoiceDetailResponse struct {
		State        string  `json:"state"`
		Description  string  `json:"description"`
		SenderName   string  `json:"senderName"`
		ReceiverName string  `json:"receiverName"`
		Amount       float64 `json:"amount"`
		Info         string  `json:"info"`
		HoldID       string  `json:"holdId"`
		ID           uint    `json:"id"`
		CreatedAt    string  `json:"createdAt"`
		AliasName    string  `json:"aliasName"`
		TerminalID   int64   `json:"terminalId"`
		BranchName   string  `json:"branchName"`
		BranchID     int     `json:"branchId"`
		OrderNumber  string  `json:"orderNumber"`
		InvoiceType  string  `json:"invoiceType"`
	}

	PocketWebhookResponse struct {
		ID           uint    `json:"id"`
		Amount       float64 `json:"amount"`
		Info         string  `json:"info"`
		InvoiceID    string  `json:"invoiceId"`
		InvoiceState string  `json:"invoiceState"`
		HeldID       string  `json:"heldId"`
		PhoneNumber  string  `json:"phoneNumber"`
		OrderNumber  string  `json:"orderNumber"`
	}

	Pageable struct {
		Sort       Sort `json:"sort"`
		Offset     int  `json:"offset"`
		PageSize   int  `json:"pageSize"`
		PageNumber int  `json:"pageNumber"`
		Paged      bool `json:"paged"`
	}

	Sort struct {
		Sorted   bool `json:"sorted"`
		Unsorted bool `json:"unsorted"`
		Empty    bool `json:"empty"`
	}
)
