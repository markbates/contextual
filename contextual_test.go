package contextual

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_String(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	type CtxKey string

	ctx := context.Background()
	ctx = context.WithValue(ctx, CtxKey("id"), "42")
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, time.Second)
	defer cancel()

	ctx = context.WithValue(ctx, CtxKey("name"), "mary")
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	ctx, cancel = context.WithDeadline(ctx, time.Now().Add(time.Hour))
	defer cancel()

	ctx = context.WithValue(ctx, CtxKey("request_id"), "abc")

	bb := &bytes.Buffer{}
	p := &Printer{
		TimePrinter: func(t time.Time) (string, error) {
			return "TIME", nil
		},
	}

	err := p.Print(ctx, bb)
	r.NoError(err)

	act := bb.String()
	act = strings.TrimSpace(act)

	// fmt.Println(act)

	exp := `context
	.Background
		.WithValue(type contextual
			.CtxKey, val 42)
				.WithCancel
					.WithDeadline(TIME)
						.WithValue(type contextual
							.CtxKey, val mary)
								.WithCancel
									.WithCancel
										.WithValue(type contextual`
	r.Equal(exp, act)
}
