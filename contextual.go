package contextual

import (
	"bytes"
	"context"
	"io"
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

func Print(ctx context.Context, w io.Writer) error {
	p := &Printer{}

	return p.Print(ctx, w)
}
