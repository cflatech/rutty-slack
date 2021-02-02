package slack

import (
	"log"
	"os"
	"strconv"
	"strings"

	// clientって外部との通信部だからビジネスロジックから依存するのは微妙な気がする

	"github.com/k1hiiragi/rutty-slack/client"
	"github.com/k1hiiragi/rutty-slack/domain/command"
	"github.com/k1hiiragi/rutty-slack/domain/rutty"
	"github.com/slack-go/slack"
)

// Run return nil
func Run() {
	// TODO: ログ出力
	api := slack.New(os.Getenv("SLACK_TOKEN"))
	rtm := api.NewRTM()
	var botID string
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			botID = rtm.GetInfo().User.ID
			log.Println(botID)
		case *slack.MessageEvent:
			channelID := ev.Channel
			log.Println(channelID)

			// 自分へのメンションか確認
			if !isMentionToBot(ev.Text, botID) {
				continue
			}

			// 送信用コマンドをParse
			command, parseErr := command.CreateCommand(ev.Text)
			log.Printf("respose = %v, error = %v", command, parseErr)
			if parseErr != nil {
				message := "入力をParseできませんでした！ごめんね！"
				sendMessage(message, channelID, rtm)
				continue
			}

			// RuttyへRequest送信
			responseData, requestErr := client.SendRuttyExecuteRequest(command)
			log.Printf("response = %v, error = %v", responseData, requestErr)
			if requestErr != nil {
				message := "Ruttyへのリクエストに失敗しました…"
				sendMessage(message, channelID, rtm)
				continue
			}

			// 結果送信
			message := makeExecResultMessage(responseData)
			log.Printf("sendMessage = \n%v", message)
			sendMessage(message, channelID, rtm)
		}
	}
}

func makeExecResultMessage(responseData rutty.ResponseData) string {
	var stdout = ""
	if len(responseData.Stdout) != 0 {
		stdout = "*stdout*\n```" + responseData.Stdout + "\n```\n"
	}

	var stderr = ""
	if len(responseData.Stderr) != 0 {
		stderr = "*stderr*\n```" + responseData.Stderr + "\n```\n"
	}

	var strReturn = ""
	if responseData.Rc != 0 {
		strReturn = "*return code*\n" + strconv.Itoa(responseData.Rc) + "\n"
	}

	return stdout + stderr + strReturn
}

func sendMessage(message string, channelID string, rtm *slack.RTM) {
	rtm.SendMessage(rtm.NewOutgoingMessage(message, channelID))
}

func isMentionToBot(message, botID string) bool {
	return strings.HasPrefix(message, "<@"+botID+">")
}
