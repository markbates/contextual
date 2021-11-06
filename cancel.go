package contextual

type CancelCtx struct {
}

func (ctx CancelCtx) String() string {
	return "WithCancel"

}
