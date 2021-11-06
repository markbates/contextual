package contextual

import "fmt"

type ValueCtx struct {
	key interface{}
	val interface{}
}

func (ctx ValueCtx) String() string {
	return fmt.Sprintf("WithValue(key: %v, value: %v)", ctx.key, ctx.val)
}
