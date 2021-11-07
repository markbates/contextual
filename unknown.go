package contextual

type Unknown string

func (u Unknown) String() string {
	return string(u)
}
