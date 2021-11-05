package contextual

import (
	"bytes"
	"context"
)

func String(ctx context.Context) (string, error) {
	bb := &bytes.Buffer{}

	p := &Printer{}

	err := p.Print(ctx, bb)
	if err != nil {
		return "", err
	}

	return bb.String(), nil
}
