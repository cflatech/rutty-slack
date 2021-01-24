package main

import (
	"bytes"
	"encoding/json"
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
	// TODO: ログ出力
	// この辺の処理Serviceとかに切り出して良さそう
	api := slack.New(os.Getenv("SLACK_TOKEN"))
	rtm := api.NewRTM()
	var botID string
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			botID = rtm.GetInfo().User.ID
		case *slack.MessageEvent:
			channelID := ev.Channel

			// 自分へのメンションか確認
			if !isMentionToBot(ev.Text, botID) {
				continue
			}

			// 送信用コマンドをParse
			command, err := command.CreateCommand(ev.Text)
			if err != nil {
				message := "入力をParseできませんでした！ごめんね！"
				sendMessage(message, channelID, rtm)
				continue
			}

			// RuttyへRequest送信
			responseData, requestErr := sendRuttyRequest(command)
			if requestErr != nil {
				message := "Ruttyへのリクエストに失敗しました…"
				sendMessage(message, channelID, rtm)
				continue
			}

			// 結果送信
			message := makeExecResultMessage(responseData)
			sendMessage(message, channelID, rtm)
		}
	}
}

func makeExecResultMessage(responseData responseData) string {
	return "# stdout: \n" + responseData.Stdout + "\n" +
		"# stderr: \n" + responseData.Stderr + "\n" +
		"# return: \n" + strconv.Itoa(responseData.Rc)
}

func sendMessage(message string, channelID string, rtm *slack.RTM) {
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
