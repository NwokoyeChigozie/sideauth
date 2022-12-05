package external_models

type AccountIDModel struct {
	AccountId int `json:"account_id"`
}

type SendOtpModel struct {
	AccountId int `json:"account_id"`
	OtpToken  int `json:"otp_token"`
}
