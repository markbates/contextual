package contextual

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Printer(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	type CtxKey string

	ctx := context.Background()
	ctx = context.WithValue(ctx, CtxKey("id"), "42")
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	ctx = context.WithValue(ctx, CtxKey("name"), "mary")
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	ctx, cancel = context.WithDeadline(ctx, time.Now().Add(time.Hour))
	defer cancel()

	ctx, cancel = signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	ctx = context.WithValue(ctx, CtxKey("request_id"), "abc")

	bb := &bytes.Buffer{}
	p := &Printer{
		DeadlinePrinter: func(deadline interface{}) (string, error) {
			return "TIME", nil
		},
		Writer: bb,
	}

	err := p.Print(ctx)
	r.NoError(err)

	act := bb.String()
	act = strings.TrimSpace(act)

	fmt.Println(act)

	exp := `WithValue(key: request_id, value: abc)
	--> SignalCtx([]os.Signal{2, 15})
		--> WithCancel
			--> WithCancel
				--> WithCancel
					--> WithValue(key: name, value: mary)
						--> WithTimeout(deadline: TIME)
							--> WithCancel
								--> WithValue(key: id, value: 42)
									--> context.backgroundCtx`
	r.Equal(exp, act)
}
