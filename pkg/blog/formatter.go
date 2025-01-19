package blog

import (
	"bytes"
	"fmt"

	"github.com/sirupsen/logrus"
)

var defaultFormat = "2006-01-02 15:04:05"

// formatter logrus formatter
type formatter struct {
}

func (f formatter) Format(e *logrus.Entry) (res []byte, err error) {
	var b *bytes.Buffer
	if e.Buffer != nil {
		b = e.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	caller := ""
	if e.HasCaller() {
		fileVal := fmt.Sprintf("%s:%d", e.Caller.File, e.Caller.Line)
		caller = " " + fileVal
	}

	fmt.Fprintf(b, "[binlogo-%.5s] [%s]%s %-44s", e.Level, e.Time.Format(defaultFormat), caller, e.Message)
	for k, v := range e.Data {
		b.WriteByte(' ')
		b.WriteString(k)
		b.WriteByte('=')
		stringVal, ok := v.(string)
		if !ok {
			stringVal = fmt.Sprint()
		}
		b.WriteString(fmt.Sprintf("%q", stringVal))
	}
	b.WriteByte('\n')
	res = b.Bytes()
	return
}
