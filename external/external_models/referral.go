package external_models

type ReferralCreateModel struct {
	AccountId    int    `json:"account_id"`
	ReferralCode string `json:"referral_code"`
}
