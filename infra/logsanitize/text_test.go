package logsanitize

import "testing"

func TestForLog(t *testing.T) {
	input := "ok\x00\t\x01\x7f"
	got := ForLog(input)
	want := `ok\0\t\x01\x7f`
	if got != want {
		t.Fatalf("unexpected sanitize result: got=%q want=%q", got, want)
	}
}

func TestForLogInvalidUtf8(t *testing.T) {
	input := string([]byte{0xff, 'a'})
	got := ForLog(input)
	if got != "�a" {
		t.Fatalf("unexpected utf8 sanitize result: got=%q", got)
	}
}

