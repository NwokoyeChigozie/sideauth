package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/vesicash/auth-ms/utility"
)

func SendRequest(logger *utility.Logger, reqType, name string, headers map[string]string, data interface{}, response interface{}) error {
	var (
		reqObject = RequestObj{}
		err       error
	)
	if reqType == "service" {
		reqObject, err = FindMicroserviceRequest(name, headers, data)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("not implemented")
	}

	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(data)
	if err != nil {
		logger.Info("encoding error", name, err.Error())
	}

	logger.Info(name, reqObject.Path, data, buf)

	client := &http.Client{}
	req, err := http.NewRequest(reqObject.Method, reqObject.Path, buf)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = json.NewDecoder(res.Body).Decode(response)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	logger.Info("response body", name, reqObject.Path, body)

	defer res.Body.Close()

	if res.StatusCode == reqObject.SuccessCode {
		return nil
	}

	if res.StatusCode < 200 && res.StatusCode > 299 {
		return fmt.Errorf("Error " + strconv.Itoa(res.StatusCode))
	}

	return nil
}
