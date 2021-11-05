package contextual

import "bytes"

type Tabs []rune

func (t Tabs) String() string {
	bb := bytes.Buffer{}
	for _, r := range t {
		bb.WriteRune(r)
	}
	return bb.String()
}

func (i Tabs) Increment() Tabs {
	return append(i, '\t')
}
