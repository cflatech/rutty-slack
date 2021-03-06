package command

import (
	"errors"
	"regexp"
	"strings"
)

// Command Rutty用の言語とソースコード
type Command struct {
	language, code string
}

// CreateCommand return Command
func CreateCommand(message string) (Command, error) {
	return parse(message)
}

func parse(message string) (Command, error) {
	// slackの仕様, PHPの実行時などに問題が発生する
	// https://api.slack.com/reference/surfaces/formatting
	message = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(message, "&amp;", "&"), "&lt;", "<"), "&gt;", ">")
	trimed := strings.TrimSpace(message)
	rep, _ := regexp.Compile("([[:space:]]|\u00a0)+")
	splited := rep.Split(trimed, 3)

	if len(splited) < 3 {
		return Command{"", ""}, errors.New("Parse Error")
	}

	return Command{splited[1], splited[2]}, nil
}

// Language return string
func (t Command) Language() string {
	return t.language
}

// Code return string
func (t Command) Code() string {
	return t.code
}
