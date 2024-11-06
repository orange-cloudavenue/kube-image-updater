package annotations

import "strings"

// * FailurePolicy

type (
	FailurePolicy struct {
		aChan aChan
		value string
	}

	AFailurePolicyKey string
)

const (
	// FailurePolicy Ignore
	FailurePolicyIgnore AFailurePolicyKey = "ignore"

	// FailurePolicy Fail
	FailurePolicyFail AFailurePolicyKey = "fail"
)

func (a *Annotation) FailurePolicy() (af *FailurePolicy) {
	af = &FailurePolicy{
		aChan: make(aChan),
	}

	if v, ok := a.annotations[string(KeyFailurePolicy)]; ok {
		af.value = v
	}

	go func() {
		for {
			select {
			case x := <-af.aChan:
				a.annotations[string(x.key)] = x.value
			case <-a.ctx.Done():
				return
			}
		}
	}()

	return af
}

func (a *FailurePolicy) Is(policy AFailurePolicyKey) bool {
	return strings.EqualFold(a.value, string(policy))
}

func (a *FailurePolicy) IsNull() bool {
	return a.value == ""
}

func (a *FailurePolicy) Get() AFailurePolicyKey {
	return AFailurePolicyKey(a.value)
}

func (a *FailurePolicy) Set(policy AFailurePolicyKey) {
	a.aChan.Send(KeyFailurePolicy, string(policy))
}
