package thirdparty

import (
	"github.com/vesicash/auth-ms/external"
	"github.com/vesicash/auth-ms/utility"
)

func GetIpInfo(ip4IpAddress string) (map[interface{}]interface{}, error) {
	logger := utility.NewLogger()
	var (
		outBoundResponse map[interface{}]interface{}
		headers          = map[string]string{
			"Content-Type": "application/json",
		}
	)

	logger.Info("get ip info", nil)

	err := external.SendRequest(logger, "third_party", "get_ip_info", headers, nil, &outBoundResponse, "?ip="+ip4IpAddress)
	if err != nil {
		logger.Info("get ip info", outBoundResponse, err)
		return outBoundResponse, err
	}
	logger.Info("get ip info", outBoundResponse)
	logger.Info("get ip info2", outBoundResponse["geoplugin_regionName"], outBoundResponse["geoplugin_countryName"])

	return outBoundResponse, nil
}
