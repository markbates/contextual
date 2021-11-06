package contextual

import (
	"bytes"
	"context"
	"io"
)

func String(ctx context.Context) (string, error) {
	bb := &bytes.Buffer{}

	p := &Printer{
		Writer: bb,
	}

	err := p.Print(ctx)
	if err != nil {
		return "", err
	}

	return bb.String(), nil
}

func Print(ctx context.Context, w io.Writer) error {
	p := &Printer{
		Writer: w,
	}

	return p.Print(ctx)
}
