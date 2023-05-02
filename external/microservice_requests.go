package external

import (
	"fmt"

	"github.com/vesicash/auth-ms/internal/config"
)

type RequestObj struct {
	Path         string
	Method       string
	Headers      map[string]string
	SuccessCode  int
	RequestData  interface{}
	DecodeMethod decodemethod
}

type (
	decodemethod string
)

var (
	JsonDecodeMethod    decodemethod = "json"
	PhpSerializerMethod decodemethod = "phpserializer"
)

func FindMicroserviceRequest(name string, headers map[string]string, data interface{}) (RequestObj, error) {
	var (
		config = config.GetConfig()
	)
	switch name {
	case "welcome_notification":
		return RequestObj{
			Path:         fmt.Sprintf("%v/v2/send/send_welcome_mail", config.Microservices.Notification),
			Method:       "POST",
			Headers:      headers,
			SuccessCode:  200,
			RequestData:  data,
			DecodeMethod: JsonDecodeMethod,
		}, nil
	case "welcome_sms_notification":
		return RequestObj{
			Path:         fmt.Sprintf("%v/v2/send/send_welcome_sms", config.Microservices.Notification),
			Method:       "POST",
			Headers:      headers,
			SuccessCode:  200,
			RequestData:  data,
			DecodeMethod: JsonDecodeMethod,
		}, nil
	case "send_otp_notification":
		return RequestObj{
			Path:         fmt.Sprintf("%v/v2/send/send_otp", config.Microservices.Notification),
			Method:       "POST",
			Headers:      headers,
			SuccessCode:  200,
			RequestData:  data,
			DecodeMethod: JsonDecodeMethod,
		}, nil
	case "welcome_password_reset_notification":
		return RequestObj{
			Path:         fmt.Sprintf("%v/v2/send/send_welcome_password_mail", config.Microservices.Notification),
			Method:       "POST",
			Headers:      headers,
			SuccessCode:  200,
			RequestData:  data,
			DecodeMethod: JsonDecodeMethod,
		}, nil
	case "email_password_reset_notification":
		return RequestObj{
			Path:         fmt.Sprintf("%v/v2/send/send_reset_password_mail", config.Microservices.Notification),
			Method:       "POST",
			Headers:      headers,
			SuccessCode:  200,
			RequestData:  data,
			DecodeMethod: JsonDecodeMethod,
		}, nil
	case "phone_password_reset_notification":
		return RequestObj{
			Path:         fmt.Sprintf("%v/v2/send/send_reset_password_sms", config.Microservices.Notification),
			Method:       "POST",
			Headers:      headers,
			SuccessCode:  200,
			RequestData:  data,
			DecodeMethod: JsonDecodeMethod,
		}, nil
	case "email_password_reset_done_notification":
		return RequestObj{
			Path:         fmt.Sprintf("%v/v2/send/send_reset_password_done_mail", config.Microservices.Notification),
			Method:       "POST",
			Headers:      headers,
			SuccessCode:  200,
			RequestData:  data,
			DecodeMethod: JsonDecodeMethod,
		}, nil
	case "phone_password_reset_done_notification":
		return RequestObj{
			Path:         fmt.Sprintf("%v/v2/send/send_reset_password_done_sms", config.Microservices.Notification),
			Method:       "POST",
			Headers:      headers,
			SuccessCode:  200,
			RequestData:  data,
			DecodeMethod: JsonDecodeMethod,
		}, nil
	case "verification_email":
		return RequestObj{
			Path:         fmt.Sprintf("%v/v2/email", config.Microservices.Verification),
			Method:       "POST",
			Headers:      headers,
			SuccessCode:  200,
			RequestData:  data,
			DecodeMethod: JsonDecodeMethod,
		}, nil
	case "verification_sms":
		return RequestObj{
			Path:         fmt.Sprintf("%v/v2/phone", config.Microservices.Verification),
			Method:       "POST",
			Headers:      headers,
			SuccessCode:  200,
			RequestData:  data,
			DecodeMethod: JsonDecodeMethod,
		}, nil
	case "get_verifications":
		return RequestObj{
			Path:         fmt.Sprintf("%v/v2/fetch", config.Microservices.Verification),
			Method:       "GET",
			Headers:      headers,
			SuccessCode:  200,
			RequestData:  data,
			DecodeMethod: JsonDecodeMethod,
		}, nil
	case "create_referral":
		return RequestObj{
			Path:         fmt.Sprintf("%v/create", config.Microservices.Referral),
			Method:       "POST",
			Headers:      headers,
			SuccessCode:  200,
			RequestData:  data,
			DecodeMethod: JsonDecodeMethod,
		}, nil
	case "get_disbursements":
		return RequestObj{
			Path:         fmt.Sprintf("%v/disbursement/user", config.Microservices.Payment),
			Method:       "POST",
			Headers:      headers,
			SuccessCode:  200,
			RequestData:  data,
			DecodeMethod: JsonDecodeMethod,
		}, nil
	default:
		return RequestObj{}, fmt.Errorf("request not found")
	}
}
