package external_models

type GetDisbursement struct {
	Status  string          `json:"status"`
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    []Disbursements `json:"data"`
}

type Disbursements struct {
	ID                    *int                      `json:"id"`
	DisbursementID        *int                      `json:"disbursement_id"`
	RecipientID           *int                      `json:"recipient_id"`
	PaymentID             *string                   `json:"payment_id"`
	BusinessID            *int                      `json:"business_id"`
	Amount                *string                   `json:"amount"`
	Narration             *string                   `json:"narration"`
	Currency              *string                   `json:"currency"`
	Reference             *string                   `json:"reference"`
	CallbackUrl           *string                   `json:"callback_url"`
	BeneficiaryName       *string                   `json:"beneficiary_name"`
	DestinationBranchCode *string                   `json:"destination_branch_code"`
	DebitCurrency         *string                   `json:"debit_currency"`
	Gateway               *string                   `json:"gateway"`
	Type                  *string                   `json:"type"`
	Status                *string                   `json:"status"`
	PaymentReleasedAt     *string                   `json:"payment_released_at"`
	DeletedAt             *string                   `json:"deleted_at"`
	CreatedAt             *string                   `json:"created_at"`
	UpdatedAt             *string                   `json:"updated_at"`
	Fee                   *int                      `json:"fee"`
	Tries                 *int                      `json:"tries"`
	TryAgainAt            *string                   `json:"try_again_at"`
	BankAccountNumber     *string                   `json:"bank_account_number"`
	BankName              *string                   `json:"bank_name"`
	UserDetails           *DisbursementsUserDetails `json:"user_details"`
}

type DisbursementsUserDetails struct {
	ID    *int    `json:"id"`
	Name  *string `json:"name"`
	Email *string `json:"email"`
}
