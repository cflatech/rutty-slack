package slack

import (
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
			responseData, requestErr := client.SendRuttyRequest(command)
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

func makeExecResultMessage(responseData rutty.ResponseData) string {
	var stdout = ""
	if len(responseData.Stdout) != 0 {
		stdout = "```" + responseData.Stdout + "\n```"
	}

	var stderr = ""
	if len(responseData.Stderr) != 0 {
		stderr = "```" + responseData.Stderr + "\n```"
	}

	return "# *stdout*: \n" + stdout + "\n" +
		"# *stderr*: \n" + stderr + "\n" +
		"# *return*: \n" + strconv.Itoa(responseData.Rc)
}

func sendMessage(message string, channelID string, rtm *slack.RTM) {
	rtm.SendMessage(rtm.NewOutgoingMessage(message, channelID))
}

func isMentionToBot(message, botID string) bool {
	return strings.HasPrefix(message, "<@"+botID+">")
}
