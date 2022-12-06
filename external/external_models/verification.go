package external_models

type PhoneEmailVerificationModel struct {
	AccountId int `json:"account_id"`
	Token     int `json:"token"`
}

type GetVerifications struct {
	Status  string          `json:"status"`
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    []Verifications `json:"data"`
}

type Verifications struct {
	ID                 *int    `json:"id"`
	AccountID          *int    `json:"account_id"`
	VerificationCodeID *int    `json:"verification_code_id"`
	VerificationDocID  *int    `json:"verification_doc_id"`
	VerificationType   *string `json:"verification_type"`
	IsVerified         *bool   `json:"is_verified"`
	VerifiedAt         *string `json:"verified_at"`
	DeletedAt          *string `json:"deleted_at"`
	CreatedAt          *string `json:"created_at"`
	UpdatedAt          *string `json:"updated_at"`
	Tries              *int    `json:"tries"`
}
