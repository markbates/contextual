package contextual

import (
	"bytes"
	"fmt"
	"reflect"
)

type SignalCtx reflect.Value

func (ctx SignalCtx) String() string {
	rv := reflect.Value(ctx)

	bb := &bytes.Buffer{}
	bb.WriteString("SignalCtx(")
	fmt.Fprintf(bb, "%#v", rv)
	bb.WriteString(")")

	return bb.String()
}
