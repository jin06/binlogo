package message

import "testing"

func TestJson(t *testing.T) {
	msg := New()
	js, err := msg.Json()
	if err != nil {
		t.Error(err)
	}
	t.Log(js)
}

func TestJsonContent(t *testing.T) {
	msg := New()
	js, err := msg.JsonContent()
	if err != nil {
		t.Error(err)
	}
	t.Log(js)
}
