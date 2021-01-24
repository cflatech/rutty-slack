package client

import "testing"

func TestMakeRequestJson(t *testing.T) {
	expect := `{"code":"puts \"hello\""}`
	t.Run("正常なリクエスト用JSONに変換できる", func(t *testing.T) {
		if result := makeRequestJSON(`puts "hello"`); string(result) != expect {
			t.Fatalf("expect: %s != result: %s", expect, result)
		}
	})
}
