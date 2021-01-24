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
	lang, err := parseLanguage(message)
	if err != nil {
		return Command{"", ""}, err
	}

	code := parseCode(message)
	return Command{lang, code}, nil
}

func parseLanguage(message string) (string, error) {
	trimed := strings.TrimSpace(message)
	rep := regexp.MustCompile(`[:space:]+`)
	replaced := rep.ReplaceAllString(trimed, " ")
	splited := strings.Fields(replaced)
	if len(splited) < 2 {
		// 言語指定がない場合、空白を返す
		return "", errors.New("Language not found")
	}

	return splited[1], nil
}

func parseCode(message string) string {
	blockStart := strings.Index(message, "```") + 3
	blockEnd := strings.LastIndex(message, "```")
	return message[blockStart:blockEnd]
}

// Language return string
func (t Command) Language() string {
	return t.language
}

// Code return string
func (t Command) Code() string {
	return t.code
}
