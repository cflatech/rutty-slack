package slack

import "testing"

func TestIsMentionToBot(t *testing.T) {
	type args struct {
		message,
		botID string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "自身へのメンションの場合Trueを返す",
			args: args{
				message: "<@U01L10XDBQR>\u00a0ruby\n```puts \"hello\"```",
				botID:   "U01L10XDBQR",
			},
			want: true,
		},
		{
			name: "自身へのメンションでない場合Falseを返す",
			args: args{
				message: "<@XXXXX0XDBQR>\u00a0ruby\n```puts \"hello\"```",
				botID:   "U01L10XDBQR",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if flag := isMentionToBot(tt.args.message, tt.args.botID); flag != tt.want {
				t.Fatalf("expect:%v != result:%v}", tt.want, flag)
			}
		})
	}

}
