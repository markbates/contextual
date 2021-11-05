package contextual

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
)

type Printer struct {
	TimePrinter TimePrinter
}

func (p *Printer) Print(ctx context.Context, w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}

	s := fmt.Sprintf("%v", ctx)

	err := p.format(w, Tabs{}, s)
	if err != nil {
		return err
	}

	return nil
}

func (p *Printer) format(w io.Writer, spaces Tabs, s string) error {
	const match = `\.[Background|WithValue|WithCancel|WithDeadline].+`
	rx, err := regexp.Compile(match)

	if err != nil {
		return err
	}
	_ = rx

	loc := rx.FindStringIndex(s)
	if len(loc) == 0 {
		return nil
	}

	cur, next := s[:loc[0]], s[loc[0]+1:]

	if len(spaces) > 0 {
		cur = "." + cur
	}

	if strings.HasPrefix(cur, ".WithDeadline") {
		var err error
		cur, err = p.deadline(cur)
		if err != nil {
			return err
		}
	}

	fmt.Fprintf(w, "%+v%s\n", spaces, cur)

	return p.format(w, spaces.Increment(), next)
}

func (p *Printer) deadline(s string) (string, error) {

	const match = `.WithDeadline\((.+)\)`

	rx, err := regexp.Compile(match)
	if err != nil {
		return s, err
	}

	res := rx.FindStringSubmatch(s)
	if len(res) == 0 {
		return s, nil
	}

	second := res[1]

	trx, err := regexp.Compile(`(\sm=.+$)`)
	if err != nil {
		return s, err
	}

	trail := trx.FindString(second)
	clean := strings.TrimSuffix(second, trail)

	tf := `2006-01-02 15:04:05.999999999 -0700 MST`

	t, err := time.Parse(tf, clean)
	if err != nil {
		return s, err
	}

	fn := p.TimePrinter
	if fn == nil {
		fn = RFC3339Nano
	}

	ts, err := fn(t)
	if err != nil {
		return s, err
	}

	return fmt.Sprintf(".WithDeadline(%s)", ts), nil
}

/*
context
	.Background
		.WithValue(type string, val 42)
			.WithCancel
				.WithDeadline(2021-11-05 16:41:02.118779 -0400 EDT m=+1.002606096 [999.63316ms])
					.WithValue(type string, val mary)
						.WithCancel
*/
