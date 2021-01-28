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
	trimed := strings.TrimSpace(message)
	rep, _ := regexp.Compile("([[:space:]]|\u00a0)+")
	splited := rep.Split(trimed, 3)

	if len(splited) < 3 {
		return Command{"", ""}, errors.New("Parse Error")
	}

	blockStart := strings.Index(splited[2], "```") + 3
	blockEnd := strings.LastIndex(splited[2], "```")
	if blockStart < 0 || blockEnd < 0 {
		return Command{"", ""}, errors.New("Code not found")
	}

	return Command{splited[1], splited[2][blockStart:blockEnd]}, nil
}

// Language return string
func (t Command) Language() string {
	return t.language
}

// Code return string
func (t Command) Code() string {
	return t.code
}
