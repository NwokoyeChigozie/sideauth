package auth

import (
	"net/http"

	"github.com/vesicash/auth-ms/external/microservice/payment"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
)

func GetDisbursementsService(db postgresql.Databases, accountID int) ([]map[string]interface{}, int, error) {
	var (
		data = []map[string]interface{}{}
	)
	disbursements, _ := payment.GetDisbursement(db.Auth, accountID)

	for _, v := range disbursements {
		sData := map[string]interface{}{
			"disbursement_id":  v.DisbursementID,
			"recipient_id":     v.RecipientID,
			"business_id":      v.BusinessID,
			"amount":           v.Amount,
			"narration":        v.Narration,
			"currency":         v.Currency,
			"reference":        v.Reference,
			"beneficiary_name": v.BeneficiaryName,
			"type":             v.Type,
			"status":           v.Status,
		}
		data = append(data, sData)
	}

	return data, http.StatusOK, nil
}
