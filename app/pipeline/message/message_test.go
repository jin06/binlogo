package message

import "testing"

func TestMessage(t *testing.T) {
	msg := New()
	if _, err := msg.Json(); err != nil {
		t.Error(err)
	}

	if _, err := msg.JsonContent(); err != nil {
		t.Error(err)
	}

	msg.Content = Content{
		Head: Head{
			Type:     "",
			Time:     0,
			Database: "",
			Table:    "",
			Position: nil,
		},
		Data: map[string]string{},
	}

	if _, err := msg.Json(); err != nil {
		t.Error(err)
	}
	if _, err := msg.JsonContent(); err != nil {
		t.Error(err)
	}
}

func TestString(t *testing.T) {
	msg := New()
	if msg.ToString() == "" {
		t.Fail()
	}
}

func TestMessageType(t *testing.T) {
	typ := TYPE_EMPTY
	t.Log(typ.String())
	typ = TYPE_INSERT
	t.Log(typ.String())
	typ = TYPE_ALTER_TABLE
	t.Log(typ.String())
	typ = TYPE_DELETE
	t.Log(typ.String())
	typ = TYPE_UPDATE
	t.Log(typ.String())
	typ = TYPE_CREATE_TABLE
	t.Log(typ.String())
}
