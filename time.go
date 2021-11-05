package contextual

import "time"

type TimePrinter func(t time.Time) (string, error)

func RFC3339Nano(t time.Time) (string, error) {
	return t.Format(time.RFC3339Nano), nil
}
