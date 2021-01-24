package command

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want Command
	}{
		{
			name: "正常系",
			arg:  "<@U01L10XDBQR>\u00a0ruby\n```puts \"hello\"```",
			want: Command{"ruby", "puts \"hello\""},
		},
		{
			name: "先頭に空白が入っている場合",
			arg:  " <@U01L10XDBQR>\u00a0ruby\n```puts \"hello\"```",
			want: Command{"ruby", "puts \"hello\""},
		},
		{
			name: "コード内に```がある場合",
			arg:  "<@U01L10XDBQR>\u00a0ruby\n```puts 1 #```\nputs 2```",
			want: Command{"ruby", "puts 1 #```\nputs 2"},
		},
		{
			name: "ブロックまでに改行が入る場合",
			arg:  "<@U01L10XDBQR> ruby\n\n\n\n```puts \"ruby\"```",
			want: Command{"ruby", "puts \"ruby\""},
		},
		{
			name: "言語指定までに複数の空白がある場合",
			arg:  "<@U01L10XDBQR>\u00a0   ruby\n```puts \"ruby\"```",
			want: Command{"ruby", "puts \"ruby\""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command, _ := parse(tt.arg)
			if command != tt.want {
				t.Fatalf("expect:%v != result:%v}", tt.want, command)
			}
		})
	}

}
func TestParseError(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "言語が指定されていない場合エラーとなる",
			arg:  "<@U01L10XDBQR>",
			want: "Language not found",
		},
		{
			name: "```が存在しない場合",
			arg:  "<@U01L10XDBQR>\u00a0   ruby\nputs \"ruby\"",
			want: "Code not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parse(tt.arg)
			if errorStr := err.Error(); errorStr != tt.want {
				t.Fatalf("expect:%v != result:%v}", tt.want, errorStr)
			}
		})
	}

}
