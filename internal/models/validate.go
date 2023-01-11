package models

type ValidateOnDBReq struct {
	Table string      `validate:"required" json:"table"`
	Type  string      `validate:"required" json:"type"`
	Query string      `validate:"required" json:"query"`
	Value interface{} `json:"value"`
}

type ValidateAuthorizationReq struct {
	Type               string `validate:"required" json:"type"`
	AuthorizationToken string `json:"authorization-token"`
	VApp               string `json:"v-app"`
	VPrivateKey        string `json:"v-private-key"`
	VPublicKey         string `json:"v-public-key"`
}
