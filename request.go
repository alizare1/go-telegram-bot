package telegrambot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const formatedAPIUrl = "https://api.telegram.org/bot%v/%v"

type TelegramResponse struct {
	Ok          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	Description string          `json:"description"`
	ErrorCode   int             `json:"error_code"`
}

func getMethodUrl(method string, token string) (reqUrl string) {
	reqUrl = fmt.Sprintf(formatedAPIUrl, token, method)
	return reqUrl
}

func sendGetRequest(reqUrl string, params map[string]string) (tgResp TelegramResponse, err error) {
	parsedUrl, _ := url.Parse(reqUrl)
	if params != nil {
		parsedParams := url.Values{}
		for k, v := range params {
			parsedParams.Add(k, v)
		}
		parsedUrl.RawQuery = parsedParams.Encode()
	}

	resp, err := http.Get(parsedUrl.String())
	if err != nil {
		return tgResp, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&tgResp)

	return tgResp, err
}

func sendPostRequest(reqUrl string, body map[string]interface{}) (tgResp TelegramResponse, err error) {
	encodedBody, err := json.Marshal(body)
	if err != nil {
		return tgResp, err
	}

	resp, err := http.Post(reqUrl, "application/json", bytes.NewBuffer(encodedBody))
	if err != nil {
		return tgResp, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&tgResp)
	return tgResp, err
}
