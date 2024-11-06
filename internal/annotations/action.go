package annotations

import "strings"

// * Action

type (
	Action struct {
		aChan aChan
		value string
	}

	AActionKey string
)

const (
	// Action Refresh
	ActionRefresh AActionKey = "refresh"

	// Action Reload
	ActionReload AActionKey = "reload"

	// Action Delete
	ActionDelete AActionKey = "delete"
)

func (a *Annotation) Action() (ac *Action) {
	ac = &Action{
		aChan: make(aChan),
	}

	if v, ok := a.annotations[string(KeyAction)]; ok {
		ac.value = v
	}

	go func() {
		for {
			select {
			case x := <-ac.aChan:
				a.annotations[string(x.key)] = x.value
			case <-a.ctx.Done():
				return
			}
		}
	}()

	return ac
}

func (a *Action) Is(action AActionKey) bool {
	return strings.EqualFold(a.value, string(action))
}

func (a *Action) IsNull() bool {
	return a.value == ""
}

func (a *Action) Get() AActionKey {
	return AActionKey(a.value)
}

func (a *Action) Set(action AActionKey) {
	a.aChan.Send(KeyAction, string(action))
}
