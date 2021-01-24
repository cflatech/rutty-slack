package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/k1hiiragi/rutty-slack/domain/command"
)

// requestData Ruttyに投げるRequestData
type requestData struct {
	Code string `json:"code"`
}

// ResponseData Ruttyから返ってきたレスポンス
type ResponseData struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Rc     int    `json:"rc"`
}

// SendRuttyRequest return (responseData, error)
func SendRuttyRequest(command command.Command) (ResponseData, error) {
	requestJSON := makeRequestJSON(command.Code())

	// Todo: 環境変数に変える
	apiURL := os.Getenv("RUTTY_API_URL")
	resp, err := http.Post(apiURL+command.Language(), "application/json", bytes.NewBuffer(requestJSON))
	if err != nil {
		return ResponseData{}, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var execResult ResponseData
	marshalErr := json.Unmarshal(body, &execResult)
	if marshalErr != nil {
		return ResponseData{}, marshalErr
	}

	return execResult, nil

}

func makeRequestJSON(code string) []byte {
	requestData := requestData{
		Code: code,
	}

	requestJSON, _ := json.Marshal(requestData)

	return requestJSON
}
