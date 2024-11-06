package annotations

import "strconv"

// * Enabled

type (
	Enabled struct {
		aChan aChan
		value bool
	}
)

func (a *Annotation) Enabled() Enabled {
	ae := Enabled{
		aChan: make(aChan),
	}

	if v, ok := a.annotations[string(KeyEnabled)]; ok {
		boolValue, _ := strconv.ParseBool(v)
		ae.value = boolValue
	}

	go func() {
		for {
			select {
			case x := <-ae.aChan:
				a.annotations[string(x.key)] = x.value
			case <-a.ctx.Done():
				return
			}
		}
	}()

	return ae
}

func (a Enabled) Get() bool {
	return a.value
}

func (a Enabled) Set(enabled bool) {
	a.aChan.Send(KeyEnabled, strconv.FormatBool(enabled))
}
