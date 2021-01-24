package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/k1hiiragi/rutty-slack/domain/command"

	"github.com/slack-go/slack"
)

type requestData struct {
	Code string `json:"code"`
}

type responseData struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Rc     int    `json:"rc"`
}

func main() {
	api := slack.New(os.Getenv("SLACK_TOKEN"))
	rtm := api.NewRTM()
	var botID, channelID string
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			botID = rtm.GetInfo().User.ID
		case *slack.MessageEvent:
			channelID = ev.Channel
			fmt.Println(channelID)

			// 自分へのメンションか確認
			if !isMentionToBot(ev.Text, botID) {
				continue
			}

			// 送信用コマンドをParse
			command, err := command.CreateCommand(ev.Text)
			if err != nil {
				continue
			}

			// RuttyへRequest送信
			responseData, requestErr := sendRuttyRequest(command)
			if requestErr != nil {
				continue
			}

			// 結果送信
			sendMessage(responseData, channelID, rtm)
		}
	}
}

func sendMessage(responseData responseData, channelID string, rtm *slack.RTM) {
	message :=
		"# stdout: \n" + responseData.Stdout + "\n" +
			"# stderr: \n" + responseData.Stderr + "\n" +
			"# return: \n" + strconv.Itoa(responseData.Rc)
	rtm.SendMessage(rtm.NewOutgoingMessage(message, channelID))
}

func sendRuttyRequest(command command.Command) (responseData, error) {
	requestJSON := makeRequestJSON(command.Code())

	// Todo: 環境変数に変える
	resp, err := http.Post("http://localhost:3000/executors/"+command.Language(), "application/json", bytes.NewBuffer(requestJSON))
	if err != nil {
		return responseData{}, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var execResult responseData
	marshalErr := json.Unmarshal(body, &execResult)
	if marshalErr != nil {
		return responseData{}, marshalErr
	}

	return execResult, nil

}

func isMentionToBot(message, botID string) bool {
	return strings.HasPrefix(message, "<@"+botID+">")
}

func makeRequestJSON(code string) []byte {
	requestData := requestData{
		Code: code,
	}

	requestJSON, _ := json.Marshal(requestData)

	return requestJSON
}
