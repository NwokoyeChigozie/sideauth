package auth

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/auth-ms/external/microservice/verification"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/auth-ms/utility"
)

func UpgradeUserTierService(db postgresql.Databases, tier int, accountID int) (int, error) {
	user := models.User{AccountID: uint(accountID)}
	code, err := user.GetUserByAccountID(db.Auth)
	if err != nil {
		return code, err
	}

	user.TierType = tier
	err = user.Update(db.Auth)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func GetUserRestrictionsService(logger *utility.Logger, db postgresql.Databases, accountID int) (map[string]interface{}, int, error) {
	var (
		tier         = 0
		empty_fields = []string{}
		tier_status  = "complete" // complete incomplete
		restrictions = gin.H{
			"tier":         tier,
			"empty_fields": empty_fields,
			"tier_status":  tier_status,
		}
	)
	user := models.User{AccountID: uint(accountID)}
	code, err := user.GetUserByAccountID(db.Auth)
	if err != nil {
		return restrictions, code, err
	}
	tier = user.TierType
	dataSlice, code, err := TierChecks(logger, tier, int(user.AccountID), db, user.LoginAccessToken)
	if err != nil {
		return restrictions, code, err
	}
	empty_fields = dataSlice
	if len(dataSlice) > 0 {
		tier_status = "incomplete"
	}

	return gin.H{
		"tier":         tier,
		"empty_fields": empty_fields,
		"tier_status":  tier_status,
	}, http.StatusOK, nil
}

func TierChecks(logger *utility.Logger, tierType, accountID int, db postgresql.Databases, accessToken string) ([]string, int, error) {
	response := []string{}
	user := models.User{AccountID: uint(accountID)}
	code, err := user.GetUserByAccountID(db.Auth)
	if err != nil {
		return response, code, err
	}

	var userMap map[string]interface{}
	inrec, err := json.Marshal(user)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}
	err = json.Unmarshal(inrec, &userMap)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	if tierType == 1 {
		fields := []string{"firstname", "lastname", "email_address", "username", "phone_number"}
		for _, v := range fields {
			if userMap[v] == "" {
				response = append(response, v)
			}
		}

	} else if tierType == 2 {
		fieldsMap := map[string]int{}
		fields := []string{"national_id", "bvn"}
		verifications, _ := verification.GetVerifications(logger, db.Auth, int(user.AccountID), accessToken)
		for _, v := range verifications {
			if v.IsVerified {
				val := fieldsMap[v.VerificationType] + 1
				fieldsMap[v.VerificationType] = val

			}
		}

		for _, f := range fields {
			if fieldsMap[f] < 1 {
				response = append(response, f)
			}
		}

	}
	return response, http.StatusOK, nil

}
