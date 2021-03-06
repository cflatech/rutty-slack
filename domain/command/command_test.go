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
			arg:  "<@U01L10XDBQR>\u00a0ruby\nputs \"hello\"",
			want: Command{"ruby", "puts \"hello\""},
		},
		{
			name: "先頭に空白が入っている場合",
			arg:  " <@U01L10XDBQR>\u00a0ruby\nputs \"hello\"",
			want: Command{"ruby", "puts \"hello\""},
		},
		{
			name: "コード内に```がある場合",
			arg:  "<@U01L10XDBQR>\u00a0ruby\nputs 1 #```\nputs 2",
			want: Command{"ruby", "puts 1 #```\nputs 2"},
		},
		{
			name: "ブロックまでに改行が入る場合",
			arg:  "<@U01L10XDBQR> ruby\n\n\n\nputs \"ruby\"",
			want: Command{"ruby", "puts \"ruby\""},
		},
		{
			name: "言語指定までに複数の空白がある場合",
			arg:  "<@U01L10XDBQR>\u00a0   ruby\nputs \"ruby\"",
			want: Command{"ruby", "puts \"ruby\""},
		},
		{
			name: "言語がPHPの場合",
			arg:  "<@U01L10XDBQR>\u00a0php\n&lt;?php\necho 'hoge';",
			want: Command{"php", "<?php\necho 'hoge';"},
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
			want: "Parse Error",
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
