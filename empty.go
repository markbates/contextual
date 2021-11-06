package contextual

type EmptyCtx struct{}

func (e EmptyCtx) String() string {
	return "Background"
}
