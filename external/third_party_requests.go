package external

import (
	"fmt"
)

func FindThirdPartyRequest(name string, headers map[string]string, data interface{}) (RequestObj, error) {
	var (
	// config = config.GetConfig()
	)
	switch name {
	case "get_ip_info":
		return RequestObj{
			Path:         "http://www.geoplugin.net/php.gp",
			Method:       "Get",
			Headers:      headers,
			SuccessCode:  200,
			RequestData:  data,
			DecodeMethod: PhpSerializerMethod,
		}, nil
	default:
		return RequestObj{}, fmt.Errorf("request not found")
	}
}
