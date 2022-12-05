package external

import (
	"fmt"

	"github.com/vesicash/auth-ms/internal/config"
)

type RequestObj struct {
	Path        string
	Method      string
	Headers     map[string]string
	SuccessCode int
	RequestData interface{}
}

func FindMicroserviceRequest(name string, headers map[string]string, data interface{}) (RequestObj, error) {
	var (
		config = config.GetConfig()
	)
	switch name {
	case "welcome_notification":
		return RequestObj{
			Path:        fmt.Sprintf("%v/email/send/welcome", config.Microservices.Notification),
			Method:      "POST",
			Headers:     headers,
			SuccessCode: 200,
			RequestData: data,
		}, nil
	case "welcome_sms_notification":
		return RequestObj{
			Path:        fmt.Sprintf("%v/phone/send/welcome", config.Microservices.Notification),
			Method:      "POST",
			Headers:     headers,
			SuccessCode: 200,
			RequestData: data,
		}, nil
	case "send_otp_notification":
		return RequestObj{
			Path:        fmt.Sprintf("%v/send_otp", config.Microservices.Notification),
			Method:      "POST",
			Headers:     headers,
			SuccessCode: 200,
			RequestData: data,
		}, nil
	case "welcome_password_reset_notification":
		return RequestObj{
			Path:        fmt.Sprintf("%v/email/send/welcome_password", config.Microservices.Notification),
			Method:      "POST",
			Headers:     headers,
			SuccessCode: 200,
			RequestData: data,
		}, nil
	case "verification_email":
		return RequestObj{
			Path:        fmt.Sprintf("%v/email", config.Microservices.Verification),
			Method:      "POST",
			Headers:     headers,
			SuccessCode: 200,
			RequestData: data,
		}, nil
	case "verification_sms":
		return RequestObj{
			Path:        fmt.Sprintf("%v/phone", config.Microservices.Verification),
			Method:      "POST",
			Headers:     headers,
			SuccessCode: 200,
			RequestData: data,
		}, nil
	case "create_referral":
		return RequestObj{
			Path:        fmt.Sprintf("%v/create", config.Microservices.Referral),
			Method:      "POST",
			Headers:     headers,
			SuccessCode: 200,
			RequestData: data,
		}, nil
	default:
		return RequestObj{}, fmt.Errorf("request not found")
	}
}
