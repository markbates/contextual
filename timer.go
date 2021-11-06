package contextual

import "fmt"

type TimerCtx struct {
	Deadline interface{}
}

func (t TimerCtx) String() string {
	return fmt.Sprintf("WithTimeout(deadline: %v)", t.Deadline)
}
