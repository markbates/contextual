package contextual

import (
	"context"
	"fmt"
	"io"
	"os"
	"reflect"
)

type Printer struct {
	io.Writer
	DeadlinePrinter DeadlinePrinter
	chain           []fmt.Stringer
}

func (p *Printer) Print(ctx context.Context) error {
	if p.Writer == nil {
		p.Writer = os.Stdout
	}

	err := p.printCtx(ctx)

	if err != nil {
		return err
	}

	chain := p.chain

	var spaces Tabs

	for i, c := range chain {
		s := c.String()

		if i > 0 {
			s = "--> " + s
		}

		line := fmt.Sprintf("%s%s\n", spaces, s)
		fmt.Fprint(p.Writer, line)
		spaces = spaces.Increment()
	}

	return nil
}

type unknown string

func (u unknown) String() string {
	return string(u)
}

func (p *Printer) printCtx(ctx context.Context) error {
	rv := reflect.ValueOf(ctx)
	rt := rv.Type()
	rvi := reflect.Indirect(rv)

	name := fmt.Sprintf("%v", rt)

	switch name {
	case "*context.valueCtx":
		return p.withValue(rvi)
	case "*context.cancelCtx":
		return p.withCancel(rvi)
	case "*context.timerCtx":
		return p.withTimer(rvi)
	case "*context.emptyCtx":
		return p.withEmpty(rvi)
	default:
		p.chain = append(p.chain, unknown(name))
	}

	return nil
}

func (p *Printer) withTimer(rv reflect.Value) error {
	tc := TimerCtx{}

	f := rv.FieldByName("deadline")
	if !f.IsValid() {
		return fmt.Errorf("deadline field not found")
	}

	s := fmt.Sprintf("%+v", f)
	tc.Deadline = s

	if fn := p.DeadlinePrinter; fn != nil {
		x, err := fn(f)
		if err != nil {
			return err
		}
		tc.Deadline = x
	}

	p.chain = append(p.chain, tc)

	f = rv.FieldByName("Context")
	if !f.IsValid() {
		return nil
	}

	ctx, ok := f.Interface().(context.Context)
	if !ok {
		return fmt.Errorf("unexpected type %v", f.Type())
	}
	return p.printCtx(ctx)
}

func (p *Printer) withEmpty(rv reflect.Value) error {

	p.chain = append(p.chain, EmptyCtx{})

	return nil
}

func (p *Printer) withCancel(rv reflect.Value) error {
	cc := CancelCtx{}

	f := rv.FieldByName("Context")
	if !f.IsValid() {
		return nil
	}

	p.chain = append(p.chain, cc)
	ctx, ok := f.Interface().(context.Context)
	if !ok {
		return fmt.Errorf("unexpected type %v", f.Type())
	}

	return p.printCtx(ctx)
}

func (p *Printer) withValue(rv reflect.Value) error {
	vc := ValueCtx{}

	f := rv.FieldByName("key")
	if f.IsValid() {
		vc.key = fmt.Sprintf("%v", f)
	}

	f = rv.FieldByName("val")
	if f.IsValid() {
		vc.val = fmt.Sprintf("%v", f)
	}

	p.chain = append(p.chain, vc)

	f = rv.FieldByName("Context")
	if !f.IsValid() {
		return nil
	}

	ctx, ok := f.Interface().(context.Context)
	if !ok {
		return fmt.Errorf("unexpected type %v", f.Type())
	}

	return p.printCtx(ctx)
}
