package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/k1hiiragi/rutty-slack/domain/command"
	"github.com/k1hiiragi/rutty-slack/domain/rutty"
)

// SendRuttyExecuteRequest return (responseData, error)
func SendRuttyExecuteRequest(command command.Command) (rutty.ResponseData, error) {
	requestJSON := makeRequestJSON(command.Code())

	apiURL := os.Getenv("RUTTY_API_URL")
	resp, err := http.Post(apiURL+"/executors/"+command.Language(), "application/json", bytes.NewBuffer(requestJSON))
	if err != nil {
		return rutty.ResponseData{}, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var execResult rutty.ResponseData
	marshalErr := json.Unmarshal(body, &execResult)
	if marshalErr != nil {
		return rutty.ResponseData{}, marshalErr
	}

	return execResult, nil

}

func makeRequestJSON(code string) []byte {
	requestData := rutty.RequestData{
		Code: code,
	}

	requestJSON, _ := json.Marshal(requestData)

	return requestJSON
}
